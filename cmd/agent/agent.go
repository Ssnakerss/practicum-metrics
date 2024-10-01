package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Ssnakerss/practicum-metrics/internal/flags"
	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/Ssnakerss/practicum-metrics/internal/report"
	"golang.org/x/sync/errgroup"
)

type sharedSlice struct {
	m     sync.Mutex
	Slice []metric.Metric
}

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
	logger.SLog.Infow("startup", "config", flags.Cfg)

	//хранилище собранных  метрик
	var mm sharedSlice

	pollTimeTicker := time.NewTicker(time.Duration(flags.Cfg.PollInterval) * time.Second)
	reportTimeTicker := time.NewTicker(time.Duration(flags.Cfg.ReportInterval) * time.Second)

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
				mm.m.Lock()
				logger.SLog.Infof("#%d poll  metrics", pollCount)
				mm.Slice = metric.CollectMetrics(pollCount)
				mm.m.Unlock()
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

				//читам метрики из хранилища для передачи
				mm.m.Lock()
				metricsToSend := mm.Slice
				mm.Slice = nil
				mm.m.Unlock()

				//Отправляем метрики
				//При ошибке -  пробуем еще раз с задержкой
				err := errors.New("trying to send")
				retry := 0
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
