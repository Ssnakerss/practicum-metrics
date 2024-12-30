package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/Ssnakerss/practicum-metrics/internal/agent"
	"github.com/Ssnakerss/practicum-metrics/internal/app"

	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/Ssnakerss/practicum-metrics/internal/metric"
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

	ctx, cancel := context.WithCancel(context.Background())

	a, err := agent.New(logger.SLog)
	if err != nil {
		logger.SLog.Fatalf("cannot create agent: %s", err.Error())
	}

	//Ждем сигнала окончаниея работы
	go app.SysCallProcess(ctx, cancel)

	a.Run(ctx)

}
