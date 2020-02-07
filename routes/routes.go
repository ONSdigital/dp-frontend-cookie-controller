package routes

import (
    "context"

	"dp-frontend-cookie-controller/config"
	"dp-frontend-cookie-controller/handlers"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

// Init initialises routes for the service
func Init(ctx context.Context, r *mux.Router, cfg *config.Config, hc health.HealthCheck) {
    log.Event(ctx, "adding routes")
	r.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)

	r.StrictSlash(true).Path("/cookie/accept-all").Methods("GET").HandlerFunc(handlers.AcceptAll())
	r.StrictSlash(true).Path("/cookie/edit").Methods("GET").HandlerFunc(handlers.Read())
	r.StrictSlash(true).Path("/cookie/edit").Methods("POST").HandlerFunc(handlers.Edit())}
