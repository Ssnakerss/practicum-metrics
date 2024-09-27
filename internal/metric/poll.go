package metric

import (
	"fmt"
	"reflect"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
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

func PollGopsMetrics(metricsToGather []string, result *[]Metric) error {
	v, _ := mem.VirtualMemory()
	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)
	// convert to JSON. String() is also implemented
	fmt.Println(v)
	c, _ := cpu.Percent(4*time.Second, true)
	fmt.Println(c)

	//TO-DO: Implement this
	//Затычка для гопса
	mm := []Metric{
		{
			Name:  "TotalMemory",
			Gauge: 16566.0,
			Type:  "gauge",
		},
		{
			Name:  "FreeMemory",
			Gauge: 1666.0,
			Type:  "gauge",
		},
		{
			Name:  "CPUutilization1",
			Gauge: 89.9,
			Type:  "gauge",
		},
	}
	time.Sleep(time.Millisecond * 50)
	*result = append(*result, mm...)
	return nil
}
