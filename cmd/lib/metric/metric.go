package metric

import (
	"fmt"
	"strconv"
)

var MetricsToGather = map[string]bool{
	"Frees":         true,
	"Alloc":         true,
	"BuckHashSys":   true,
	"GCCPUFraction": true,
	"GCSys":         true,
	"HeapAlloc":     true,
	"HeapIdle":      true,
	"HeapInuse":     true,
	"HeapObjects":   true,
	"HeapReleased":  true,
	"HeapSys":       true,
	"LastGC":        true,
	"Lookups":       true,
	"MCacheInuse":   true,
	"MCacheSys":     true,
	"MSpanInuse":    true,
	"MSpanSys":      true,
	"Mallocs":       true,
	"NextGC":        true,
	"NumForcedGC":   true,
	"NumGC":         true,
	"OtherSys":      true,
	"PauseTotalNs":  true,
	"StackInuse":    true,
	"StackSys":      true,
	"Sys":           true,
	"TotalAlloc":    true,
	"testCounter":   true,
	"testGauge":     true,
	"PollCount":     true,
	"RandomValue":   true,
}

type Metric struct {
	//metric name - Alloc,	BuckHashSys etc
	Name string
	//metric value : float64  to to store float64, uint64 uint 32
	value float64
	//metric value type : to proper convert from  float64  to  uint64 uint 32 when necessary
	vType string
	//metric type - gauge, counter : use when saving to DB
	MType string
}

// check metric name and type by allowed values
func (m *Metric) IsValid(name string, mType string) bool {
	switch mType {
	case "gauge", "counter":
		return MetricsToGather[name]
	default:
		return false
	}
}

// convert metric value to float64 to keep
func (m *Metric) convertValue(value string, vType string) (float64, error) {
	switch vType {
	case "float32", "float64":
		if v, err := strconv.ParseFloat(value, 64); err != nil {
			return v, nil
		}
	case "uint32", "uint64":
		if v, err := strconv.ParseUint(value, 10, 64); err != nil {
			return float64(v), nil
		}
	default:
		return 0, fmt.Errorf("unknown value type: %s -> %s", value, vType)
	}
	return 0, fmt.Errorf("type convertion error: %s -> %s", value, vType)
}

// set metric values
func (m *Metric) Set(name string, value string, vType string, mType string) (bool, error) {
	if m.IsValid(name, mType) {
		if v, err := m.convertValue(value, vType); err != nil {
			m.Name, m.value, m.vType, m.MType = name, v, vType, mType
			return true, nil
		}
		return false, fmt.Errorf("type convertion error: %s -> %s", value, vType)
	}
	return false, fmt.Errorf("invalid metric name %s", name)
}
