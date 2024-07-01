package main

import (
	"net/http"
)

var Stor Storage

func main() {
	Stor.New()
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, updateHandler)
	//mux.HandleFunc(`/`, mainPage)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
