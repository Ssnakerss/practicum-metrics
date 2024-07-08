package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/Ssnakerss/practicum-metrics/internal/report"

	"github.com/caarlos0/env/v6"
)

const (
// pollInterval = 2
// reportInterval = 10
)

type Config struct {
	endPointAddress string `env:"ADDRESS"`
	reportInterval  int    `env:"REPORT_INTERVAL"`
	pollInterval    int    `env:"POLL_INTERVAL"`
}

func main() {
	//переменные окружения
	// 	•	ADDRESS отвечает за адрес эндпоинта HTTP-сервера.
	// •	REPORT_INTERVAL позволяет переопределять reportInterval.
	// •	POLL_INTERVAL позволяет переопределять pollInterval.
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	if cfg.endPointAddress == "" {
		// //•	Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
		flag.StringVar(&cfg.endPointAddress, "a", "localhost:8080", "endpoint address")
	}
	if cfg.reportInterval == 0 {
		// //•	Флаг -r=<ЗНАЧЕНИЕ> позволяет переопределять reportInterval — частоту отправки метрик на сервер (по умолчанию 10 секунд).
		flag.IntVar(&cfg.reportInterval, "r", 10, "report interval")
	}
	if cfg.pollInterval == 0 {
		// //•	Флаг -p=<ЗНАЧЕНИЕ> позволяет переопределять pollInterval — частоту опроса метрик из пакета runtime (по умолчанию 2 секунды).
		flag.IntVar(&cfg.pollInterval, "p", 2, "poll interval")
	}
	flag.Parse()

	var gatheredMetrics [29]metric.Metric
	//Initialize metrics array for use
	for idx := range gatheredMetrics {
		var m metric.Metric
		m.Set("testgauge", "0", "gauge")
		gatheredMetrics[idx] = m
	}
	fmt.Println("Agent started")
	fmt.Printf("Poll: %dsec, report: %dsec, endpoint:%s\n\r", cfg.pollInterval, cfg.reportInterval, cfg.endPointAddress)

	var cnt uint64 = 0
	rp := 0
	for {
		if rp == cfg.reportInterval {
			//It's time to report metrics
			fmt.Print("Reporting metrics ... \r")
			report.ReportMetrics(gatheredMetrics[:], cfg.endPointAddress)
			rp = 0
		}

		time.Sleep(time.Duration(cfg.pollInterval) * time.Second)
		rp += cfg.pollInterval

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
