package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/Ssnakerss/practicum-metrics/internal/storage"
	"github.com/go-chi/chi/v5"
)

var Stor storage.Storage

func init() {
	Stor.New()
	log.Println("Initialize storage ....")
}

func ChiUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var m metric.Metric
	w.Header().Set("Content-Type", "text/plain")
	//Make metric params case insensitive
	mType := strings.ToLower(chi.URLParam(r, "type"))
	mName := strings.ToLower(chi.URLParam(r, "name"))
	mValue := strings.ToLower(chi.URLParam(r, "value"))
	//Checking metric type and name
	if err := m.Set(mName, mValue, mType); err == nil {
		//Processing metrics values
		switch mType {
		case "gauge":
			err := Stor.Update(m)
			if err != nil {
				log.Printf("error update metric: %s\r\n", err)
				log.Printf("metric value: %v\r\n", m)
			}
		case "counter":
			err := Stor.Insert(m)
			if err != nil {
				log.Printf("error insert metric: %s\r\n", err)
				log.Printf("metric value: %v\r\n", m)
			}
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func ChiGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	//Make metric params case insensitive
	mType := strings.ToLower(chi.URLParam(r, "type"))
	mName := strings.ToLower(chi.URLParam(r, "name"))

	if metric.IsValid(mType, "0") {
		//Selecting metrics from storage
		results := make(map[string]metric.Metric)
		//--------------------------------------
		Stor.Select(results, mName)
		if m, ok := results[mName]; ok {
			//Initialized with default values - "","",0,0
			if m.Name != "" {
				w.Write([]byte(m.Value()))
				return
			}
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	results := make(map[string]metric.Metric)
	var body string
	//Return all values
	Stor.Select(results)
	for _, v := range results {
		body += fmt.Sprintf("Name: %s  Type: %s Value: %s \r\n", v.Name, v.Type, v.Value())
	}
	w.Write([]byte(body))

}
