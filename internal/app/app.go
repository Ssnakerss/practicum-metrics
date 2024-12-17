package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ssnakerss/practicum-metrics/internal/logger"
)

//create channel for exit signal
//wait for signal and canlcel global context

func SysCallProcess(ctx context.Context,
	cancel context.CancelFunc,
	ff ...func(),
) {
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
		logger.Log.Info("shutting down")
		cancel()
	case <-ctx.Done():
		logger.Log.Info("shutting down")
	}
	logger.SLog.Info("performing pre-shutdown tasks")
	for _, f := range ff {
		f()
	}
}
