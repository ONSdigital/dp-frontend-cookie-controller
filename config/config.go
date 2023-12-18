package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-frontend-cookie-controller
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	Debug                      bool          `envconfig:"DEBUG"`
	SiteDomain                 string        `envconfig:"SITE_DOMAIN"`
	PatternLibraryAssetsPath   string        `envconfig:"PATTERN_LIBRARY_ASSETS_PATH"`
	SupportedLanguages         [2]string     `envconfig:"SUPPORTED_LANGUAGES"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	OTExporterOTLPEndpoint     string        `envconfig:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	OTServiceName              string        `envconfig:"OTEL_SERVICE_NAME"`
	OTBatchTimeout             time.Duration `envconfig:"OTEL_BATCH_TIMEOUT"`
}

var cfg *Config

// Get returns the default config with any modifications through environment variables
func Get() (*Config, error) {
	config, err := get()
	if err != nil {
		return nil, err
	}

	if config.Debug {
		config.PatternLibraryAssetsPath = "http://localhost:9002/dist/assets"
	} else {
		config.PatternLibraryAssetsPath = "//cdn.ons.gov.uk/dp-design-system/e0a75c3"
	}
	return config, nil
}

func get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   ":24100",
		Debug:                      false,
		SiteDomain:                 "localhost",
		SupportedLanguages:         [2]string{"en", "cy"},
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		OTExporterOTLPEndpoint:     "localhost:4317",
		OTServiceName:              "dp-frontend-cookie-controller",
		OTBatchTimeout:             5 * time.Second,
	}

	return cfg, envconfig.Process("", cfg)
}
