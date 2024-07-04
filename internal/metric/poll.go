package metric

import (
	"fmt"
	"reflect"
	"runtime"
)

func PollMemStatsMetrics(metricsToGather []string, result []Metric) (bool, error) {
	var memoryStat runtime.MemStats
	runtime.ReadMemStats(&memoryStat)
	val := reflect.ValueOf(memoryStat)
	idx := 0
	for _, k := range metricsToGather {

		field := val.FieldByName(k)
		//For better code readability
		var name, value, vtype string
		nm, ok := val.Type().FieldByName(k)
		if ok {
			name = nm.Name
		} else {
			panic("?????")
		}
		value = fmt.Sprintf("%v", field)
		vtype = field.Type().String()
		//----------------------------
		result[idx].Set(name, value, vtype, "gauge")
		idx++
	}
	return true, nil
}
