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

	st.metrics[m.Name+m.Type] = m
	return nil
}

// Update - update existing value or add new if missing
func (st *Storage) Insert(m metric.Metric) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	if v, ok := st.metrics[m.Name+m.Type]; ok {
		v.Counter += m.Counter
		st.metrics[m.Name+m.Type] = v
	} else {
		st.metrics[m.Name+m.Type] = m
	}
	return nil
}

// namesAndTypes =metricName@metricType
func (st *Storage) Select(results map[string]metric.Metric, names ...metric.Metric) int {
	found := 0
	for _, m := range names {
		//Return specific values
		if m, ok := st.metrics[m.Name+m.Type]; ok {
			results[m.Name+m.Type] = m
			found++
		}
	}

	if len(names) == 0 {
		//Return all values
		for k := range st.metrics {
			results[k] = st.metrics[k]
			found++
		}
	}
	return found
}
