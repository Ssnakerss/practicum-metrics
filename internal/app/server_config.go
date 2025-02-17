package app

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v6"
)

type serverJSONConfig struct {
	Address       string `json:"address"`
	CryptoKey     string `json:"crypto_key"`
	TrustedSubnet string `json:"trusted_subnet"`

	GrpcAddress string `json:"grpc_address"`

	Restore       bool   `json:"restore"`
	StoreInterval string `json:"store_interval"`
	StoreFile     string `json:"store_file"`
	DatabaseDSN   string `json:"database_dsn"`
}

func parseServerConfig(sjcfg *serverJSONConfig, scfg *ServerConfig) {
	scfg.Address = sjcfg.Address
	scfg.CryptoKey = sjcfg.CryptoKey
	scfg.Restore = sjcfg.Restore

	storeInterval, err := strconv.ParseUint(strings.TrimRight(sjcfg.StoreInterval, "s"), 10, 64)
	if err == nil {
		scfg.StoreInterval = storeInterval
	}

	scfg.StoreFile = sjcfg.StoreFile
	scfg.DatabaseDSN = sjcfg.DatabaseDSN
}

type ServerConfig struct {
	//server params
	// -i=<ЗНАЧЕНИЕ> интервал времени в секундах, по истечении которого текущие показания
	//сервера сохраняются на диск (по умолчанию 300 секунд, значение 0 делает запись синхронной)
	StoreInterval uint64 `env:"STORE_INTERVAL"`
	// -f=<ЗНАЧЕНИЕ> путь до файла, куда сохраняются текущие значения.
	StoreFile string `env:"FILE_STORAGE_PATH"`
	// -r=<ЗНАЧЕНИЕ>  булево значение (true/false), определяющее, загружать или нет ранее
	//сохранённые значения из указанного файла при старте сервера (по умолчанию true)
	Restore bool `env:"RESTORE"`
	// -d=<значение> -  адрес подключения к БД / string
	//flag.StringVar(&Cfg.DatabaseDSN, "d", `postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable`, "database dsn adress")
	DatabaseDSN string `env:"DATABASE_DSN"`
	//--common config section  -----
	// - a ADDRESS отвечает за адрес эндпоинта HTTP-сервера
	Address string `env:"ADDRESS"`
	// -k Ключ для SHA256 хэширования
	Key string `env:"KEY"`
	// -crypto-key Путь дофайла с Ключом для шифрования
	CryptoKey string `env:"CRYPTO_KEY"`

	// -c / -config Путь до файла с конфигом
	CFile      string `env:"CONFIG"`
	ConfigFile string `env:"CONFIG"`

	//-t trusted subnet - адрес доверенной сети для приема данных от агента
	TrustedSubnet string `env:"TRUSTED_SUBNET"`

	//-g addree:port of grpc server to receive incoming connectiob
	GrpcAddress string `env:"GRPC_ADDRESS"`

	// -e=<ЗНАЧЕНИЕ> отвечает за среду DEV или PROD
	Env string `env:"ENV"`
}

func (c *ServerConfig) read() error {
	//prepare and read commanline params
	//without default values
	flag.Uint64Var(&c.StoreInterval, "i", c.StoreInterval, "data store interval, sec")
	flag.StringVar(&c.StoreFile, "f", c.StoreFile, "file storage path")
	flag.BoolVar(&c.Restore, "r", c.Restore, "restore data on startup")
	flag.StringVar(&c.DatabaseDSN, "d", c.DatabaseDSN, "database dsn adress")

	flag.StringVar(&c.Address, "a", c.Address, "endpoint address")
	flag.StringVar(&c.Key, "k", c.Key, "sha256 key")
	flag.StringVar(&c.CryptoKey, "crypto-key", c.CryptoKey, "rsa key file path")
	flag.StringVar(&c.Env, "e", c.Env, "environment")

	flag.StringVar(&c.TrustedSubnet, "t", c.TrustedSubnet, "trusted subnet")

	flag.StringVar(&c.GrpcAddress, "g", c.GrpcAddress, "grpc server address")

	flag.StringVar(&c.CFile, "c", "", "config file path")
	flag.StringVar(&c.ConfigFile, "config", "", "config file path")

	flag.Parse()

	//reading environment params
	//overwrites preavious values from command line
	return env.Parse(c)
}

func MakeServerConfig() *ServerConfig {
	s := ServerConfig{
		StoreInterval: 300,
		Restore:       true,
		StoreFile:     "default",
		DatabaseDSN:   "default",
		Address:       "localhost:8080",
		Key:           "",
		CryptoKey:     "",
		TrustedSubnet: "",
		GrpcAddress:   "",
		Env:           "PROD",
	}
	//First check for json file path parameter
	cfgFilePath := ""
	// -c|-config - path to config file
	s.read()
	if s.CFile != "" {
		cfgFilePath = s.CFile
	}
	if s.ConfigFile != "" {
		cfgFilePath = s.ConfigFile
	}

	//If jsonfile path is not empty is set - load params from json file as default values
	if cfgFilePath != "" {
		cJSON := serverJSONConfig{}
		err := loadJSONConfig(cfgFilePath, &cJSON)
		if err == nil {
			parseServerConfig(&cJSON, &s)
		}
	}
	fmt.Println("1: ", s)

	//Reading params from command line and environment after setting default values
	// flag.Parse()
	fmt.Println("2: ", s)
	// env.Parse(s)
	fmt.Println("3: ", s)

	return &s
}
