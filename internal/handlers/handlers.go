package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Ssnakerss/practicum-metrics/internal/compression"
	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/Ssnakerss/practicum-metrics/internal/storage"
	"github.com/go-chi/chi/v5"
)

var Stor storage.Storage

func InitStorage(filePath string, syncWrite bool) {
	Stor.New(filePath, syncWrite)
	log.Println("Initialize storage ....")
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	results := make(map[string]*metric.Metric)
	var body string
	//Return all values
	Stor.Select(results)
	for _, v := range results {
		body += fmt.Sprintf("Name: %s  Type: %s Value: %s \r\n", v.Name, v.Type, v.Value())
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(body))
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
		if err = Stor.SaveMetric(&m); err == nil {
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
}

// Получаем и обрабатываем метрику в URL params
func GetDataTextHandler(w http.ResponseWriter, r *http.Request) {
	var m metric.Metric
	w.Header().Set("Content-Type", "text/plain")
	//Make metric params case insensitive
	mType := strings.ToLower(chi.URLParam(r, "type"))
	mName := strings.ToLower(chi.URLParam(r, "name"))
	if err := m.Set(mName, "0", mType); err == nil {
		//Selecting metrics from storage
		results := make(map[string]*metric.Metric)
		//--------------------------------------
		Stor.Select(results, &m)
		if m, ok := results[m.Name+m.Type]; ok {
			if m.Name != "" {
				w.Write([]byte(m.Value()))
				return
			}
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

// Получаем и обрабатываем метрику в JSON
func SetDataJSONHandler(w http.ResponseWriter, r *http.Request) {
	if m, statusCode, err := checkRequestAndGetMetric(w, r, "setdata"); err == nil {
		//Все ОК - сохраняем метрику
		if err := Stor.SaveMetric(m); err != nil {
			logger.SLog.Infow("fail to process", "metric", m)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.SLog.Infow("aquire new", "metric", m)

		//Надо вернуть метрику с обновленным значением Value
		//Выбираем метрику из хранилища
		results := make(map[string]*metric.Metric)
		Stor.Select(results, m)

		if nm, ok := results[m.Name+m.Type]; ok {
			mj := metric.ConvertMetricS2I(nm)
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
	} else {
		http.Error(w, "error", statusCode)
	}
}

func GetDataJSONHandler(w http.ResponseWriter, r *http.Request) {
	m, err := checkRequestAndGetMetric(w, r, "getdata")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//Надо вернуть метрику с обновленным значением Value
	//Выбираем метрику из хранилища
	results := make(map[string]*metric.Metric)
	if found := Stor.Select(results, m); found > 0 {
		if nm, ok := results[m.Name+m.Type]; ok {
			mj := metric.ConvertMetricS2I(nm)
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
	}
	http.Error(w, "metric not found", http.StatusNotFound)
}

// func readMetric(m *metric.Metric) ([]byte, int, error) {
// 	results := make(map[string]*metric.Metric)
// 	if found := Stor.Select(results, m); found > 0 {
// 		if nm, ok := results[m.Name+m.Type]; ok {
// 			mj := metric.ConvertMetricS2I(nm)
// 			if jsonMetric, err := json.Marshal(mj); err != nil {
// 				logger.Log.Error("something went wrong")
// 				return nil, http.StatusInternalServerError,
// 					fmt.Errorf("cannot marshal metric %v", mj)
// 			} else {
// 				return jsonMetric, http.StatusOK, nil
// 			}
// 		}
// 	}
// 	return nil, http.StatusNotFound, fmt.Errorf("metric not found")
// }

// ------------------------------
func checkRequestAndGetMetric(w http.ResponseWriter,
	r *http.Request, rtype string) (*metric.Metric, int, error) {
	ct := r.Header.Get("content-type")
	logger.SLog.Infoln(">> got request ", "request type", rtype, "content-type", ct)

	if ct != "application/json" {
		logger.SLog.Errorw("incorrect content-type:", "content-type", ct)
		return nil, http.StatusBadRequest, fmt.Errorf("incorrect content type")
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.SLog.Errorw("error reading body:", "body", r.Body)
		return nil, http.StatusBadRequest, err
	}

	//Decompression
	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		body, err = compression.Decompress(body)
		if err != nil {
			logger.SLog.Errorw("failed gzip decompression", "error", err)
			return nil, http.StatusBadRequest, err
		}
		logger.SLog.Infow("gzip content decompressed", "body", body)
	}
	//

	var mi metric.MetricJSON
	err = json.Unmarshal(body, &mi)
	if err != nil {
		logger.SLog.Errorw("fail to unmarshall", "body", string(body))
		return nil, http.StatusBadRequest, fmt.Errorf("fail to unmarshall json")
	}
	m := metric.ConvertMetricI2S(&mi)
	return m, http.StatusOK, nil
}
