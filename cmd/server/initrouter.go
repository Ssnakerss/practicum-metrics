package main

import (
	"net/http"

	"github.com/Ssnakerss/practicum-metrics/internal/compression"
	"github.com/Ssnakerss/practicum-metrics/internal/dtadapter"
	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/go-chi/chi/v5"
)

func NewRouter(da *dtadapter.Adapter) *chi.Mux {

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

	return r
}
