package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type procF func(metricType string, metricName string, value float64, mem MemStorage) bool

type metricsValues []float64

type MemStorage struct {
	metrics map[string]metricsValues
	addItem procF
}

var methods map[string]procF
var storage MemStorage

func newMetricDataProcessing(metricType string, metricName string, value float64, mem MemStorage) bool {
	switch metricType {
	case "gauge":
		if mem.metrics[metricName] == nil {
			mem.metrics[metricName] = make([]float64, 1)
		}
		mem.metrics[metricName][0] = value
		return true
	case "counter":
		if mem.metrics[metricName] == nil {
			mem.metrics[metricName] = make([]float64, 0)
		}
		mem.metrics[metricName] = append(mem.metrics[metricName], value)
		return true
	default:
		return false
	}
}

func main() {
	//Initialize used methods
	methods = make(map[string]procF)
	methods["gauge"] = newMetricDataProcessing
	methods["counter"] = newMetricDataProcessing

	// vals := make([]float64)
	storage = MemStorage{make(map[string]metricsValues), newMetricDataProcessing}

	mux := http.NewServeMux()
	//diagnistics
	mux.HandleFunc(`/`, updateHandler)
	//-----------
	mux.HandleFunc(`/update/`, updateHandler)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	if r.Method != http.MethodPost {
		//diagnostic
		body := ""
		body += fmt.Sprintf("MemStorage == %v\r\n", storage)
		body += fmt.Sprintf("MemStorage == %#v\r\n", storage)
		w.Write([]byte(body))
		//-----------

		// w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	params := strings.Split(r.URL.Path, "/")
	if len(params) != 5 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		if p := methods[params[2]]; p == nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			//Processing metrics values
			val, err := strconv.ParseFloat(params[4], 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			storage.addItem(params[2], params[3], val, storage)
			w.WriteHeader(http.StatusOK)
		}
	}
	//Diagnostics
	// body := ""
	// for idx, p := range params {
	// 	body += fmt.Sprintf("param %d - %s\r\n", idx, p)
	// }
	// body += fmt.Sprintf("Method == %s\r\n", r.Method)
	// body += fmt.Sprintf("MemStorage == %v\r\n", storage)
	// body += fmt.Sprintf("MemStorage == %#v\r\n", storage)

	// w.Write([]byte(body))

}

// type Middleware func(http.Handler) http.Handler
// func Conveyor(h http.Handler, middlewares ...Middleware) http.Handler {
// 	for _, middleware := range middlewares {
// 		h = middleware(h)
// 	}
// 	return h
// }

// func checkMethod(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 		return
// 	}
// }
