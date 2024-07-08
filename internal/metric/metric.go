package metric

import (
	"fmt"
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

var ExtraMetrics = map[string]bool{
	"PollCount":   true,
	"RandomValue": true,
}

type Metric struct {
	//metric name - Alloc,	BuckHashSys etc
	Name    string
	Type    string
	Gauge   float64
	Counter int64
}

// IsValid - Check metric name and type by allowed values
func IsValid(mType string, mValue string) bool {
	mType = strings.ToLower(mType)
	switch mType {
	case "gauge", "counter":
		if _, err := strconv.ParseFloat(mValue, 64); err == nil {
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
func (m *Metric) Set(mName string, mValue string, mType string) (bool, error) {
	if !IsValid(mType, mValue) {
		return false, fmt.Errorf("invalid metric type or value: %s, %s", mType, mValue)
	}
	m.Type = mType
	m.Name = mName

	var err error = nil
	switch mType {
	case "gauge":
		if m.Gauge, err = strconv.ParseFloat(mValue, 64); err != nil {
			return false, err
		}
	case "counter":
		if m.Counter, err = strconv.ParseInt(mValue, 0, 64); err != nil {
			return false, err
		}
	}
	return true, nil
}
