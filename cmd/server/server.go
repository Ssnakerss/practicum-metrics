package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ssnakerss/practicum-metrics/internal/flags"
	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"golang.org/x/sync/errgroup"
)

func main() {
	// cоздаем логгер ZAP
	// не получится - проолжать не имеет смысла, fatal
	if err := logger.Initialize("DEBUG"); err != nil {
		log.Fatal("FATAL: cannot initialize LOGGER: ", err)
	}
	defer logger.Log.Sync()
	//Перехватываем паники
	defer func() {
		if err := recover(); err != nil {
			logger.SLog.Fatalf(
				"error heppened while operating -> program will exit",
				"error", err)
		}
	}()

	//Читаем конфигурацию
	if err := flags.ReadServerConfig(); err != nil {
		logger.SLog.Warnw("error getting env params", "error", err)
	}

	//-----------------------------------------------------------------
	//Основновной контекст
	ctx, cancel := context.WithCancel(context.Background())

	//Создаем адаптер для хэндлеров и работы с хранилищем
	da, err := InitAdapter(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//Configuring CHI
	r := NewRouter(da)

	//https://habr.com/ru/articles/771626/
	//Ждем сигнал заверешения работы для остановки сервисов
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		sig := <-exit
		logger.SLog.Infow("received termination", "signal", sig)
		logger.Log.Info("sync storage")
		da.DoSync()
		logger.Log.Info("cancel operations")
		cancel()

		//Закрываем хранилище - актуально для БД
		logger.Log.Info("closing storage")
		da.Ds.Close()
		logger.Log.Fatal("program operation stopped")
	}()
	//-----------------------------------------------------------------------
	httpServer := &http.Server{
		Addr:    flags.Cfg.EndPointAddress,
		Handler: r,
	}
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		logger.SLog.Infow("starting server at ", "address", flags.Cfg.EndPointAddress)
		return httpServer.ListenAndServe() //Запускаем сервер
	})
	g.Go(func() error {
		<-gCtx.Done() //Ожидаем завершения контекста
		logger.Log.Info("shutting server down")
		return httpServer.Shutdown(context.Background()) //Завершаем сервер
	})
	if err := g.Wait(); err != nil {
		logger.SLog.Warnw("server stopped", "error", err)
	}

}
