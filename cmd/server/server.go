package main

import (
	"fmt"
	"net/http"

	"github.com/Ssnakerss/practicum-metrics/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	fmt.Println("server started ...")
	r := chi.NewRouter()

	r.Get("/", handlers.MainPage)
	r.Get("/value/{type}/{name}", handlers.ChiGetHandler)

	r.Post("/update/{type}/{name}/{value}", handlers.ChiUpdateHandler)

	err := http.ListenAndServe(`:8080`, r)

	// mux := http.NewServeMux()
	// mux.HandleFunc(`/update/`, handlers.UpdateHandler)
	// mux.HandleFunc(`/`, handlers.MainPage)
	// err := http.ListenAndServe(`:8080`, mux)

	if err != nil {
		panic(err)
	}
}
