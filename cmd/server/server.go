package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Ssnakerss/practicum-metrics/internal/app"

	"github.com/Ssnakerss/practicum-metrics/internal/logger"
	"github.com/Ssnakerss/practicum-metrics/internal/server"
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

	// creating globall loger ZAP
	// if fail -  quit program
	if err := logger.Initialize("DEBUG"); err != nil {
		log.Fatal("FATAL: cannot initialize LOGGER: ", err)
	}
	defer logger.Log.Sync()

	//Capture app panics if any.
	defer func() {
		if err := recover(); err != nil {
			logger.SLog.Fatalf(
				"error heppened while operating -> program will exit",
				"error", err)
		}
	}()

	//-----------------------------------------------------------------
	//Creating main app context.
	ctx, cancel := context.WithCancel(context.Background())

	//Create and inialise server.
	s, err := server.New(ctx, logger.SLog)
	if err != nil {
		logger.SLog.Fatalf("cannot create server: %v", err)
	}
	//start  Listening for SysCal events
	go app.SysCallProcess(ctx, cancel)

	//Starting server
	s.Run(ctx)
}
