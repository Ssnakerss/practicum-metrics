package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ssnakerss/practicum-metrics/internal/flags"
	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/Ssnakerss/practicum-metrics/internal/report"
)

func main() {

	// cоздаем логгер ZAP
	// не получится - проолжать не имеет смысла, fatal
	if err := logger.Initialize("DEBUG"); err != nil {
		log.Fatal("FATAL: cannot initialize LOGGER: ", err)
	}
	defer logger.Log.Sync()

	if err := flags.ReadAgentConfig(); err != nil {
		logger.SLog.Warnw("error getting env params", "error", err)
	}

	gatheredMetrics := make(map[string]metric.Metric)
	pollTimeTicker := time.NewTicker(time.Duration(flags.Cfg.PollInterval) * time.Second)
	reportTimeTicker := time.NewTicker(time.Duration(flags.Cfg.ReportInterval) * time.Second)

	logger.SLog.Infow("Agent started", "poll interval", flags.Cfg.PollInterval,
		"report interval", flags.Cfg.ReportInterval,
		"endpoint address", flags.Cfg.EndPointAddress)

	go func() {
		cnt := 0
		for range pollTimeTicker.C {
			//Собираем метрики
			//MemStat metric - получаем из runtime.MemStats
			logger.SLog.Info("Gathering MemStatsMetrics")
			if err := metric.PollMemStatsMetrics(metric.MemStatsMetrics, gatheredMetrics); err != nil {
				logger.SLog.Errorw("polling metrics", "erorr", err)
			}
			// Кастомные метрики -  получаем вызывая функции из определения метрики
			logger.SLog.Info("Gathering ExtraMetrics")
			for n, p := range metric.ExtraMetrics {
				var m metric.Metric
				m.Set(n, p.MFunc(cnt), p.MType)
				gatheredMetrics[n] = m
			}
			cnt++
		}
	}()

	go func() {
		for range reportTimeTicker.C {
			err := errors.New("trying to send")
			retry := 0
			//Отправляем метрики
			//При ошибке -  пробуем еще раз с задежкой
			for err != nil {
				//Тормозим тикеры на время передачи данных
				pollTimeTicker.Stop()
				reportTimeTicker.Stop()

				logger.Log.Info("reporting metric")
				time.Sleep(time.Duration(flags.RetryIntervals[retry]) * time.Second)
				err = report.ReportMetrics(gatheredMetrics, flags.Cfg.EndPointAddress)
				if err == nil || retry == len(flags.RetryIntervals)-1 {
					break
				}
				retry++
				logger.SLog.Warnf("error reporting, retry in %d seconds", flags.RetryIntervals[retry])

			}
			if err != nil {
				logger.SLog.Errorw("error reporting", "err", err)
			}

			//Запускаем таймеры снова
			reportTimeTicker.Reset(time.Duration(flags.Cfg.ReportInterval) * time.Second)
			pollTimeTicker.Reset(time.Duration(flags.Cfg.PollInterval) * time.Second)
		}
	}()

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT)
	<-terminateSignals

	reportTimeTicker.Stop()
	pollTimeTicker.Stop()

	logger.Log.Info("agent stopped")
}
