package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ONSdigital/dp-frontend-cookie-controller/assets"
	"github.com/ONSdigital/dp-frontend-cookie-controller/config"
	"github.com/ONSdigital/dp-frontend-cookie-controller/routes"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpnethttp "github.com/ONSdigital/dp-net/v3/http"
	dpotelgo "github.com/ONSdigital/dp-otel-go"
	render "github.com/ONSdigital/dp-renderer/v2"
	"github.com/ONSdigital/dp-renderer/v2/middleware/renderror"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Error(ctx, "unable to run application", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	log.Namespace = "dp-frontend-cookie-controller"

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config.Get()
	if err != nil {
		log.Error(ctx, "unable to retrieve service configuration", err)
		return err
	}

	log.Info(ctx, "got service configuration", log.Data{"config": cfg})

	if cfg.OtelEnabled {
		// Set up OpenTelemetry
		otelConfig := dpotelgo.Config{
			OtelServiceName:          cfg.OTServiceName,
			OtelExporterOtlpEndpoint: cfg.OTExporterOTLPEndpoint,
			OtelBatchTimeout:         cfg.OTBatchTimeout,
		}

		otelShutdown, err := dpotelgo.SetupOTelSDK(ctx, otelConfig)

		if err != nil {
			log.Error(ctx, "error setting up OpenTelemetry - hint: ensure OTEL_EXPORTER_OTLP_ENDPOINT is set", err)
		}

		// Handle shutdown properly so nothing leaks.
		defer func() {
			err = errors.Join(err, otelShutdown(context.Background()))
		}()
	}

	versionInfo, err := health.NewVersionInfo(
		BuildTime,
		GitCommit,
		Version,
	)
	if err != nil {
		log.Error(ctx, "failed to create service version information", err)
		return err
	}

	r := mux.NewRouter()

	if cfg.OtelEnabled {
		r.Use(otelmux.Middleware(cfg.OTServiceName))
	}

	rendC := render.NewWithDefaultClient(assets.Asset, assets.AssetNames, cfg.PatternLibraryAssetsPath, cfg.SiteDomain)

	healthcheck := health.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	routes.Init(ctx, r, healthcheck, rendC)
	healthcheck.Start(ctx)

	middleware := []alice.Constructor{
		renderror.Handler(rendC),
	}

	newAlice := alice.New(middleware...).Then(r)

	var s *dpnethttp.Server
	if cfg.OtelEnabled {
		otelHandler := otelhttp.NewHandler(newAlice, "/")
		s = dpnethttp.NewServer(cfg.BindAddr, otelHandler)
	} else {
		s = dpnethttp.NewServer(cfg.BindAddr, newAlice)
	}

	s.HandleOSSignals = false

	log.Info(ctx, "Starting server", log.Data{"config": cfg})

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Error(ctx, "failed to start http listen and serve", err)
			return
		}
	}()

	//nolint:gosimple // ignore this as intention is to continue to listen for signals
	for {
		select {
		case <-signals:
			log.Info(ctx, "os signal received")
			return gracefulShutdown(cfg, s, healthcheck)
		}
	}
}

func gracefulShutdown(cfg *config.Config, s *dpnethttp.Server, hc health.HealthCheck) error {
	log.Info(context.Background(), fmt.Sprintf("shutdown with timeout: %s", cfg.GracefulShutdownTimeout))
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GracefulShutdownTimeout)
	log.Info(ctx, "shutting service down gracefully")
	defer cancel()

	// Stop health check tickers
	hc.Stop()
	if err := s.Server.Shutdown(ctx); err != nil {
		log.Error(ctx, "failed to shutdown http server", err)
		return err
	}
	log.Info(ctx, "graceful shutdown complete successfully")
	return nil
}
