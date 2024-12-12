package server

import (
	"net/http"

	"github.com/Ssnakerss/practicum-metrics/internal/app"
	"github.com/Ssnakerss/practicum-metrics/internal/compression"
	"github.com/Ssnakerss/practicum-metrics/internal/dtadapter"
	"github.com/Ssnakerss/practicum-metrics/internal/hash"
	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "net/http/pprof"
)

func NewRouter(da *dtadapter.Adapter, c *app.ServerConfig) *chi.Mux {

	//Configuring CHI
	r := chi.NewRouter()

	r.Use(logger.WithLogging)
	r.Use(compression.GzipHandle)

	h := hash.New(c.Key)
	r.Use(h.Handle)

	//Добваляем обработчики для pprof
	r.Mount("/debug", middleware.Profiler())

	r.Get("/", http.HandlerFunc(da.MainPage))

	//Проверка соединения с хранилищем
	r.Get("/ping", http.HandlerFunc(da.Ping))

	// JSON handlers
	r.Post("/update/", http.HandlerFunc(da.SetDataJSONHandler))
	r.Post("/updates/", http.HandlerFunc(da.SetDataJSONSliceHandler))

	r.Post("/value/", http.HandlerFunc(da.GetDataJSONHandler))

	//TEXT handlers
	r.Get("/value/{type}/{name}", http.HandlerFunc(da.GetDataTextHandler))
	r.Post("/update/{type}/{name}/{value}", http.HandlerFunc(da.SetDataTextHandler))

	return r
}
