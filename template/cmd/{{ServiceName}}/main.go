package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"{{RepoBase}}/{{RepoGroup}}/{{RepoName}}/internal/common/config"
	"{{RepoBase}}/{{RepoGroup}}/{{RepoName}}/internal/common/logger"
	"{{RepoBase}}/{{RepoGroup}}/{{RepoName}}/internal/metrics"
	"{{RepoBase}}/{{RepoGroup}}/{{RepoName}}/internal/service"
)

var (
	_ConfigFile = flag.String("conf", "./etc/{{ServiceName}}-dev.json", "config file path")
)

func main() {
	// Seed random number generator, is deprecated since Go 1.20.
	// rand.Seed(time.Now().UnixNano())

	flag.Parse()
	config.LoadConfigFileOrPanic(*_ConfigFile)
	defer SetupTeardown()()

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

func SetupTeardown() func() {
	logrus.Debug("Run service-initialization.")
	SetupRuntimeEnvironment(config.GetConfig())
	return func() {
		logrus.Debug("Run service-cleanup.")
		ClearRuntimeEnvironment(config.GetConfig())
	}
}

func SetupRuntimeEnvironment(conf *config.Config) {
	logger.SetGlobalLogger(conf)
	if len(conf.ServiceMetricsEndpoint) > 0 {
		go func() {
			metrics.Register()
			http.Handle("/metrics", promhttp.Handler())
			logger.GetGlobalLogger().Error(http.ListenAndServe(conf.ServiceMetricsEndpoint, nil))
		}()
	}
	// Add more service initialization here.
	service.Setup{{ServiceNameInCamelCase}}Impl()
}

func ClearRuntimeEnvironment(_ *config.Config) {
	service.Close{{ServiceNameInCamelCase}}Impl()
	// Add more service cleanup here.
}
