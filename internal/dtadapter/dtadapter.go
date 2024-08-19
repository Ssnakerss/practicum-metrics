package dtadapter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Ssnakerss/practicum-metrics/internal/compression"
	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/Ssnakerss/practicum-metrics/internal/storage"
	"github.com/go-chi/chi/v5"
)

type Adapter struct{ ds storage.DataStorage }

func (da *Adapter) New(ds storage.DataStorage) {
	da.ds = ds
}

func (da *Adapter) MainPage(w http.ResponseWriter, r *http.Request) {
	mcs := make([]metric.Metric, 0)
	_, err := da.ds.ReadAll(&mcs)

	logger.SLog.Debugw("ADAPTER", "metric", mcs)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body := ""
	for _, v := range mcs {
		body += fmt.Sprintf("Name: %s  Type: %s Value: %s \r\n", v.Name, v.Type, v.Value())
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(body))
}

// Handler to save metric received with JSON
func (da *Adapter) SetDataJSONHandler(w http.ResponseWriter, r *http.Request) {
	m, err := da.checkRequestAndGetMetric(r)
	if err != nil {
		logger.SLog.Errorw("fail to receive data", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//Сохраняем метрику в хранилище
	if err = da.ds.Write(m); err != nil {
		logger.SLog.Errorw("fail to save metric", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.SLog.Infow("receive new", "metric", m)

	//Возвращаем метрику из хранилища с обновленным Value
	mj, err := da.readMetricAndMarshal(m)
	if err != nil {
		logger.SLog.Errorw("", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(mj)
	// w.WriteHeader(http.StatusOK)
}

// Handler to save metric received with JSON
func (da *Adapter) GetDataJSONHandler(w http.ResponseWriter, r *http.Request) {
	m, err := da.checkRequestAndGetMetric(r)
	if err != nil {
		logger.SLog.Errorw("fail to receive data", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//Возвращаем метрику из хранилища с обновленным Value
	mj, err := da.readMetricAndMarshal(m)
	if err != nil {
		logger.SLog.Errorw("", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(mj)
	// w.WriteHeader(http.StatusOK)
}

// Read metric and convert to data interface type
func (da *Adapter) readMetricAndMarshal(m *metric.Metric) ([]byte, error) {
	err := da.ds.Read(m)
	if err != nil {
		return nil, fmt.Errorf("fail to read metric: %w", err)
	}
	mi := metric.ConvertMetricS2I(m)
	mj, err := json.Marshal(mi)
	if err != nil {
		return nil, fmt.Errorf("fail to convert saved metric: %w", err)
	}
	return mj, nil
}

// Cheking if request correct and extract metric from Body
func (da *Adapter) checkRequestAndGetMetric(r *http.Request) (*metric.Metric, error) {
	ct := r.Header.Get("content-type")
	if !strings.Contains(ct, "application/json") {
		return nil, fmt.Errorf("incorrect content type: %v", ct)
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read request body: %w", err)
	}
	//Decompression -> TODO: Change to MiddleWare
	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		body, err = compression.Decompress(body)
		if err != nil {
			return nil, fmt.Errorf("fail to un-gzip body %w", err)
		}
	}
	var mi metric.MetricJSON
	err = json.Unmarshal(body, &mi)
	if err != nil {
		return nil, fmt.Errorf("fail to convert json: %w", err)
	}
	return metric.ConvertMetricI2S(&mi), nil
}

// Set metric  via post url
func (da *Adapter) SetDataTextHandler(w http.ResponseWriter, r *http.Request) {
	var m metric.Metric
	var err error
	w.Header().Set("Content-Type", "text/plain")
	//Make metric params case insensitive
	mType := strings.ToLower(chi.URLParam(r, "type"))
	mName := strings.ToLower(chi.URLParam(r, "name"))
	mValue := strings.ToLower(chi.URLParam(r, "value"))
	//Checking metric type and name
	if err = m.Set(mName, mValue, mType); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	//Processing metrics values
	if err = da.ds.Write(&m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Get request for metrci via URL
func (da *Adapter) GetDataTextHandler(w http.ResponseWriter, r *http.Request) {
	var m metric.Metric
	w.Header().Set("Content-Type", "text/plain")
	//Make metric params case insensitive
	mType := strings.ToLower(chi.URLParam(r, "type"))
	mName := strings.ToLower(chi.URLParam(r, "name"))
	if err := m.Set(mName, "0", mType); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//Selecting metrics from storage
	if err := da.ds.Read(&m); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write([]byte(m.Value()))
	w.WriteHeader(http.StatusOK)
}
