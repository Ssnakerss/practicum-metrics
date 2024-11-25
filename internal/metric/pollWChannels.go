package metric

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

func GenMemStats(metricsToGather []string) chan Metric {
	result := make(chan Metric)

	go func() {
		defer close(result)
		var memoryStat runtime.MemStats
		runtime.ReadMemStats(&memoryStat)
		val := reflect.ValueOf(memoryStat)
		for _, k := range metricsToGather {
			field := val.FieldByName(k)
			var name, value string
			nm, ok := val.Type().FieldByName(k)
			if ok {
				name = nm.Name
				value = fmt.Sprintf("%v", field)
				//----------------------------
				var m Metric
				m.Set(name, value, "gauge")
				result <- m
			}
		}
	}()

	return result
}

func GenGopsStats() chan Metric {
	result := make(chan Metric)

	go func() {
		defer close(result)
		v, _ := mem.VirtualMemory()
		mm := Metric{
			Name:  "TotalMemory",
			Gauge: float64(v.Total),
			Type:  "gauge",
		}
		result <- mm

		mm = Metric{
			Name:  "FreeMemory",
			Gauge: float64(v.Free),
			Type:  "gauge",
		}
		result <- mm

		cpu, _ := cpu.Percent(0, true)

		for i, c := range cpu {
			m := Metric{
				Name:  fmt.Sprintf("CPUutilization%d", i),
				Gauge: float64(c),
				Type:  "gauge",
			}
			result <- m
		}

	}()
	return result
}

func GenExtraStats(pollCount int) chan Metric {
	result := make(chan Metric)
	go func() {
		defer close(result)
		for n, p := range ExtraMetrics {
			var m Metric
			m.Set(n, p.MFunc(pollCount), p.MType)
			result <- m
		}
	}()
	return result
}

func CollectMetrics(pollCount int) []Metric {
	//собираем каналы  с метриками  в массив
	channelsToRead := [](chan Metric){
		GenMemStats(MemStatsMetrics),
		GenGopsStats(),
		GenExtraStats(pollCount),
	}

	//канал в который горутины будут писавть метрики
	resultsChannel := make(chan Metric)
	//канал в который горутина пишет что закончила работу
	doneChannel := make(chan struct{})
	for _, ch := range channelsToRead {
		go channelReader(ch, resultsChannel, doneChannel)
	}
	cnt := 0
	res := make([]Metric, 0)
	for {
		select {
		case <-doneChannel:
			cnt++
		case m := <-resultsChannel:
			res = append(res, m)
		}
		if cnt == len(channelsToRead) {
			break
		}
	}
	close(doneChannel)
	close(resultsChannel)
	return res
}

func channelReader(ch chan Metric,
	out chan Metric,
	done chan struct{}) {
	for m := range ch {
		out <- m
	}
	done <- struct{}{}
}
