package main

import (
	"net/http"
	"strconv"
	"strings"
)

type procF func(metricType string, metricName string, value float64, mem *MemStorage) bool

type metricsValues []float64

type MemStorage struct {
	metrics map[string]metricsValues
	addItem procF
}

var methods map[string]procF
var storage MemStorage
var metricsAllowed map[string]bool

func newMetricDataProcessing(metricType string, metricName string, value float64, mem *MemStorage) bool {
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

func initializeMetrics() {
	metricsAllowed = make(map[string]bool)
	metricsAllowed["Frees"] = true
	metricsAllowed["Alloc"] = true
	metricsAllowed["BuckHashSys"] = true
	metricsAllowed["GCCPUFraction"] = true
	metricsAllowed["GCSys"] = true
	metricsAllowed["HeapAlloc"] = true
	metricsAllowed["HeapIdle"] = true
	metricsAllowed["HeapInuse"] = true
	metricsAllowed["HeapObjects"] = true
	metricsAllowed["HeapReleased"] = true
	metricsAllowed["HeapSys"] = true
	metricsAllowed["LastGC"] = true
	metricsAllowed["Lookups"] = true
	metricsAllowed["MCacheInuse"] = true
	metricsAllowed["MCacheSys"] = true
	metricsAllowed["MSpanInuse"] = true
	metricsAllowed["MSpanSys"] = true
	metricsAllowed["Mallocs"] = true
	metricsAllowed["NextGC"] = true
	metricsAllowed["NumForcedGC"] = true
	metricsAllowed["NumGC"] = true
	metricsAllowed["OtherSys"] = true
	metricsAllowed["PauseTotalNs"] = true
	metricsAllowed["StackInuse"] = true
	metricsAllowed["StackSys"] = true
	metricsAllowed["Sys"] = true
	metricsAllowed["TotalAlloc"] = true
	metricsAllowed["testCounter"] = false
	metricsAllowed["testGauge"] = true
	metricsAllowed["PollCount"] = true
	metricsAllowed["RandomValue"] = true
}

func main() {
	//fmt.Println("Server started ... ")
	//Initialize
	methods = make(map[string]procF)
	methods["gauge"] = newMetricDataProcessing
	methods["counter"] = newMetricDataProcessing
	initializeMetrics()
	storage = MemStorage{make(map[string]metricsValues), newMetricDataProcessing}

	mux := http.NewServeMux()
	//diagnistics
	//mux.HandleFunc(`/`, updateHandler)
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
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	params := strings.Split(r.URL.Path, "/")
	if len(params) != 5 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		if methods[params[2]] == nil {
			w.WriteHeader(http.StatusBadRequest)
		} else if !metricsAllowed[params[3]] {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			//Processing metrics values
			val, err := strconv.ParseFloat(params[4], 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			storage.addItem(params[2], params[3], val, &storage)
			w.WriteHeader(http.StatusOK)
		}
	}
}
