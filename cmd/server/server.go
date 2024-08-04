package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/Ssnakerss/practicum-metrics/internal/compression"
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

	endPointAddress := ""
	//переменные окружения
	//ADDRESS отвечает за адрес эндпоинта HTTP-сервера.
	if endPointAddress = os.Getenv("ADDRESS"); endPointAddress == "" {
		//Не нашли переменные окружения
		//Параметры командной строки
		//Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
		ep := flag.String("a", "localhost:8080", "endpoint address")
		flag.Parse()
		endPointAddress = *ep
		logger.SLog.Infow(
			"use CMD or DEFAULT for config",
			"endPointAddress", endPointAddress,
		)
	} else {
		logger.SLog.Infow(
			"use ENV VAR for config",
			"ADDRESS", endPointAddress,
		)
	}
	//Configuring CHI
	r := chi.NewRouter()
	r.Get("/", logger.WithLogging(http.HandlerFunc(handlers.MainPage)))

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

	logger.SLog.Infow(
		"starting server at",
		"addr", endPointAddress,
	)

	err := http.ListenAndServe(endPointAddress, r)
	if err != nil {
		logger.SLog.Fatalf(
			"failed to start server -> program will exit",
			"address", endPointAddress,
			"error", err,
		)
	}

}
