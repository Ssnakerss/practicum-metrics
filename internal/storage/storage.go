package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/Ssnakerss/practicum-metrics/internal/metric"
)

// Key for MAP:   matricName@metricType
type Storage struct {
	metrics   map[string]*metric.Metric
	syncWrite bool
	filePath  string
}

// New - initialize storage
func (st *Storage) New(file string, sync bool) {
	st.metrics = make(map[string]*metric.Metric)
	st.syncWrite = sync
	st.filePath = file
}

func (st *Storage) SaveMetric(m *metric.Metric) error {
	// Processing metrics values
	var err error
	switch m.Type {
	case "gauge":
		err = st.Update(m)
	case "counter":
		err = st.Insert(m)
	}
	return err
}

func (st *Storage) ReadMetric(mm *metric.Metric) error {
	results := make(map[string]*metric.Metric)
	if found := st.Select(results, mm); found > 0 {
		if nm, ok := results[mm.Name+mm.Type]; ok {
			mm = nm
			return nil
		}
	}
	return fmt.Errorf("metric not found %v", *mm)
}

// Insert - add new value record
func (st *Storage) Update(m *metric.Metric) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	st.metrics[m.Name+m.Type] = m
	//сохраняем изменения в файл
	if st.syncWrite {
		err := st.Save()
		if err != nil {
			logger.SLog.Warnw("failed file", "save", err)
		}
	}
	return nil
}

// Update - update existing value or add new if missing
func (st *Storage) Insert(m *metric.Metric) (err error) {
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
	//сохраняем изменения в файл
	if st.syncWrite {
		err := st.Save()
		if err != nil {
			logger.SLog.Warnw("failed file", "save", err)
		}
	}
	return nil
}

// namesAndTypes =metricName@metricType
func (st *Storage) Select(results map[string]*metric.Metric, names ...*metric.Metric) int {
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

// File oparations
// Читаем из файла
func (st *Storage) Restore() error {
	file, err := os.OpenFile(st.filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return err
	}
	data := scanner.Bytes()
	if err = json.Unmarshal(data, &st.metrics); err != nil {
		return err
	}
	return nil
}

// Пишем в файл
func (st *Storage) Save() error {
	var saveError error
	file, err := os.OpenFile(st.filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && saveError == nil {
			saveError = closeErr
		}
	}()

	writer := bufio.NewWriter(file)
	data, err := json.Marshal(&st.metrics)
	if err != nil {
		return err
	}
	if _, err := writer.Write(data); err != nil {
		return err
	}
	if err := writer.WriteByte('\n'); err != nil {
		return err
	}
	if err := writer.Flush(); err != nil {
		return err
	}
	return saveError
}
