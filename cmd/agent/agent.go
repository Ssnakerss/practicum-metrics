package main

import (
	"flag"
	"log"
	"time"

	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/Ssnakerss/practicum-metrics/internal/report"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	EndPointAddress string `env:"ADDRESS"`
	ReportInterval  int    `env:"REPORT_INTERVAL"`
	PollInterval    int    `env:"POLL_INTERVAL"`
}

func main() {
	//переменные окружения
	//ADDRESS отвечает за адрес эндпоинта HTTP-сервера.
	//REPORT_INTERVAL позволяет переопределять reportInterval.
	//POLL_INTERVAL позволяет переопределять pollInterval.
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Printf("error getting env prams: %s, contimue with cmd line or default\n\r", err)

	}

	if cfg.EndPointAddress == "" {
		//Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
		flag.StringVar(&cfg.EndPointAddress, "a", "localhost:8080", "endpoint address")
	}
	if cfg.ReportInterval == 0 {
		//Флаг -r=<ЗНАЧЕНИЕ> позволяет переопределять reportInterval — частоту отправки метрик на сервер (по умолчанию 10 секунд).
		flag.IntVar(&cfg.ReportInterval, "r", 10, "report interval")
	}
	if cfg.PollInterval == 0 {
		//Флаг -p=<ЗНАЧЕНИЕ> позволяет переопределять pollInterval — частоту опроса метрик из пакета runtime (по умолчанию 2 секунды).
		flag.IntVar(&cfg.PollInterval, "p", 2, "poll interval")
	}
	flag.Parse()

	gatheredMetrics := make(map[string]metric.Metric)
	log.Println("Agent started")
	log.Printf("Poll: %dsec, report: %dsec, endpoint:%s\n\r", cfg.PollInterval, cfg.ReportInterval, cfg.EndPointAddress)

	cnt := 0
	rp := 0
	for {
		if rp == cfg.ReportInterval {
			//It's time to report metrics
			log.Print("Reporting metrics ... \r")
			report.ReportMetrics(gatheredMetrics, cfg.EndPointAddress)
			rp = 0
		}

		time.Sleep(time.Duration(cfg.PollInterval) * time.Second)
		rp += cfg.PollInterval
		//MemStat metric - получаем из runtime.MemStats
		log.Printf("%d:Gathering MemStatsMetrics ... \r", cnt)

		if err := metric.PollMemStatsMetrics(metric.MemStatsMetrics, gatheredMetrics); err != nil {
			log.Printf("error polling metrics: %s, continue...\r\n", err)
		}
		// Кастомные метрики -  получаем вызывая функции из определения метрики
		log.Printf("%d:Gathering ExtraMetrics ... \r", cnt)
		for n, p := range metric.ExtraMetrics {
			var m metric.Metric
			m.Set(n, p.MFunc(cnt), p.MType)
			gatheredMetrics[n] = m
		}
		cnt++
	}
}
