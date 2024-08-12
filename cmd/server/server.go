package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ssnakerss/practicum-metrics/internal/compression"
	"github.com/Ssnakerss/practicum-metrics/internal/flags"
	"github.com/Ssnakerss/practicum-metrics/internal/handlers"
	"github.com/Ssnakerss/practicum-metrics/internal/logger"

	"github.com/go-chi/chi/v5"
)

func main() {
	// cоздаем логгер ZAP
	// не получится - проолжать не имеет смысла, fatal
	if err := logger.Initialize("DEBUG"); err != nil {
		log.Fatal("FATAL: cannot initialize LOGGER: ", err)
	}

	defer logger.Log.Sync()
	//If any panic happened during opeartion
	defer func() {
		if err := recover(); err != nil {
			logger.SLog.Fatalf(
				"error heppened while operating -> program will exit",
				"error", err)
		}
	}()
	//Reading configuration

	if err := flags.ReadServerConfig(); err != nil {
		logger.SLog.Warnw("error getting env params", "error", err)
	}

	//Configuring CHI
	r := chi.NewRouter()
	r.Get("/",
		logger.WithLogging(
			compression.GzipHandle(
				http.HandlerFunc(handlers.MainPage))))

	r.Post("/update/", logger.WithLogging(http.HandlerFunc(handlers.SetDataJSONHandler)))

	r.Post("/value/",
		logger.WithLogging(
			compression.GzipHandle(
				http.HandlerFunc(handlers.GetDataJSONHandler))))

	r.Get("/value/{type}/{name}",
		logger.WithLogging(
			compression.GzipHandle(
				http.HandlerFunc(handlers.GetDataTextHandler))))

	r.Post("/update/{type}/{name}/{value}", logger.WithLogging(http.HandlerFunc(handlers.SetDataTextHandler)))

	handlers.InitStorage(flags.Cfg.FileStoragePath, flags.Cfg.StoreInterval == 0)
	if flags.Cfg.Restore {
		err := handlers.Stor.Restore()
		if err != nil {
			logger.SLog.Warnw("data failed", "restore", err)
		}
	}

	go intervalSave(flags.Cfg.StoreInterval)

	//------------Program exit code------------------
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-exit
		logger.SLog.Infow("received termination", "signal", sig)
		handlers.Stor.Save()
		logger.Log.Fatal("program will exit")
	}()
	//---------------------------------

	logger.SLog.Infow("starting server ", "config", flags.Cfg)
	if err := http.ListenAndServe(flags.Cfg.EndPointAddress, r); err != nil {
		logger.SLog.Fatalf(
			"failed to start server -> program will exit",
			"address", flags.Cfg.EndPointAddress,
			"error", err,
		)
	}
}

// ----------------------------------------
func intervalSave(interval int) {
	if interval < 1 {
		return
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for {
		<-ticker.C
		handlers.Stor.Save()
	}
}
