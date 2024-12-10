package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Ssnakerss/practicum-metrics/internal/app"
	"github.com/Ssnakerss/practicum-metrics/internal/flags"
	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/Ssnakerss/practicum-metrics/internal/server"
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

func main() {
	//Print app build info
	PrintAppInfo()

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
	da, err := server.InitAdapter(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//Configuring CHI
	r := server.NewRouter(da)

	//Ждем сигнал заверешения работы для остановки сервисов

	go app.CtrlC(ctx, cancel, da.DoSync, da.Ds.Close)

	//-----------------------------------------------------------------------

	//https://habr.com/ru/articles/771626/
	httpServer := &http.Server{
		Addr:    flags.Cfg.EndPointAddress,
		Handler: r,
	}
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		logger.SLog.Infow("startup", "config", flags.Cfg)
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
