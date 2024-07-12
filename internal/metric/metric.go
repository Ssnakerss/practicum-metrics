package metric

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
)

var MemStatsMetrics = [27]string{
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

type Metric struct {
	//metric name - Alloc,	BuckHashSys etc
	Name, Type string
	Gauge      float64
	Counter    int64
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
