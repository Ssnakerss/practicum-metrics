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

	var st storage.DataStorage

	da := dtadapter.Adapter{}
	var filest *storage.FileStorage

	//Если задан DSN - используем БД в качестве хранилища
	if flags.Cfg.DatabaseDSN != "default" {
		st = &storage.DBStorage{}
		//Ставим таймаут 60 секунд
		if err := st.New(flags.Cfg.DatabaseDSN, "120"); err != nil {
			logger.SLog.Fatalf(
				"error initialize db -> program will exit",
				"dsn", flags.Cfg.DatabaseDSN,
				"error", err)
		}
		//Очищаем таблицу   -  ???
		st.Truncate()
		logger.SLog.Info("using db as storage")
	} else {
		//Иначе используем хранение в памяти
		st = &storage.MemStorage{}
		st.New()
		logger.SLog.Info("using memory as storage")

		//Если задан путь к файлу - добавляем фаловое хранилище
		if flags.Cfg.FileStoragePath != "default" {
			filest = &storage.FileStorage{}
			if err := filest.New(flags.Cfg.FileStoragePath); err != nil {
				logger.SLog.Warnw("file creation failure", "path", flags.Cfg.FileStoragePath, "err", err)
			} else {
				if flags.Cfg.Restore {
					//Восстанавливаем значения из файла
					err := da.CopyState(filest, st)
					logger.SLog.Infow("restoring data from ", "file", filest.Filename)
					if err != nil {
						logger.SLog.Warnw("data restore", "failed", err)
					}
				}
			}
		}
	}

	da.New(st)
	if filest != nil {
		//Добавляем хранилище и включаем синхронизацию
		//0 - пишем в оба сразе, > 0 - по расписанию
		da.Sync(flags.Cfg.StoreInterval, filest)
		logger.SLog.Infow("using a sync storage", "file", filest.Filename)
	}

	//--------------------------------------------

	//Configuring CHI
	r := chi.NewRouter()
	r.Get("/",
		logger.WithLogging(
			compression.GzipHandle(
				http.HandlerFunc(da.MainPage))))

	//Проверка соединения с хранилищем
	r.Get("/ping",
		logger.WithLogging(
			http.HandlerFunc(da.Ping)))

	// JSON handlers
	r.Post("/update/",
		logger.WithLogging(
			http.HandlerFunc(da.SetDataJSONHandler)))

	r.Post("/updates/",
		logger.WithLogging(
			http.HandlerFunc(da.SetDataJSONSliceHandler)))

	r.Post("/value/",
		logger.WithLogging(
			compression.GzipHandle(
				http.HandlerFunc(da.GetDataJSONHandler))))

	//TEXT handlers
	r.Get("/value/{type}/{name}",
		logger.WithLogging(
			compression.GzipHandle(
				http.HandlerFunc(da.GetDataTextHandler))))

	r.Post("/update/{type}/{name}/{value}",
		logger.WithLogging(
			http.HandlerFunc(da.SetDataTextHandler)))

	//Сохраняем состояние оперативного хранилища на диске при выходе из программы-
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-exit
		logger.SLog.Infow("received termination", "signal", sig)
		if da.SyncStorage != nil {
			da.SyncStorage.Truncate()
			logger.Log.Info("saving storage state to disk")
			da.CopyState(st, da.SyncStorage)
		}
		//Закрываем хранилище - актуально для БД
		da.Ds.Close()
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
