package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/Ssnakerss/practicum-metrics/internal/storage"
	"github.com/go-chi/chi/v5"
)

var Stor storage.Storage

func init() {
	Stor.New()
	log.Println("Initialize storage ....")
}

func SetDataTextHandler(w http.ResponseWriter, r *http.Request) {
	var m metric.Metric
	w.Header().Set("Content-Type", "text/plain")
	//Make metric params case insensitive
	mType := strings.ToLower(chi.URLParam(r, "type"))
	mName := strings.ToLower(chi.URLParam(r, "name"))
	mValue := strings.ToLower(chi.URLParam(r, "value"))
	//Checking metric type and name
	if err := m.Set(mName, mValue, mType); err == nil {
		//Processing metrics values
		if err = storage.ProcessMetric(m, &Stor); err == nil {
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
}

// Получаем и обрабатываем метрику в JSON
func SetDataJSONHandler(w http.ResponseWriter, r *http.Request) {
	if m, err := checkRequestAndGetMetric(w, r, "setdata"); err == nil {
		//Все ОК - сохраняем метрику
		//Все отличие отсюда ↓
		if err := storage.ProcessMetric(*m, &Stor); err != nil {
			logger.SLog.Infow("fail to process", "metric", m)

			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.SLog.Infow("aquire new", "metric", m)
		//И до сюда ↑

		//Надо вернуть метрику с обновленным значением Value
		//Выбираем метрику из хранилища
		results := make(map[string]metric.Metric)
		//--------------------------------------
		Stor.Select(results, m.Name)
		if nm, ok := results[m.Name]; ok {
			mj := metric.ConvertMetricS2I(&nm)
			b, err := json.Marshal(mj)
			if err != nil {
				logger.Log.Error("something went wrong")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(b)
			return
		}
		http.Error(w, "cannot retrieve metric from storage", http.StatusInternalServerError)
	}
}

// Получаем и обрабатываем метрику в URL params
func GetDataTextHandler(w http.ResponseWriter, r *http.Request) {
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

func GetDataJSONHandler(w http.ResponseWriter, r *http.Request) {

	if m, err := checkRequestAndGetMetric(w, r, "getdata"); err == nil {
		//Надо вернуть метрику с обновл↑енным значением Value
		//Выбираем метрику из хранилища
		results := make(map[string]metric.Metric)
		//--------------------------------------
		Stor.Select(results, m.Name)
		if nm, ok := results[m.Name]; ok {
			mj := metric.ConvertMetricS2I(&nm)
			b, err := json.Marshal(mj)
			if err != nil {
				logger.Log.Error("something went wrong")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(b)
			return
		}
		http.Error(w, "cannot retrieve metric from storage", http.StatusInternalServerError)
	}
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

// ------------------------------
func checkRequestAndGetMetric(w http.ResponseWriter,
	r *http.Request, rtype string) (*metric.Metric, error) {
	ct := r.Header.Get("content-type")
	logger.SLog.Infoln(">> got request ", "request type", rtype, "content-type", ct)

	if ct != "application/json" {
		logger.SLog.Errorw("incorrect content-type:", "content-type", ct)

		http.Error(w, "incorrect content type", http.StatusBadRequest)
		return nil, fmt.Errorf("bad request")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.SLog.Errorw("error reading body:", "body", r.Body)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	var mi metric.MetricJSON
	err = json.Unmarshal(body, &mi)
	if err != nil {
		logger.SLog.Errorw("fail to unmarshall", "body", string(body))

		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, fmt.Errorf("fail to unmarshall")
	}
	m := metric.ConvertMetricI2S(&mi)
	return m, nil
}
