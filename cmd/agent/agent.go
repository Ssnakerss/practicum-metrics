package main

import (
	"log"
	"time"

	"github.com/Ssnakerss/practicum-metrics/internal/flags"
	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/Ssnakerss/practicum-metrics/internal/report"
)

func main() {

	if err := flags.ReadAgentConfig(); err != nil {
		logger.SLog.Warnw("error getting env params", "error", err)
	}

	gatheredMetrics := make(map[string]metric.Metric)
	log.Println("Agent started")
	log.Printf("Poll: %dsec, report: %dsec, endpoint:%s\n\r",
		flags.Cfg.PollInterval,
		flags.Cfg.ReportInterval,
		flags.Cfg.EndPointAddress)

	cnt := 0
	rp := 0
	for {
		if rp == flags.Cfg.ReportInterval {
			//It's time to report metrics
			log.Print("Reporting metrics ... \r")
			report.ReportMetrics(gatheredMetrics, flags.Cfg.EndPointAddress)
			rp = 0
		}

		time.Sleep(time.Duration(flags.Cfg.PollInterval) * time.Second)
		rp += flags.Cfg.PollInterval
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
