package storage

import (
	"fmt"

	"github.com/Ssnakerss/practicum-metrics/internal/metric"
)

type MemStorage struct {
	// DataStorage
	metrics map[string]*metric.Metric
}

// New - initialize storage
func (memst *MemStorage) New(p ...string) error {
	memst.metrics = make(map[string]*metric.Metric)
	return nil
}

func (memst *MemStorage) Write(m *metric.Metric) error {
	if v, ok := memst.metrics[m.Name+m.Type]; ok {
		v.Counter += m.Counter
		v.Gauge = m.Gauge
	} else {
		nm := metric.CopyMetric(m)
		memst.metrics[m.Name+m.Type] = nm
	}
	return nil
}

func (memst *MemStorage) WriteAll(mm *([]metric.Metric)) (int, error) {
	cnt := 0
	for _, m := range *mm {
		err := memst.Write(&m)
		if err != nil {
			return cnt, err
		}
		cnt++
	}
	return cnt, nil
}

func (memst *MemStorage) Read(m *metric.Metric) error {
	if sm, ok := memst.metrics[m.Name+m.Type]; ok {
		m.Gauge = sm.Gauge
		m.Counter = sm.Counter
		return nil
	}
	return fmt.Errorf("metric not found %v", m)
}

func (memst *MemStorage) ReadAll(mm *([]metric.Metric)) (int, error) {
	cnt := 0
	for _, m := range memst.metrics {
		*mm = append(*mm, *m)
		cnt++
	}
	return cnt, nil
}
