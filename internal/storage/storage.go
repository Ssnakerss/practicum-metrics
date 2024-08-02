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

	st.metrics[m.Name] = m
	return nil
}

// Update - update existing value or add new if missing
func (st *Storage) Insert(m metric.Metric) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	if v, ok := st.metrics[m.Name]; ok {
		v.Counter += m.Counter
		st.metrics[m.Name] = v
	} else {
		st.metrics[m.Name] = m
	}

	return nil
}

// namesAndTypes =metricName@metricType
func (st *Storage) Select(results map[string]metric.Metric, names ...string) (found int, err error) {
	found = 0
	for _, n := range names {
		//Return specific values
		results[n] = st.metrics[n]
		found++
	}

	if len(names) == 0 {
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
