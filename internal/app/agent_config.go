package app

import (
	"flag"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v6"
)

type agentJSONConfig struct {
	Address   string `json:"address"`
	CryptoKey string `json:"crypto_key"`

	GrpcAddress string `json:"grpc_address"`

	ReportInterval string `json:"report_interval"`
	PollInterval   string `json:"poll_interval"`
}

func parseAgentConfig(ajcfg *agentJSONConfig, acfg *AgentConfig) {
	acfg.Address = ajcfg.Address
	acfg.CryptoKey = ajcfg.CryptoKey

	reportInterval, err := strconv.ParseUint(strings.TrimRight(ajcfg.ReportInterval, "s"), 10, 64)
	if err == nil {
		acfg.ReportInterval = reportInterval
	}

	pollInterval, err := strconv.ParseUint(strings.TrimRight(ajcfg.PollInterval, "s"), 10, 64)
	if err == nil {
		acfg.PollInterval = pollInterval
	}
}

type AgentConfig struct {
	// -r позволяет переопределять reportInterval — частоту отправки метрик на сервер (по умолчанию 10 секунд).
	ReportInterval uint64 `env:"REPORT_INTERVAL"`
	// -p  позволяет переопределять pollInterval. default = 2 sec
	PollInterval uint64 `env:"POLL_INTERVAL"`
	// -l=<ЗНАЧЕНИЕ> отвечает за количество исходящих запросов к серверу. default =  1 sec
	RateLimit int `env:"RATE_LIMIT"`
	//--common config section  -----
	// - a ADDRESS отвечает за адрес эндпоинта HTTP-сервера
	Address string `env:"ADDRESS"`
	// -k Ключ для SHA256 хэширования
	Key string `env:"KEY"`
	// -crypto-key Путь дофайла с Ключом для шифрования
	CryptoKey string `env:"CRYPTO_KEY"`

	//-g address:port of grpc server to send data to
	GrpcAddress string `env:"GRPC_ADDRESS"`

	// -c / -config Путь до файла с конфигом
	CFile      string `env:"CONFIG"`
	ConfigFile string `env:"CONFIG"`

	// -e=<ЗНАЧЕНИЕ> отвечает за среду DEV или PROD
	Env string `env:"ENV"`
}

func (c *AgentConfig) read() error {
	//prepare and read commanline params
	//without default values
	flag.IntVar((&c.RateLimit), "l", 1, "rate limit")
	flag.Uint64Var(&c.ReportInterval, "r", 10, "report interval")
	flag.Uint64Var(&c.PollInterval, "p", 2, "poll interval")

	flag.StringVar(&c.Address, "a", c.Address, "endpoint address")
	flag.StringVar(&c.Key, "k", c.Key, "sha256 key")
	flag.StringVar(&c.CryptoKey, "crypto-key", c.CryptoKey, "rsa key file path")
	flag.StringVar(&c.Env, "e", c.Env, "environment")

	flag.StringVar(&c.GrpcAddress, "g", c.GrpcAddress, "grpc server address")

	flag.StringVar(&c.CFile, "c", "", "conifg file path")
	flag.StringVar(&c.ConfigFile, "config", "", "conifg file path")

	flag.Parse()

	//reading environment params
	//overwrites preavious values from command line
	return env.Parse(c)
}

func MakeAgentConfig() *AgentConfig {
	a := AgentConfig{
		ReportInterval: 10,
		PollInterval:   2,
		RateLimit:      1,
		Address:        "localhost:8080",
		Key:            "",
		CryptoKey:      "",

		GrpcAddress: "",

		Env: "DEV",
	}
	//First check for json file path parameter
	cfgFilePath := ""
	// -c|-config - path to config file
	a.read()

	if a.CFile != "" {
		cfgFilePath = a.CFile
	}
	if a.ConfigFile != "" {
		cfgFilePath = a.ConfigFile
	}

	//If jsonfile path is not empty is set - load params from json file as default values
	if cfgFilePath != "" {
		cJSON := agentJSONConfig{}
		err := loadJSONConfig(cfgFilePath, &cJSON)
		if err == nil {
			parseAgentConfig(&cJSON, &a)
		}
	}
	//Reading params from command line and environment after setting default values
	flag.Parse()
	env.Parse(a)

	return &a
}
