package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/Ssnakerss/practicum-metrics/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {

	// •	Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
	endPointAddress := flag.String("a", `localhost:8080`, "endpoint address")

	flag.Parse()

	fmt.Printf("server started at %s\r\n", *endPointAddress)

	r := chi.NewRouter()
	r.Get("/", handlers.MainPage)
	r.Get("/value/{type}/{name}", handlers.ChiGetHandler)
	r.Post("/update/{type}/{name}/{value}", handlers.ChiUpdateHandler)

	err := http.ListenAndServe(*endPointAddress, r)
	if err != nil {
		panic(err)
	}
}
