package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/Ssnakerss/practicum-metrics/internal/metric"
)

type MemStorage struct {
	//Сделаем мемстор потокобезопасным - добавим мютексы
	mx sync.Mutex
	// DataStorage
	metrics map[string]*metric.Metric
}

// New - initialize storage
func (memst *MemStorage) New(ctx context.Context, p ...string) error {
	memst.metrics = make(map[string]*metric.Metric)
	return nil
}

func (memst *MemStorage) Write(m *metric.Metric) error {
	memst.mx.Lock()
	defer memst.mx.Unlock()

	if v, ok := memst.metrics[m.Name+m.Type]; ok {
		v.Counter += m.Counter
		v.Gauge = m.Gauge
	} else {
		nm := m.Copy()
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
	memst.mx.Lock()
	defer memst.mx.Unlock()

	if sm, ok := memst.metrics[m.Name+m.Type]; ok {
		m.Gauge = sm.Gauge
		m.Counter = sm.Counter
		return nil
	}
	return fmt.Errorf("metric not found %v", m)
}

func (memst *MemStorage) ReadAll(mm *([]metric.Metric)) (int, error) {
	memst.mx.Lock()
	defer memst.mx.Unlock()

	cnt := 0
	for _, m := range memst.metrics {
		*mm = append(*mm, *m)
		cnt++
	}
	return cnt, nil
}

func (memst *MemStorage) Truncate() error {
	memst.mx.Lock()
	defer memst.mx.Unlock()

	//Чистим мапу путем создания новой
	memst.metrics = make(map[string]*metric.Metric)
	return nil
}

func (memst *MemStorage) CheckStorage() error {
	return fmt.Errorf("usnig memory storage, db connection unavailable ")
}

func (memst *MemStorage) Close() {

}
