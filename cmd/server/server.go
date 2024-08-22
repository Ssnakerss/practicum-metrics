package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ssnakerss/practicum-metrics/internal/compression"
	"github.com/Ssnakerss/practicum-metrics/internal/dtadapter"
	"github.com/Ssnakerss/practicum-metrics/internal/flags"
	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/Ssnakerss/practicum-metrics/internal/storage"

	"github.com/go-chi/chi/v5"
)

func main() {
	// cоздаем логгер ZAP
	// не получится - проолжать не имеет смысла, fatal
	if err := logger.Initialize("DEBUG"); err != nil {
		log.Fatal("FATAL: cannot initialize LOGGER: ", err)
	}

	defer logger.Log.Sync()
	//Перехватываем паники
	defer func() {
		if err := recover(); err != nil {
			logger.SLog.Fatalf(
				"error heppened while operating -> program will exit",
				"error", err)
		}
	}()
	//Читаем конфигурацию
	if err := flags.ReadServerConfig(); err != nil {
		logger.SLog.Warnw("error getting env params", "error", err)
	}

	//Создаем хранилище
	//Хранилище должно соответствовать интерфейсу storage.DataStorage
	memst := &storage.MemStorage{}
	memst.New()

	filest := &storage.FileStorage{}
	filest.New(flags.Cfg.FileStoragePath)
	// filest.Truncate()

	da := dtadapter.Adapter{}

	if flags.Cfg.Restore {
		err := da.CopyState(filest, memst)
		if err != nil {
			logger.SLog.Warnw("data failed", "restore", err)
		}
	}

	da.New(memst)
	//Добавляем хранилище и включаем синхронизацию
	//0 - пишем в оба сразе, > 0 - по расписанию
	da.Sync(flags.Cfg.StoreInterval, filest)

	//Configuring CHI
	r := chi.NewRouter()
	r.Get("/",
		logger.WithLogging(
			compression.GzipHandle(
				http.HandlerFunc(da.MainPage))))
	// JSON handlers
	r.Post("/update/", logger.WithLogging(
		http.HandlerFunc(da.SetDataJSONHandler)))

	r.Post("/value/",
		logger.WithLogging(
			compression.GzipHandle(
				http.HandlerFunc(da.GetDataJSONHandler))))

	//TEXT handlers
	r.Get("/value/{type}/{name}",
		logger.WithLogging(
			compression.GzipHandle(
				http.HandlerFunc(da.GetDataTextHandler))))

	r.Post("/update/{type}/{name}/{value}", logger.WithLogging(http.HandlerFunc(da.SetDataTextHandler)))

	//Сохраняем состояние оперативного хранилища на диске при выходе из программы-
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-exit
		logger.SLog.Infow("received termination", "signal", sig)
		filest.Truncate()
		da.CopyState(memst, filest)
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
