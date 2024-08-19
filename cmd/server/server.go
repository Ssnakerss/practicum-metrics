package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ssnakerss/practicum-metrics/internal/compression"
	"github.com/Ssnakerss/practicum-metrics/internal/dataadapter"
	"github.com/Ssnakerss/practicum-metrics/internal/flags"
	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/Ssnakerss/practicum-metrics/internal/metric"
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

	//Configuring storage
	da := dataadapter.Adapter{}
	memst := &storage.MemStorage{}
	// //       ↑
	// //https://stackoverflow.com/questions/40823315/x-does-not-implement-y-method-has-a-pointer-receiver
	memst.New()
	// //cannot use memst (variable of type storage.MemStorage) as storage.DataStorage value in
	// // argument to da.New: storage.MemStorage does not implement storage.DataStorage
	// //(method Insert has pointer receiver)memst.New()
	// //       ↓
	da.New(memst)

	//Используем файл для хранения метрик
	// filest := &storage.FileStorage{}
	// filest.New(`flags.Cfg.FileStoragePath`)
	// da.New(filest)

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

	filest := &storage.FileStorage{}
	filest.New(flags.Cfg.FileStoragePath)

	if flags.Cfg.Restore {
		err := copyState(filest, memst)
		if err != nil {
			logger.SLog.Warnw("data failed", "restore", err)
		}
	}

	go intervalSave(flags.Cfg.StoreInterval, memst, filest)

	//------------Program exit code------------------
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-exit
		logger.SLog.Infow("received termination", "signal", sig)
		copyState(memst, filest)
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

// Сохранение и восстановление состояния хранилища в/из файла
func intervalSave(interval int, src storage.DataStorage, dst storage.DataStorage) {
	if interval < 1 {
		return
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for {
		<-ticker.C
		copyState(src, dst)
	}
}

func copyState(src storage.DataStorage, dst storage.DataStorage) error {
	mm := make([]metric.Metric, 0)
	readcnt, err := src.ReadAll(&mm)
	if err != nil {
		return err
	}
	writecnt, err := dst.WriteAll(&mm)
	if err != nil {
		return err
	}
	if readcnt != writecnt {
		return fmt.Errorf("read count %d not equal to write count%d", readcnt, writecnt)
	}

	return nil
}
