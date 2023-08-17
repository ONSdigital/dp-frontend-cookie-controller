package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ONSdigital/dp-frontend-cookie-controller/assets"
	"github.com/ONSdigital/dp-frontend-cookie-controller/config"
	"github.com/ONSdigital/dp-frontend-cookie-controller/routes"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	render "github.com/ONSdigital/dp-renderer/v2"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"

	dpnethttp "github.com/ONSdigital/dp-net/v2/http"
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

	versionInfo, err := health.NewVersionInfo(
		BuildTime,
		GitCommit,
		Version,
	)

	r := mux.NewRouter()

	//nolint:typecheck
	rendC := render.NewWithDefaultClient(assets.Asset, assets.AssetNames, cfg.PatternLibraryAssetsPath, cfg.SiteDomain)

	healthcheck := health.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	routes.Init(ctx, r, healthcheck, rendC)
	healthcheck.Start(ctx)

	s := dpnethttp.NewServer(cfg.BindAddr, r)
	s.HandleOSSignals = false

	log.Info(ctx, "Starting server", log.Data{"config": cfg})

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Error(ctx, "failed to start http listen and serve", err)
			return
		}
	}()

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
