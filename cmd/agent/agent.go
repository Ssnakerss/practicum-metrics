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
	"github.com/Ssnakerss/practicum-metrics/internal/storage"
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

	//хранилище собранных  метрик
	var metricsStorage storage.MemStorage
	metricsStorage.New()

	pollTimeTicker := time.NewTicker(time.Duration(flags.Cfg.PollInterval) * time.Second)
	reportTimeTicker := time.NewTicker(time.Duration(flags.Cfg.ReportInterval) * time.Second)

	logger.SLog.Infow("Agent started", "poll interval", flags.Cfg.PollInterval,
		"report interval", flags.Cfg.ReportInterval,
		"endpoint address", flags.Cfg.EndPointAddress)

	go func() {
		pollCount := 0
		for range pollTimeTicker.C {
			//Собираем метрики
			//делаем новый слайс для метрик
			metricsGathered := make([]metric.Metric, 0)
			//MemStat metric - получаем из runtime.MemStats
			logger.SLog.Info("Gathering MemStatsMetrics")
			if err := metric.PollMemStatsMetrics(metric.MemStatsMetrics, &metricsGathered); err != nil {
				logger.SLog.Errorw("polling metrics", "erorr", err)
			}
			// Кастомные метрики -  получаем вызывая функции из определения метрики
			logger.SLog.Info("Gathering ExtraMetrics")
			for n, p := range metric.ExtraMetrics {
				var m metric.Metric
				m.Set(n, p.MFunc(pollCount), p.MType)
				metricsGathered = append(metricsGathered, m)
			}

			//сохраним собраные метрики в хранилще
			metricsStorage.WriteAll(&metricsGathered)

			//очищаем слайс
			metricsGathered = nil
			pollCount++
		}
	}()

	go func() {
		for range reportTimeTicker.C {
			err := errors.New("trying to send")
			retry := 0
			//Тормозим тикеры на время передачи данных
			// pollTimeTicker.Stop()
			// reportTimeTicker.Stop()
			//читам метрики из хранилища для передачи
			metricsToSend := make([]metric.Metric, 0)
			metricsStorage.ReadAll(&metricsToSend)

			//Отправляем метрики
			//При ошибке -  пробуем еще раз с задежкой
			for err != nil {
				logger.Log.Info("reporting metric")
				time.Sleep(time.Duration(flags.RetryIntervals[retry]) * time.Second)

				err = report.ReportMetrics(metricsToSend, flags.Cfg.EndPointAddress)
				if err == nil || retry == len(flags.RetryIntervals)-1 {
					break
				}
				retry++
				logger.SLog.Warnf("error reporting, retry in %d seconds", flags.RetryIntervals[retry])

			}
			if err != nil {
				logger.SLog.Errorw("error reporting", "err", err)
			} else {
				logger.SLog.Infof("reported %d metrics", len(metricsToSend))
			}

			//очищаем мапу передачи
			metricsToSend = nil

			//Запускаем таймеры снова
			// reportTimeTicker.Reset(time.Duration(flags.Cfg.ReportInterval) * time.Second)
			// pollTimeTicker.Reset(time.Duration(flags.Cfg.PollInterval) * time.Second)
		}
	}()

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT)
	<-terminateSignals

	reportTimeTicker.Stop()
	pollTimeTicker.Stop()

	logger.Log.Info("agent stopped")
}
