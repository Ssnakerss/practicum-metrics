package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/Ssnakerss/practicum-metrics/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	var endPointAddress string = ""
	//переменные окружения
	//•	ADDRESS отвечает за адрес эндпоинта HTTP-сервера.

	if endPointAddress = os.Getenv("ADDRESS"); endPointAddress == "" {
		//Параметры командной строки
		// •	Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
		ep := flag.String("a", "localhost:8080", "endpoint address")
		flag.Parse()
		endPointAddress = *ep
	}

	fmt.Printf("server started at %s\r\n", endPointAddress)

	r := chi.NewRouter()
	r.Get("/", handlers.MainPage)
	r.Get("/value/{type}/{name}", handlers.ChiGetHandler)
	r.Post("/update/{type}/{name}/{value}", handlers.ChiUpdateHandler)

	err := http.ListenAndServe(endPointAddress, r)
	if err != nil {
		panic(err)
	}
}
