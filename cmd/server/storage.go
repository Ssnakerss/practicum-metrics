package main

import (
	"fmt"
	"lib/metric"
)

type Storage struct {
	metrics map[string][]metric.Metric
}

// New - initialize storage
func (st *Storage) New() {
	st.metrics = make(map[string][]metric.Metric)
}

// Insert - add new value record
func (st *Storage) Insert(m metric.Metric) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	st.metrics[m.Name] = append(st.metrics[m.Name], m)
	return nil
}

// Update - update existing value or add new if missing
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
