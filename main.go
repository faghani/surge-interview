package main

import (
	"fmt"

	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/faghani/surge-interview/internal/database"
	"github.com/faghani/surge-interview/internal/redis"
	"github.com/faghani/surge-interview/internal/nats"
	"github.com/faghani/surge-interview/internal/router"
	"github.com/faghani/surge-interview/internal/service"
	log "github.com/sirupsen/logrus"
)

const (
	apiVersion = "/v1"
	timeWindow = 120 * time.Second
	keyExpiration = 3 * time.Minute
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	rc, err := redis.New(os.Getenv("REDIS_ADDRESS"))
	if err != nil {
		log.WithError(err).Fatal("failed to connect to redis")
	}

	nc, err := nats.New(os.Getenv("NATS_ADDRESS"))
	if err != nil {
		log.WithError(err).Fatal("failed to connect to nats")
	}

	db, err := database.New(
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASS"),
		os.Getenv("MYSQL_NAME"),
		os.Getenv("MYSQL_PORT"),
		)
	if err != nil {
		log.WithError(err).Fatal("failed to connect to mysql")
	}

	if err := db.Migrate(os.Getenv("MIGRATION_DIRECTORY")); err != nil {
		log.WithError(err).Fatal("failed to run database migration")
	}

	if err := db.Setup(); err != nil {
		log.WithError(err).Fatal("failed to setup database")
	}

	r := &router.Router{
		DemandService: &service.DemandService{
			Redis: rc,
			Nats: nc,
			Db: db,
			Timestep: timeWindow,
			Expiration: keyExpiration,
		},
		Db: db,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT")),
		Handler:      r.Handler(apiVersion),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go server.ListenAndServe()
	log.Infof("started router, listening on %s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))

	signalReceiver := make(chan os.Signal, 1)
	signal.Notify(signalReceiver, os.Interrupt)

	<-signalReceiver
	server.Close()
	rc.Close()
	nc.Close()
	db.DB.Close()
	log.Error("router stopped")
}