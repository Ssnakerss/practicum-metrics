package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Ssnakerss/practicum-metrics/internal/app"
	"github.com/Ssnakerss/practicum-metrics/internal/encrypt"
	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/Ssnakerss/practicum-metrics/internal/metric"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type sharedSlice struct {
	m     sync.Mutex
	Slice []metric.Metric
}

type Agent struct {
	c *app.AgentConfig
	l *zap.SugaredLogger
	s sharedSlice
	e encrypt.Coder //Кодировка

}

func New(l *zap.SugaredLogger) (*Agent, error) {
	c := app.MakeAgentConfig()
	e := encrypt.Coder{}
	e.LoadPublicKey(c.CryptoKey)

	return &Agent{
		c: c,
		l: l,
		e: e,
	}, nil
}

func (a *Agent) Run(ctx context.Context) {
	a.l.Infow("startup", "config", a.c)

	pollTimeTicker := time.NewTicker(time.Duration(a.c.PollInterval) * time.Second)
	reportTimeTicker := time.NewTicker(time.Duration(a.c.ReportInterval) * time.Second)

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
				a.s.m.Lock()
				logger.SLog.Infof("#%d poll  metrics", pollCount)
				a.s.Slice = CollectMetrics(pollCount)
				a.s.m.Unlock()
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
				a.s.m.Lock()
				metricsToSend := a.s.Slice
				a.s.Slice = nil
				a.s.m.Unlock()

				//Отправляем метрики
				a.SendMetrics(gCtx, metricsToSend)
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
