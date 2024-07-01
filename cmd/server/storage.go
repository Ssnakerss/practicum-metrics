package main

import (
	"fmt"
	"lib/metric"
)

type Storage struct {
	metrics map[string][]metric.Metric
}

func (st *Storage) New() {
	st.metrics = make(map[string][]metric.Metric)
}

func (st *Storage) Insert(m metric.Metric) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	st.metrics[m.Name] = append(st.metrics[m.Name], m)
	return nil
}

func (st *Storage) Update(m metric.Metric) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	if len(st.metrics[m.Name]) == 0 {
		st.metrics[m.Name] = append(st.metrics[m.Name], m)
	} else {
		st.metrics[m.Name][0] = m
	}

	return nil
}

// func (mem MemStorage) addItem(metricType string, metricName string, value float64) bool {
// 	switch metricType {
// 	case "gauge":
// 		if mem.metrics[metricName] == nil {
// 			mem.metrics[metricName] = make([]float64, 1)
// 		}
// 		mem.metrics[metricName][0] = value
// 		return true
// 	case "counter":
// 		if mem.metrics[metricName] == nil {
// 			mem.metrics[metricName] = make([]float64, 0)
// 		}
// 		mem.metrics[metricName] = append(mem.metrics[metricName], value)
// 		return true
// 	default:
// 		return false
// 	}
// }
