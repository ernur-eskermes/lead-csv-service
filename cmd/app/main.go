package main

import (
	"context"
	"encoding/csv"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/gocarina/gocsv"
	log "github.com/sirupsen/logrus"

	"github.com/ernur-eskermes/lead-csv-service/pkg/database/postgresql"

	_ "github.com/ernur-eskermes/lead-csv-service/docs"
	"github.com/ernur-eskermes/lead-csv-service/internal/transport/rest"
	restHandler "github.com/ernur-eskermes/lead-csv-service/internal/transport/rest/handlers"

	"github.com/ernur-eskermes/lead-csv-service/internal/config"
	"github.com/ernur-eskermes/lead-csv-service/internal/service"
	"github.com/ernur-eskermes/lead-csv-service/internal/storage"
	"github.com/jackc/pgx/v4/log/logrusadapter"
)

func init() {
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})

	gocsv.SetCSVWriter(func(out io.Writer) *gocsv.SafeCSVWriter {
		writer := csv.NewWriter(out)
		writer.Comma = ';'

		return gocsv.NewSafeCSVWriter(writer)
	})
}

// @title Product app REST-API
// @version 1.0
// @description Application for adding/getting products and download CSV-file

// @host localhost:8000
// @BasePath /api/

// Run initializes whole application.
func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgresql.NewClient(context.TODO(), postgresql.StorageConfig{
		ConnStr:     cfg.Postgres.URI,
		MaxAttempts: 5,
		Logger:      logrusadapter.NewLogger(log.StandardLogger()),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	storages := storage.New(db)
	services := service.New(service.Deps{
		ProductStorage: storages.Product,
	})

	restHandlers := restHandler.New(restHandler.Deps{
		ProductService: services.Product,
	})
	restSrv := rest.NewServer(cfg, restHandlers)

	go func() {
		log.Info("Starting HTTP server")

		if err = restSrv.ListenAndServe(cfg.HTTP.Port); err != nil {
			log.Error("HTTP ListenAndServer error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	log.Info("Shutting down server")

	if err = restSrv.Stop(); err != nil {
		log.Errorf("failed to stop http server: %v", err)
	}
}
