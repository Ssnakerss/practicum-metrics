package storage

import (
	"fmt"

	"github.com/Ssnakerss/practicum-metrics/internal/metric"
)

// Key for MAP:   matricName@metricType
type Storage struct {
	metrics map[string]metric.Metric
}

// New - initialize storage
func (st *Storage) New() {
	st.metrics = make(map[string]metric.Metric)
}

// Insert - add new value record
func (st *Storage) Update(m metric.Metric) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	mKey := m.Name + "@" + m.MType
	st.metrics[mKey] = m
	return nil
}

// Update - update existing value or add new if missing
func (st *Storage) Insert(m metric.Metric) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	mKey := m.Name + "@" + m.MType
	if v, ok := st.metrics[mKey]; ok {
		v.Value = append(v.Value, m.Value[0])
		st.metrics[mKey] = v
	} else {
		st.metrics[mKey] = m
	}

	return nil
}

// namesAndTypes =metricName@metricType
func (st *Storage) Select(results map[string]metric.Metric, namesAndTypes ...string) (found int, err error) {
	found = 0
	for _, n := range namesAndTypes {
		//Return specific values
		results[n] = st.metrics[n]
		found++
	}

	if len(namesAndTypes) == 0 {
		//Return all values
		for k := range st.metrics {
			results[k] = st.metrics[k]
			found++
		}
	}

	if found == 0 {
		return -1, fmt.Errorf("data not found")
	}
	return found, nil

}
