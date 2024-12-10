package metric

import (
	"math/rand/v2"
	"strconv"
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

var GopsMetrics = []string{
	"TotalMemory",
	"FreeMemery",
	"CPUutilization1",
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
		Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	}

	Metric struct { // Оставляем для обратной совместимости
		Name    string  `json:"name"`              //ID
		Type    string  `json:"type"`              //MType
		Gauge   float64 `json:"gauge,omitempty"`   //Value
		Counter int64   `json:"counter,omitempty"` //Delta
	}
)
