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

// First attempt .....  {params} ???   who knows
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var m metric.Metric

	w.Header().Set("Content-Type", "text/plain")
	//Only POST method allowed
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	//Parsing URL to get parameters
	params := strings.Split(r.URL.Path, "/")
	if len(params) != 5 {
		//URL has extra or less parameters
		w.WriteHeader(http.StatusNotFound)
	} else {
		//Checking metric type and name
		if ok, _ := m.Set(params[3], params[4], "float64", params[2]); ok {
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
}

func ChiUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var m metric.Metric

	w.Header().Set("Content-Type", "text/plain")
	mType := chi.URLParam(r, "type")
	mName := chi.URLParam(r, "name")
	mValue := chi.URLParam(r, "value")

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
	mType := chi.URLParam(r, "type")
	mName := chi.URLParam(r, "name")
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
	if s == "" {
		Stor.Select(results)
	} else {
		Stor.Select(results, s)
	}

	var body string
	for _, v := range results {
		body += fmt.Sprintf("Name: %s Value: %v \r\n", v.Name, v.Value)
	}
	return body
}
