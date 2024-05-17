package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

var conf *Config

func SetConfig(c *Config) {
	conf = c
}

func GetConfig() *Config {
	return conf
}

type Config struct {
	DeploymentEnv                 string   `json:"deployment_env"`
	ServiceName                   string   `json:"service_name"`
	ServiceGrpcEndpoint                   string   `json:"service_grpc_endpoint"`
	EnabletHttpGateway                     bool     `json:"enable_http_gateway"`
	ServiceHttpEndpoint                   string   `json:"service_http_endpoint"`
	ServiceMetricsEndpoint                   string   `json:"service_metrics_endpoint"`
	LogLevel                      string   `json:"log_level"`
	LogPretty                     bool     `json:"log_pretty"`
	LogSentryDSN                  string   `json:"log_sentry_dsn"`
	LogPrinter                  string   `json:"log_printer"`
	LogPrinterFilename                  string   `json:"log_printer_filename"`
	EnableReflection                     bool     `json:"enable_reflection"`
}

func loadConfigFile(fn string, ptr interface{}) error {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, ptr); err != nil {
		return err
	}
	return nil
}

func LoadConfigFileOrPanic(fn string, ptr interface{}) {
	if err := loadConfigFile(fn, ptr); err != nil {
		logrus.WithError(err).Fatalf("Failed to load config file:%s.", fn)
	}
}
