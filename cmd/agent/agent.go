package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/Ssnakerss/practicum-metrics/internal/report"
)

const (
	pollInterval   = 2
	reportInterval = 10
)

func main() {

	var gatheredMetrics [29]metric.Metric
	//Initialize metrics array for use
	for idx := range gatheredMetrics {
		var m metric.Metric
		m.Set("testgauge", "0", "gauge")
		gatheredMetrics[idx] = m
	}

	var cnt uint64 = 0
	rp := 0
	for {

		if rp == reportInterval {
			//It's time to report metrics
			fmt.Print("Reporting metrics ... \r")
			report.ReportMetrics(gatheredMetrics[:])
			rp = 0
		}
		time.Sleep(pollInterval * time.Second)
		rp += pollInterval

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
