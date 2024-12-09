package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Ssnakerss/practicum-metrics/internal/app"
	"github.com/Ssnakerss/practicum-metrics/internal/flags"
	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"github.com/Ssnakerss/practicum-metrics/internal/report"
	"golang.org/x/sync/errgroup"
)

// global variable for build versioninfo
var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func PrintAppInfo() {
	fmt.Println("Build version: ", buildVersion)
	fmt.Println("Build date: ", buildDate)
	fmt.Println("Build commit: ", buildCommit)

}

type sharedSlice struct {
	m     sync.Mutex
	Slice []metric.Metric
}

func main() {
	//Print app build info
	PrintAppInfo()

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

	go app.CtrlC(ctx, cancel)

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
				logger.SLog.Infof("sending  metrics")
				//читам метрики из хранилища для передачи
				mm.m.Lock()
				metricsToSend := mm.Slice
				mm.Slice = nil
				mm.m.Unlock()

				//Отправляем метрики
				report.SendMetrics(gCtx, metricsToSend)
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
