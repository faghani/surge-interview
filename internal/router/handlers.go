package router

import (
	"encoding/json"
	"github.com/faghani/surge-interview/internal/database"
	"github.com/faghani/surge-interview/internal/osm"
	"github.com/faghani/surge-interview/internal/service"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type AddDemandHandler struct {
	DemandService *service.DemandService
	Db *database.Database
}

func (d AddDemandHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	dbyte, _ := ioutil.ReadAll(r.Body)

	var loc Location
	if err := json.Unmarshal(dbyte, &loc); err != nil {
		res, _ := json.Marshal(StatusMessage{false})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	go func() {
		relationId := osm.FindOsmRelationID(loc.Latitude, loc.Longitude)
		if relationId != "" {
			if err := d.DemandService.Increment(osm.NormalizeOsmRelationId(relationId)); err != nil {
				log.WithError(err).Errorf("failed to increment demand for %d", osm.NormalizeOsmRelationId(relationId))
			}
		}
	}()

	res, _ := json.Marshal(StatusMessage{true})
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

type RuleHandler struct {
	Db *database.Database
}

func (h RuleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	dbyte, _ := ioutil.ReadAll(r.Body)
	var req RulesRequest
	if err := json.Unmarshal(dbyte, &req); err != nil {
		res, _ := json.Marshal(StatusMessage{false})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	err := h.Db.SaveRules(req.Rules)
	if err != nil {
		res, _ := json.Marshal(StatusMessage{false})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	res, _ := json.Marshal(StatusMessage{true})
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

type OsmRelationHandler struct {
}

func (o OsmRelationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	dbyte, _ := ioutil.ReadAll(r.Body)
	var loc Location
	if err := json.Unmarshal(dbyte, &loc); err != nil {
		res, _ := json.Marshal(OsmRelationIdMessage{
			Status: StatusMessage{false},
		})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	relationId := osm.FindOsmRelationID(loc.Latitude, loc.Longitude)
	res, _ := json.Marshal(OsmRelationIdMessage{
		Status: StatusMessage{true},
		RelationId: relationId,
	})

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}