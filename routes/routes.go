package routes

import (
	"context"

	"dp-frontend-cookie-controller/config"
	"dp-frontend-cookie-controller/handlers"
	"github.com/ONSdigital/dp-api-clients-go/renderer"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

// Init initialises routes for the service
func Init(ctx context.Context, r *mux.Router, cfg *config.Config, hc health.HealthCheck) {
	log.Event(ctx, "adding routes")

	rendC := renderer.New(cfg.RendererURL)

	r.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)

	r.StrictSlash(true).Path("/cookies/accept-all").Methods("GET").HandlerFunc(handlers.AcceptAll())
	r.StrictSlash(true).Path("/cookies/edit").Methods("GET").HandlerFunc(handlers.Read(rendC))
	r.StrictSlash(true).Path("/cookies/edit").Methods("POST").HandlerFunc(handlers.Edit())
}
