package flags

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
)

// Config contains application  configuration
type Config struct {
	//server params
	StoreInterval   uint   `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	//Agent params
	//REPORT_INTERVAL позволяет переопределять reportInterval.
	ReportInterval int `env:"REPORT_INTERVAL"`
	//POLL_INTERVAL позволяет переопределять pollInterval.
	PollInterval int `env:"POLL_INTERVAL"`

	//Ключи общие для сесрвера и клиента
	//ADDRESS отвечает за адрес эндпоинта HTTP-сервера
	EndPointAddress string `env:"ADDRESS"`
	//Ключ для SHA256 хэширования
	Key string `env:"KEY"`
	//Количество исходящих запросо к серверу
	//По имуолчанию 1 - отправлять одним пакетом
	RateLimit int `env:"RATE_LIMIT"`

	//Среда - окружение разроботка или продакшан
	Env string `env:"ENV"`
}

// Cfg - application configuration parameters read from command line or Env
var Cfg Config

// RetryInterval in time to repeat functions in case of connection  errors
// Ингтервалы для повторений при ошибках соединения и ввода-вывода
var RetryIntervals = []time.Duration{0, 1, 3, 5}

// Read server config from commanline parameters or Env
func ReadServerConfig() error {
	//Сначала считаем командную строку если есть или заполним конфиг дефолтом

	//Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080)
	//`env:"ADDRESS"`
	flag.StringVar(&Cfg.EndPointAddress, "a", ":8080", "endpoint address")
	flag.StringVar(&Cfg.Key, "k", ``, "sha256 key")

	//Server
	//Флаг -i=<ЗНАЧЕНИЕ> интервал времени в секундах, по истечении которого текущие показания
	//сервера сохраняются на диск (по умолчанию 300 секунд, значение 0 делает запись синхронной)
	flag.UintVar(&Cfg.StoreInterval, "i", 300, "data store interval, sec")

	//Флаг -f=<ЗНАЧЕНИЕ> путь до файла, куда сохраняются текущие значения.
	// flag.StringVar(&Cfg.FileStoragePath, "f", `d:\temp\filest.txt`, "file storage path")
	flag.StringVar(&Cfg.FileStoragePath, "f", `default`, "file storage path")

	//Флаг -r=<ЗНАЧЕНИЕ>  булево значение (true/false), определяющее, загружать или нет ранее
	//сохранённые значения из указанного файла при старте сервера (по умолчанию true)
	flag.BoolVar(&Cfg.Restore, "r", true, "restore data on startup")

	//Флаг -d=<значение> -  адрес подключения к БД / string
	//flag.StringVar(&Cfg.DatabaseDSN, "d", `postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable`, "database dsn adress")
	flag.StringVar(&Cfg.DatabaseDSN, "d", `default`, "database dsn adress")

	//Флаг -e=<ЗНАЧЕНИЕ> позволяет переопределять среду DEV или PROD
	flag.StringVar(&Cfg.Env, "e", "PROD", "environment")

	flag.Parse()
	//Читаем переменные среды и если есть -  перезаписываем параметра ком строки или дефолты
	//в соответствии с условием задания - высший приоритет у переменных окружения
	return env.Parse(&Cfg)
}

// Read agent config from  commandline parameters or Env
func ReadAgentConfig() error {
	//Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080)
	flag.StringVar(&Cfg.EndPointAddress, "a", ":8080", "endpoint address")
	flag.StringVar(&Cfg.Key, "k", ``, "sha256 key")
	flag.IntVar((&Cfg.RateLimit), "l", 1, "rate limit")

	//Agent
	//Флаг -r=<ЗНАЧЕНИЕ> позволяет переопределять reportInterval — частоту отправки метрик на сервер (по умолчанию 10 секунд).
	flag.IntVar(&Cfg.ReportInterval, "r", 10, "report interval")
	//Флаг -p=<ЗНАЧЕНИЕ> позволяет переопределять pollInterval — частоту опроса метрик из пакета runtime (по умолчанию 2 секунды).
	flag.IntVar(&Cfg.PollInterval, "p", 2, "poll interval")

	//Флаг -e=<ЗНАЧЕНИЕ> позволяет переопределять среду DEV или PROD
	flag.StringVar(&Cfg.Env, "e", "PROD", "environment")

	flag.Parse()

	return env.Parse(&Cfg)
}
