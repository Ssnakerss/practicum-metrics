package dtadapter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/go-chi/chi/v5"
)

// Возвращаем список метрик
func (da *Adapter) MainPage(w http.ResponseWriter, r *http.Request) {
	mcs := make([]metric.Metric, 0)
	err := da.ReadAll(&mcs)

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

// Handler для получения и сохранения массива метрик
func (da *Adapter) SetDataJSONSliceHandler(w http.ResponseWriter, r *http.Request) {
	//Проверяем запрос и распаковываем body
	ct := r.Header.Get("content-type")
	if !strings.Contains(ct, "application/json") {
		logger.SLog.Errorf("bad request", "content-type", ct)
		http.Error(w, "incorrect content type", http.StatusBadRequest)
		return
	}
	//Распаковка боди уже в мидлваре
	//Теперь у нас есть распакованная строка, пробуем сконвертить ее в массив метрик
	mcsj := make([]metric.MetricJSON, 0)
	mcs := make([]metric.Metric, 0)
	if err := json.NewDecoder(r.Body).Decode(&mcsj); err != nil {
		logger.SLog.Warn("cannot convert json to []metric, check rsa")
		http.Error(w, "cannot convert json to []metric", http.StatusBadRequest)
		return
	}
	for _, m := range mcsj {
		mcs = append(mcs, *metric.ConvertMetricI2S(&m))
	}
	//Записываем получившийся массив в хранилище
	err := da.WriteAll(&mcs)
	if err != nil {
		logger.SLog.Error("error saving to storage")
		http.Error(w, "error saving to storage", http.StatusInternalServerError)
		return
	}
	logger.SLog.Infow("received new [] of metrics", "count", len(mcs))
	w.WriteHeader(http.StatusOK)
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
	if err = da.Write(m); err != nil {
		logger.SLog.Errorw("fail to save metric", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.SLog.Infow("receive new", "metric", m)

	//Возвращаем метрику из хранилища с обновленным Value
	mj, err := da.readMetricAndMarshal(m)
	if err != nil {
		logger.SLog.Warnw("SetDataJSONHandler", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(mj)

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
		logger.SLog.Warnw("GetDataJSONHandler", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(mj)

}

// Set metric  via post url
func (da *Adapter) SetDataTextHandler(w http.ResponseWriter, r *http.Request) {
	var m metric.Metric
	var err error

	mType := (chi.URLParam(r, "type"))
	mName := (chi.URLParam(r, "name"))
	mValue := (chi.URLParam(r, "value"))

	//Checking metric type and name
	if err = m.Set(mName, mValue, mType); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//Processing metrics values
	if err = da.Write(&m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

// Get request for metrci via URL
func (da *Adapter) GetDataTextHandler(w http.ResponseWriter, r *http.Request) {
	var m metric.Metric

	mType := (chi.URLParam(r, "type"))
	mName := (chi.URLParam(r, "name"))

	if err := m.Set(mName, "0", mType); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	//Selecting metrics from storage
	if err := da.Read(&m); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(m.Value()))
}

// Проверяем соединение с базой данных
func (da *Adapter) Ping(w http.ResponseWriter, r *http.Request) {

	//Не важно, какой у нас тип хранилища используется
	//CheckStorage вернет nil для memory & file
	//если используется db -  проверит состояние подключения

	if err := da.Ds.CheckStorage(); err != nil {
		logger.SLog.Warnw("db storage connection check", "fail", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
