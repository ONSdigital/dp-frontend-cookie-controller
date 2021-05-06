package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-frontend-cookie-controller
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	RendererURL                string        `envconfig:"RENDERER_URL"`
	APIRouterURL               string        `envconfig:"API_ROUTER_URL"`
	SiteDomain                 string        `envconfig:"SITE_DOMAIN"`
	PatternLibraryAssetsPath   string        `envconfig:"PATTERN_LIBRARY_ASSETS_PATH"`
	SupportedLanguages         [2]string     `envconfig:"SUPPORTED_LANGUAGES"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
}

var cfg *Config

// Get returns the default config with any modifications through environment variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg := &Config{
		BindAddr:                   ":24100",
		RendererURL:                "http://localhost:20010",
		APIRouterURL:               "http://localhost:22400",
		SiteDomain:                 "localhost",
		SupportedLanguages:         [2]string{"en", "cy"},
		PatternLibraryAssetsPath:   "http://localhost:9000/dist",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
	}

	return cfg, envconfig.Process("", cfg)
}
