package app

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ssnakerss/practicum-metrics/internal/logger"
)

// RetryInterval in time to repeat functions in case of connection  errors
// Ингтервалы для повторений при ошибках соединения и ввода-вывода
var RetryIntervals = []time.Duration{0, 1, 3, 5}

func loadJsonConfig(path string, cfg any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, cfg)
	if err != nil {
		return err
	}
	return nil
}

// create channel for exit signal
// wait for signal and canlcel global context
func SysCallProcess(ctx context.Context,
	cancel context.CancelFunc,
	ff ...func(),
) {

	defer cancel()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	)
	select {
	case s := <-exit:
		logger.SLog.Info("received signal: ", "syscal", s.Signal)
	case <-ctx.Done():
	}

	logger.Log.Info("shutting down")
	if len(ff) > 0 {
		logger.SLog.Info("performing pre-shutdown tasks")
		for _, f := range ff {
			f()
		}
	}
}
