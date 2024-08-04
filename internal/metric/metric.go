package metric

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
)

var MemStatsMetrics = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

type etcMetrics struct {
	MType string
	MFunc func(p ...int) string
}

var ExtraMetrics = map[string]etcMetrics{
	"PollCount": {
		MType: "counter",
		MFunc: func(p ...int) string {
			return strconv.Itoa(p[0])
		},
	},
	"RandomValue": {
		MType: "gauge",
		MFunc: func(p ...int) string {
			return strconv.FormatFloat(rand.Float64(), 'f', -1, 64)
		},
	},
}

type (
	MetricJSON struct {
		ID    string   `json:"id"`              // имя метрики
		MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
		Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
		//    ↑ это видимо для возможности делать EMPTY!
		Delta *int64 `json:"delta,omitempty"` // значение метрики в случае передачи counter
	}

	Metric struct { // Оставляем для обратной совместимости
		Name    string  `json:"id"`              //ID
		Type    string  `json:"type"`            //MType
		Gauge   float64 `json:"value,omitempty"` //Value
		Counter int64   `json:"delta,omitempty"` //Delta
	}
)

// Convert metric from [S]torage format to [I]nterface format
func ConvertMetricS2I(ms *Metric) *MetricJSON {
	mi := MetricJSON{
		ID:    ms.Name,
		MType: ms.Type,
	}
	switch ms.Type {
	case "gauge":
		mi.Value = &ms.Gauge
	case "counter":
		mi.Delta = &ms.Counter
	}
	return &mi
}

// Convert metric from [S]torage format to [I]nterface format
func ConvertMetricI2S(mi *MetricJSON) *Metric {
	ms := Metric{
		Name: mi.ID,
		Type: mi.MType,
	}
	switch mi.MType {
	case "gauge":
		ms.Gauge = *mi.Value
	case "counter":
		ms.Counter = *mi.Delta
	}
	return &ms
}

// IsValid - Check metric name and type by allowed values
func IsValid(mType string, mValue string) bool {
	mType = strings.ToLower(mType)
	switch mType {
	case "gauge":
		if _, err := strconv.ParseFloat(mValue, 64); err == nil {
			return true
		}
	case "counter":
		if _, err := strconv.ParseInt(mValue, 10, 64); err == nil {
			return true
		}
	}
	return false
}

// Value -  current metric value regardless type
func (m *Metric) Value() string {
	switch m.Type {
	case "gauge":
		return strconv.FormatFloat(m.Gauge, 'f', -1, 64)
	case "counter":
		return strconv.FormatInt(m.Counter, 10)
	}
	return ""
}

// Set metric values
func (m *Metric) Set(
	mName, mValue, mType string,
) error {
	if !IsValid(mType, mValue) {
		return fmt.Errorf("invalid metric type or value: %s, %s", mType, mValue)
	}
	m.Type = mType
	m.Name = mName

	var err error
	switch mType {
	case "gauge":
		if m.Gauge, err = strconv.ParseFloat(mValue, 64); err != nil {
			return err
		}
	case "counter":
		if m.Counter, err = strconv.ParseInt(mValue, 0, 64); err != nil {
			return err
		}
	}
	return nil
}
