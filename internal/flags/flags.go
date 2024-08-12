package flags

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	//server params
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	//Agent params
	//REPORT_INTERVAL позволяет переопределять reportInterval.
	//POLL_INTERVAL позволяет переопределять pollInterval.
	ReportInterval int `env:"REPORT_INTERVAL"`
	PollInterval   int `env:"POLL_INTERVAL"`

	//Common config value
	//ADDRESS отвечает за адрес эндпоинта HTTP-сервера
	EndPointAddress string `env:"ADDRESS"`
}

var Cfg Config

func ReadServerConfig() error {
	//Сначала считаем командную строку если есть или заполним конфиг дефолтом

	//Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080)
	flag.StringVar(&Cfg.EndPointAddress, "a", "localhost:8080", "endpoint address")

	//Server
	//Флаг -i=<ЗНАЧЕНИЕ> интервал времени в секундах, по истечении которого текущие показания
	//сервера сохраняются на диск (по умолчанию 300 секунд, значение 0 делает запись синхронной)
	flag.IntVar(&Cfg.StoreInterval, "i", 300, "data store interval, sec")
	//Флаг -f=<ЗНАЧЕНИЕ> путь до файла, куда сохраняются текущие значения.
	flag.StringVar(&Cfg.FileStoragePath, "f", "", "file storage path")
	//Флаг -r=<ЗНАЧЕНИЕ>  булево значение (true/false), определяющее, загружать или нет ранее
	//сохранённые значения из указанного файла при старте сервера (по умолчанию true)
	flag.BoolVar(&Cfg.Restore, "r", true, "restore data on startup")

	flag.Parse()

	return env.Parse(&Cfg)
}

func ReadAgentConfig() error {
	//Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080)
	flag.StringVar(&Cfg.EndPointAddress, "a", "localhost:8080", "endpoint address")

	//Agent
	//Флаг -r=<ЗНАЧЕНИЕ> позволяет переопределять reportInterval — частоту отправки метрик на сервер (по умолчанию 10 секунд).
	flag.IntVar(&Cfg.ReportInterval, "r", 10, "report interval")
	//Флаг -p=<ЗНАЧЕНИЕ> позволяет переопределять pollInterval — частоту опроса метрик из пакета runtime (по умолчанию 2 секунды).
	flag.IntVar(&Cfg.PollInterval, "p", 2, "poll interval")

	flag.Parse()

	return env.Parse(&Cfg)
}
