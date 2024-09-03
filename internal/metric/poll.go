package metric

import (
	"fmt"
	"reflect"
	"runtime"
)

func PollMemStatsMetrics(metricsToGather []string,
	result *[]Metric) error {

	var memoryStat runtime.MemStats
	runtime.ReadMemStats(&memoryStat)
	val := reflect.ValueOf(memoryStat)
	idx := 0
	for _, k := range metricsToGather {
		field := val.FieldByName(k)
		var name, value string
		nm, ok := val.Type().FieldByName(k)
		if ok {
			name = nm.Name
		} else {
			return fmt.Errorf("error metric not found: %s\n\r", k)
		}
		value = fmt.Sprintf("%v", field)
		//----------------------------
		var m Metric
		m.Set(name, value, "gauge")
		*result = append((*result), m)
		idx++
	}
	if idx == 0 {
		return fmt.Errorf("no metric with name %v found", metricsToGather)
	}
	if idx != len(metricsToGather) {
		return fmt.Errorf("not all metrics found")
	}
	return nil
}
