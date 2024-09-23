package main

import (
	"context"
	"errors"
	"fmt"
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
	"golang.org/x/sync/errgroup"
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
	metricsStorage.New(context.TODO())

	pollTimeTicker := time.NewTicker(time.Duration(flags.Cfg.PollInterval) * time.Second)
	reportTimeTicker := time.NewTicker(time.Duration(flags.Cfg.ReportInterval) * time.Second)

	logger.SLog.Infow("startup", "config", flags.Cfg)

	ctx, cancel := context.WithCancel(context.Background())

	//Ждем сигнала окончаниея работы
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		<-exit
		logger.SLog.Info("exit signal received")
		cancel()
	}()

	//Собираем метрики по таймеру
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		pollCount := 0
		for {
			select {
			case <-gCtx.Done():
				err := fmt.Errorf("poll metric process stopped")
				logger.SLog.Info(err.Error())
				return err
			case <-pollTimeTicker.C:
				//Собираем метрики
				//делаем новый слайс для метрик
				metricsGathered := make([]metric.Metric, 0)
				//MemStat metric - получаем из runtime.MemStats
				logger.SLog.Info("Gathering MemStatsMetrics")
				if err := metric.PollMemStatsMetrics(metric.MemStatsMetrics, &metricsGathered); err != nil {
					logger.SLog.Errorw("polling metrics", "erorr", err)
					return err
				}
				// Кастомные метрики -  получаем вызывая функции из определения метрики
				logger.SLog.Info("Gathering ExtraMetrics")
				for n, p := range metric.ExtraMetrics {
					var m metric.Metric
					m.Set(n, p.MFunc(pollCount), p.MType)
					metricsGathered = append(metricsGathered, m)
				}
				//сохраним собраные метрики в хранилще
				_, err := metricsStorage.WriteAll(&metricsGathered)
				if err != nil {
					logger.SLog.Errorw("saving metrics to storage", "erorr", err)
					return err
				}
				//очищаем слайс
				metricsGathered = nil
				pollCount++
			}
		}
	})
	//Отправляем метрики по таймеру
	g.Go(func() error {
		for {
			select {
			case <-gCtx.Done():
				err := fmt.Errorf("report metric process stopped")
				logger.SLog.Info(err.Error())
				return err
			case <-reportTimeTicker.C:
				err := errors.New("trying to send")
				retry := 0

				//читам метрики из хранилища для передачи
				metricsToSend := make([]metric.Metric, 0)
				if _, err := metricsStorage.ReadAll(&metricsToSend); err != nil {
					return err
				}

				//Отправляем метрики
				//При ошибке -  пробуем еще раз с задежкой
				for err != nil {
					logger.Log.Info("reporting metric")
					//TODO поменять time.Sleep на TimeTicker
					time.Sleep(time.Duration(flags.RetryIntervals[retry]) * time.Second)
					err = report.ReportMetrics(metricsToSend, flags.Cfg.EndPointAddress)
					if err == nil ||
						retry == len(flags.RetryIntervals)-1 ||
						gCtx.Err() != nil {
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
			}
		}
	})

	if err := g.Wait(); err != nil {
		logger.SLog.Errorw("agent", " error ", err)
	}

	reportTimeTicker.Stop()
	pollTimeTicker.Stop()
	logger.Log.Info("agent stopped")
}
