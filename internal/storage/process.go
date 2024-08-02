package storage

import (
	"github.com/Ssnakerss/practicum-metrics/internal/metric"
)

func ProcessMetric(m metric.Metric, stor *Storage) error {
	//Processing metrics values
	var err error
	switch m.Type {
	case "gauge":
		err = stor.Update(m)
	case "counter":
		err = stor.Insert(m)
	}
	return err
}
