package main

import (
	"context"
	"flag"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/evalphobia/logrus_sentry"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"{{RepoBase}}/{{RepoGroup}}/{{RepoName}}/internal/common/config"
	"{{RepoBase}}/{{RepoGroup}}/{{RepoName}}/internal/common/logger"
	"{{RepoBase}}/{{RepoGroup}}/{{RepoName}}/internal/metrics"
	"{{RepoBase}}/{{RepoGroup}}/{{RepoName}}/internal/service"
)

var (
	_ConfigFile = flag.String("conf", "./etc/{{ServiceName}}-dev.json", "the config file")
)

func main() {
	rand.Seed(time.Now().UnixNano())

	flag.Parse()

	var conf config.Config
	config.LoadConfigFileOrPanic(*_ConfigFile, &conf)
	config.SetConfig(&conf)

	prepareEnv(conf)
	defer cleanEnv(conf)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stopCh := make(chan struct{})
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	wg.Add(1)
	go setupGrpcService(ctx, wg, stopCh)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-sigCh

	close(stopCh)
}

func setGlobalLogger(conf config.Config) {
	var _logger *logrus.Logger
	var _level logrus.Level
	var err error

	if len(conf.LogLevel) > 0 {
		_level, err = logrus.ParseLevel(conf.LogLevel)
		if err != nil {
			_level = logrus.DebugLevel
		}
	} else {
		_level = logrus.DebugLevel
	}

	_logger = logrus.New()
	_logger.SetLevel(_level)
	if conf.LogPretty {
		_logger.SetFormatter(&logrus.JSONFormatter{
			DisableHTMLEscape: true,
			PrettyPrint:       true,
		})
	} else {
		_logger.SetFormatter(&logrus.TextFormatter{})
	}

	if len(conf.LogSentryDSN) > 0 {
		hook, err := logrus_sentry.NewSentryHook(conf.LogSentryDSN, []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		})
		if err != nil {
			logrus.WithError(err).Fatal("Failed to set sentry hook for logrus logger.")
		} else {
			_logger.Hooks.Add(hook)
		}
	}

	logger.SetGlobalLogger(_logger, conf.ServiceName)
}

func prepareEnv(conf config.Config) {
	if len(conf.ServiceMetricsEndpoint) > 0 {
		go runMetricsServer(conf.ServiceMetricsEndpoint)
	}
	setGlobalLogger(conf)
	service.Setup{{ServiceNameInCamelCase}}Impl()
}

func cleanEnv(conf config.Config) {
	service.Close{{ServiceNameInCamelCase}}Impl()
}

func runMetricsServer(ep string) {
	metrics.Register()
	http.Handle("/metrics", promhttp.Handler())
	logger.GetGlobalLogger().Error(http.ListenAndServe(ep, nil))
}
