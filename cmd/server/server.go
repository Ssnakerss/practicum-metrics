package main

import (
	"fmt"
	"net/http"

	"github.com/Ssnakerss/practicum-metrics/internal/handlers"
)

func main() {
	fmt.Println("server started ...")
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, handlers.UpdateHandler)
	//mux.HandleFunc(`/`, mainPage)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
