package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/Ssnakerss/practicum-metrics/internal/report"
)

const (
// pollInterval = 2
// reportInterval = 10
)

func main() {

	// //•	Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
	endPointAddress := flag.String("a", `http://localhost:8080/`, "endpoint address")
	// //•	Флаг -r=<ЗНАЧЕНИЕ> позволяет переопределять reportInterval — частоту отправки метрик на сервер (по умолчанию 10 секунд).
	reportInterval := flag.Int("r", 10, "report interval")
	// //•	Флаг -p=<ЗНАЧЕНИЕ> позволяет переопределять pollInterval — частоту опроса метрик из пакета runtime (по умолчанию 2 секунды).
	pollInterval := flag.Int("p", 2, "poll interval")
	flag.Parse()

	var gatheredMetrics [29]metric.Metric
	//Initialize metrics array for use
	for idx := range gatheredMetrics {
		var m metric.Metric
		m.Set("testgauge", "0", "gauge")
		gatheredMetrics[idx] = m
	}
	fmt.Println("Agent started")
	fmt.Printf("Poll: %dsec, report: %dsec, endpoint:%s\n\r", *pollInterval, *reportInterval, *endPointAddress)

	var cnt uint64 = 0
	rp := 0
	for {
		if rp == *reportInterval {
			//It's time to report metrics
			fmt.Print("Reporting metrics ... \r")
			report.ReportMetrics(gatheredMetrics[:], *endPointAddress)
			rp = 0
		}

		time.Sleep(time.Duration(*pollInterval) * time.Second)
		rp += *pollInterval

		fmt.Printf("%d:Gathering metrics ... \r", cnt)
		_, err := metric.PollMemStatsMetrics(metric.MemStatsMetrics[:], gatheredMetrics[:])
		if err != nil {
			panic(err)
		}
		cnt++
		gatheredMetrics[27].Set("PollCount", strconv.FormatUint(cnt, 10), "counter")
		gatheredMetrics[28].Set("RandomValue", strconv.FormatFloat(rand.Float64(), 'f', -1, 64), "gauge")
	}
}
