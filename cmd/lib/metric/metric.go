package metric

import (
	"fmt"
	"strconv"
	"strings"
)

var AllowedMetrics = map[string]bool{
	"alloc":         true,
	"buckhashsys":   true,
	"frees":         true,
	"gccpufraction": true,
	"gcsys":         true,
	"heapalloc":     true,
	"heapidle":      true,
	"heapinuse":     true,
	"heapobjects":   true,
	"heapreleased":  true,
	"heapsys":       true,
	"lastgc":        true,
	"lookups":       true,
	"mcacheinuse":   true,
	"mcachesys":     true,
	"mspaninuse":    true,
	"mspansys":      true,
	"mallocs":       true,
	"nextgc":        true,
	"numforcedgc":   true,
	"numgc":         true,
	"othersys":      true,
	"pausetotalns":  true,
	"stackinuse":    true,
	"stacksys":      true,
	"sys":           true,
	"totalalloc":    true,

	"pollcount":   true,
	"randomvalue": true,

	"testcounter": true,
	"testgauge":   true,
}

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
	Name string
	//metric value : float64  to to store float64, uint64 uint 32
	Value float64
	//metric value type : to proper convert from  float64  to  uint64 uint 32 when necessary
	VType string
	//metric type - gauge, counter : use when saving to DB
	MType string
}

// IsValid - Check metric name and type by allowed values
func (m *Metric) IsValid(name string, mType string) bool {
	name = strings.ToLower(name)
	mType = strings.ToLower(mType)
	switch mType {
	case "gauge", "counter":
		return (AllowedMetrics[name])
	default:
		return false
	}
}

// convertValue - Convert metric value to float64 to keep
func (m *Metric) convertValue(value string, vType string) (float64, error) {
	switch vType {
	case "float32", "float64":
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			return v, nil
		} else {
			fmt.Println(err)
		}

	case "uint32", "uint64":
		if v, err := strconv.ParseUint(value, 10, 64); err == nil {
			return float64(v), nil
		} else {
			fmt.Println(err)
		}

	default:
		return 0, fmt.Errorf("unknown value type: %s -> %s", value, vType)
	}
	fmt.Printf("error converting %s TO %s\n\r", value, vType)
	return 0, fmt.Errorf("type convertion error: %s -> %s", value, vType)
}

// Set metric values
func (m *Metric) Set(name string, value string, vType string, mType string) (bool, error) {
	if !m.IsValid(name, mType) {
		return false, fmt.Errorf("invalid name or type: %s, %s", name, mType)
	}

	if v, err := m.convertValue(value, vType); err == nil {
		m.Name, m.Value, m.VType, m.MType = name, v, vType, mType
		return true, nil
	}
	return false, fmt.Errorf("type convertion error: %s -> %s", value, vType)

}
