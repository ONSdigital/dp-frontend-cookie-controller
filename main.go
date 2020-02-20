package main

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-api-clients-go/renderer"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"os"
	"os/signal"
	"syscall"

	"dp-frontend-cookie-controller/config"
	"dp-frontend-cookie-controller/routes"
	"github.com/ONSdigital/go-ns/server"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
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
		log.Event(ctx, "unable to run application", log.Error(err))
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	log.Namespace = "dp-frontend-cookie-controller"

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config.Get()
	if err != nil {
		log.Event(ctx, "unable to retrieve service configuration", log.Error(err))
		return err
	}

	log.Event(ctx, "got service configuration", log.Data{"config": cfg})

	versionInfo, err := health.NewVersionInfo(
		BuildTime,
		GitCommit,
		Version,
	)

	r := mux.NewRouter()

	rend := renderer.New(cfg.RendererURL)

	healthcheck := health.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	if err = registerCheckers(ctx, &healthcheck, rend); err != nil {
		return err
	}
	routes.Init(ctx, r, cfg, healthcheck)

	healthcheck.Start(ctx)

	s := server.New(cfg.BindAddr, r)
	s.HandleOSSignals = false

	log.Event(ctx, "Starting server", log.Data{"config": cfg})

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Event(ctx, "failed to start http listen and serve", log.Error(err))
			return
		}
	}()

	for {
		select {
		case <-signals:
			log.Event(nil, "os signal received")
			return gracefulShutdown(cfg, s, healthcheck)
		}
	}
	// protective programming, shouldn't get to this... but just in case
	// nil translates to exit code 0
	return nil
}

func gracefulShutdown(cfg *config.Config, s *server.Server, hc health.HealthCheck) error {
	log.Event(nil, fmt.Sprintf("shutdown with timeout: %s", cfg.GracefulShutdownTimeout))
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GracefulShutdownTimeout)
	log.Event(ctx, "shutting service down gracefully")
	defer cancel()

	// Stop health check tickers
	hc.Stop()
	if err := s.Server.Shutdown(ctx); err != nil {
		log.Event(ctx, "failed to shutdown http server", log.Error(err))
	}
	return nil
}

func registerCheckers(ctx context.Context, h *health.HealthCheck, r *renderer.Renderer) (err error) {
	if err = h.AddCheck("frontend renderer", r.Checker); err != nil {
		log.Event(ctx, "failed to add frontend renderer checker", log.Error(err))
	}
	return
}
