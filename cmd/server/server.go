package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/Ssnakerss/practicum-metrics/internal/compression"
	"github.com/Ssnakerss/practicum-metrics/internal/handlers"
	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
)

type Config struct {
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	//ADDRESS отвечает за адрес эндпоинта HTTP-сервера.
	EndPointAddress string `env:"ADDRESS"`
}

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
	cfg := Config{
		StoreInterval: -1,
	}

	//Сначала считаем командную строку если есть или заполним конфиг дефолтом
	//Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080)
	flag.StringVar(&cfg.EndPointAddress, "a", "localhost:8080", "endpoint address")
	//Флаг -i=<ЗНАЧЕНИЕ> интервал времени в секундах, по истечении которого текущие показания
	//сервера сохраняются на диск (по умолчанию 300 секунд, значение 0 делает запись синхронной)
	flag.IntVar(&cfg.StoreInterval, "i", 300, "data store interval, sec")
	//Флаг -f=<ЗНАЧЕНИЕ> путь до файла, куда сохраняются текущие значения.
	flag.StringVar(&cfg.FileStoragePath, "f", "", "file storage path")
	//Флаг -r=<ЗНАЧЕНИЕ>  булево значение (true/false), определяющее, загружать или нет ранее
	//сохранённые значения из указанного файла при старте сервера (по умолчанию true)
	flag.BoolVar(&cfg.Restore, "r", true, "restore data on startup")
	flag.Parse()

	//Потом пробуем прочитать ENV
	//Если есть - перепишут имеющиеся значения (?)
	err := env.Parse(&cfg)
	if err != nil {
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

	handlers.InitStorage(cfg.FileStoragePath, cfg.StoreInterval == 0)
	if cfg.Restore {
		err = handlers.Stor.Restore()
		if err != nil {
			logger.SLog.Warnw("data failed", "restore", err)
		}
	}

	go intervalSave(cfg.StoreInterval)

	defer handlers.Stor.Save()

	logger.SLog.Infow("starting server ", "config", cfg)
	if err := http.ListenAndServe(cfg.EndPointAddress, r); err != nil {
		logger.SLog.Fatalf(
			"failed to start server -> program will exit",
			"address", cfg.EndPointAddress,
			"error", err,
		)
	}
}

// ----------------------------------------
func intervalSave(interval int) {
	if interval < 1 {
		return
	}
	for {
		time.Sleep(time.Duration(interval) * time.Second)
		handlers.Stor.Save()
	}
}
