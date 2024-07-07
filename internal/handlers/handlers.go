package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/Ssnakerss/practicum-metrics/internal/storage"
	"github.com/go-chi/chi/v5"
)

var Stor storage.Storage

func init() {
	Stor.New()
	fmt.Println("Initialize storage ....")
}

func ChiUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var m metric.Metric

	w.Header().Set("Content-Type", "text/plain")
	//Make metric params case insensitive
	mType := strings.ToLower(chi.URLParam(r, "type"))
	mName := strings.ToLower(chi.URLParam(r, "name"))
	mValue := strings.ToLower(chi.URLParam(r, "value"))

	//Checking metric type and name
	if ok, _ := m.Set(mName, mValue, "float64", mType); ok {
		//Processing metrics values
		switch m.MType {
		case "gauge":
			err := Stor.Update(m)
			if err != nil {
				//panic("UPDATE ERROR")
				fmt.Println(err)
			}
		case "counter":
			err := Stor.Insert(m)
			if err != nil {
				// panic("INSERT ERROR")
				fmt.Println(err)
			}
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func ChiGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	//Make metric params case insensitive
	mType := strings.ToLower(chi.URLParam(r, "type"))
	mName := strings.ToLower(chi.URLParam(r, "name"))

	if metric.IsAllowed(mName, mType) {
		w.Write([]byte(prepareBody(mName + "@" + mType)))
		// w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(prepareBody("")))
	// w.WriteHeader(http.StatusOK)
}

func prepareBody(s string) string {
	results := make(map[string]metric.Metric)
	var body string

	if s == "" {
		Stor.Select(results)
		for _, v := range results {
			body += fmt.Sprintf("Name: %s Value: %v \r\n", v.Name, v.Value)
		}
	} else {
		Stor.Select(results, s)
		for _, v := range results {
			if len(v.Value) > 0 {
				// body += fmt.Sprintf("%s", v.Value[len(v.Value)-1])
				body += v.Value[len(v.Value)-1]
			}
		}
	}

	return body
}
