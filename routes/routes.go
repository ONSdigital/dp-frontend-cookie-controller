package routes

import (
	"context"

	"dp-frontend-cookie-controller/assets"

	"dp-frontend-cookie-controller/config"
	"dp-frontend-cookie-controller/handlers"

	"github.com/rav-pradhan/test-modules/render"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

// Init initialises routes for the service
func Init(ctx context.Context, r *mux.Router, cfg *config.Config, hc health.HealthCheck) {
	log.Event(ctx, "adding api routes")

	rendC := render.New(cfg.PatternLibraryAssetsPath, cfg.SiteDomain, assets.Asset, assets.AssetNames)

	r.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)

	r.StrictSlash(true).Path("/cookies").Methods("GET").HandlerFunc(handlers.Read(cfg, rendC))
	r.StrictSlash(true).Path("/cookies").Methods("POST").HandlerFunc(handlers.Edit(cfg, rendC))
}
