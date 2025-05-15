package routes

import (
	"context"

	"github.com/ONSdigital/dp-frontend-cookie-controller/handlers"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	render "github.com/ONSdigital/dp-renderer/v2"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// Init initialises routes for the service
func Init(ctx context.Context, r *mux.Router, hc health.HealthCheck, rendC *render.Render) {
	log.Info(ctx, "adding api routes")

	r.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)

	r.StrictSlash(true).Path("/cookies").Methods("GET").HandlerFunc(handlers.Read(rendC))
	r.StrictSlash(true).Path("/cookies").Methods("POST").HandlerFunc(handlers.Edit(rendC))
}
