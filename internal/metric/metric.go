package metric

import (
	"fmt"
	"strconv"
	"strings"
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

	if mi.Value != nil || mi.Delta != nil {

		switch mi.MType {
		case "gauge":
			ms.Gauge = *mi.Value
		case "counter":
			ms.Counter = *mi.Delta
		}
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
func (m *Metric) Copy() *Metric {
	cpm := Metric{
		Name:    m.Name,
		Type:    m.Type,
		Counter: m.Counter,
		Gauge:   m.Gauge,
	}
	return &cpm
}

func (m *Metric) Clear() {
	m.Name = ""
	m.Type = ""
	m.Gauge = 0
	m.Counter = 0
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
