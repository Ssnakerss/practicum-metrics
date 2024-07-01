package main

import (
	"fmt"
	"lib/metric"
	"runtime"
)

func pollMetrics(mm *[]metric.Metric) (bool, error) {
	var memoryStat runtime.MemStats
	runtime.ReadMemStats(&memoryStat)

	fields := &memoryStat
	for idx, d := range []interface{}{
		fields.Alloc,
		fields.BuckHashSys,
		fields.Frees,
		fields.GCCPUFraction,
		fields.GCSys,
		fields.HeapAlloc,
		fields.HeapIdle,
		fields.HeapInuse,
		fields.HeapObjects,
		fields.HeapReleased,
		fields.HeapSys,
		fields.LastGC,
		fields.Lookups,
		fields.MCacheInuse,
		fields.MCacheSys,
		fields.MSpanInuse,
		fields.MSpanSys,
		fields.Mallocs,
		fields.NextGC,
		fields.NumForcedGC,
		fields.NumGC,
		fields.OtherSys,
		fields.PauseTotalNs,
		fields.StackInuse,
		fields.StackSys,
		fields.Sys,
		fields.TotalAlloc,
	} {
		switch v := d.(type) {
		case uint64:
			(*mm)[idx].metricValue = float64(d.(uint64))
		case uint32:
			(*mm)[idx].metricValue = float64(d.(uint32))
		case float64:
			(*mm)[idx].metricValue = d.(float64)
		default:
			return false, fmt.Errorf("unexpected type %v", v)
		}
	}
	(*mm)[27].metricValue += 1
	(*mm)[28].metricValue = rand.Float64()
	return true, nil
}
