package router

import (
	"github.com/faghani/surge-interview/internal/database"
	"github.com/faghani/surge-interview/internal/service"
	"github.com/gorilla/mux"
	"net/http"
)

type Location struct {
	Latitude float64 `json:"lat"`
	Longitude float64 `json:"long"`
}

type StatusMessage struct {
	Success bool `json:"success"`
}

type OsmRelationIdMessage struct {
	Status StatusMessage `json:"status"`
	RelationId string `json:"relation_id"`
}

type RulesRequest struct {
	Rules []database.Rule `json:"rules"`
}

type Router struct {
	DemandService *service.DemandService
	Db *database.Database
}

func (r *Router) Handler(apiVersion string) *mux.Router {
	router := mux.NewRouter().PathPrefix(apiVersion).Subrouter()
	router.Methods(http.MethodPost).Path("/demand").Handler(AddDemandHandler{r.DemandService, r.Db})
	router.Methods(http.MethodPost).Path("/rule").Handler(RuleHandler{r.Db})
	router.Methods(http.MethodPost).Path("/osm-relation-id").Handler(OsmRelationHandler{})
	router.Methods(http.MethodGet).Path("/health").HandlerFunc(r.handleHealth)
	return router
}

func (ro *Router) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}