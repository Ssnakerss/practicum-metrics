package agent

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

func GenMemStats(metricsToGather []string) chan metric.Metric {
	result := make(chan metric.Metric)

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
				var m metric.Metric
				m.Set(name, value, "gauge")
				result <- m
			}
		}
	}()

	return result
}

func GenGopsStats() chan metric.Metric {
	result := make(chan metric.Metric)

	go func() {
		defer close(result)
		v, _ := mem.VirtualMemory()
		mm := metric.Metric{
			Name:  "TotalMemory",
			Gauge: float64(v.Total),
			Type:  "gauge",
		}
		result <- mm

		mm = metric.Metric{
			Name:  "FreeMemory",
			Gauge: float64(v.Free),
			Type:  "gauge",
		}
		result <- mm

		cpu, _ := cpu.Percent(0, true)

		for i, c := range cpu {
			m := metric.Metric{
				Name:  fmt.Sprintf("CPUutilization%d", i),
				Gauge: float64(c),
				Type:  "gauge",
			}
			result <- m
		}

	}()
	return result
}

func GenExtraStats(pollCount int) chan metric.Metric {
	result := make(chan metric.Metric)
	go func() {
		defer close(result)
		for n, p := range metric.ExtraMetrics {
			var m metric.Metric
			m.Set(n, p.MFunc(pollCount), p.MType)
			result <- m
		}
	}()
	return result
}

func CollectMetrics(pollCount int) []metric.Metric {
	//собираем каналы  с метриками  в массив
	channelsToRead := [](chan metric.Metric){
		GenMemStats(metric.MemStatsMetrics),
		GenGopsStats(),
		GenExtraStats(pollCount),
	}

	//канал в который горутины будут писавть метрики
	resultsChannel := make(chan metric.Metric)
	//канал в который горутина пишет что закончила работу
	doneChannel := make(chan struct{})
	for _, ch := range channelsToRead {
		go channelReader(ch, resultsChannel, doneChannel)
	}
	cnt := 0
	res := make([]metric.Metric, 0)
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

func channelReader(ch chan metric.Metric,
	out chan metric.Metric,
	done chan struct{}) {
	for m := range ch {
		out <- m
	}
	done <- struct{}{}
}
