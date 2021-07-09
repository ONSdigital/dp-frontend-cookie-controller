package routes

import (
	"context"

	"dp-frontend-cookie-controller/handlers"

	render "github.com/ONSdigital/dp-renderer"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

// Init initialises routes for the service
func Init(ctx context.Context, r *mux.Router, hc health.HealthCheck, rendC *render.Render) {
	log.Event(ctx, "adding api routes")

	r.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)

	r.StrictSlash(true).Path("/cookies").Methods("GET").HandlerFunc(handlers.Read(rendC))
	r.StrictSlash(true).Path("/cookies").Methods("POST").HandlerFunc(handlers.Edit(rendC, rendC.SiteDomain))
}
