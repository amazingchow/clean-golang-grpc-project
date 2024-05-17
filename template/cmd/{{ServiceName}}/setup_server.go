package main

import (
	"context"
	"net"
	"sync"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"{{RepoBase}}/{{RepoGroup}}/{{RepoName}}/internal/common/config"
	"{{RepoBase}}/{{RepoGroup}}/{{RepoName}}/internal/common/logger"
	"{{RepoBase}}/{{RepoGroup}}/{{RepoName}}/internal/proto_gens"
	"{{RepoBase}}/{{RepoGroup}}/{{RepoName}}/internal/service"
	"{{RepoBase}}/{{RepoGroup}}/{{RepoName}}/internal/service/interceptor"
)

const (
	_DefaultMaxSendMsgSize         = 8 * 1024 * 1024
	_DefaultMaxRecvMsgSize         = 8 * 1024 * 1024
	_DefaultCliMinPingIntervalTime = 3 * time.Minute
	_DefaultSrvKeepaliveTime       = 5 * time.Minute
	_DefaultSrvKeepaliveTimeout    = 2 * time.Minute
)

func setupGrpcService(ctx context.Context, wg *sync.WaitGroup, stopCh chan struct{}) {
	defer wg.Done()

	l, err := net.Listen("tcp", config.GetConfig().ServiceGrpcEndpoint)
	if err != nil {
		logger.GetGlobalLogger().WithError(err).Fatal("Failed to start grpc service.")
	}

	opts := []grpc.ServerOption{
		grpc.MaxSendMsgSize(_DefaultMaxSendMsgSize),
		grpc.MaxRecvMsgSize(_DefaultMaxRecvMsgSize),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             _DefaultCliMinPingIntervalTime,
			PermitWithoutStream: false,
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    _DefaultSrvKeepaliveTime,
			Timeout: _DefaultSrvKeepaliveTimeout,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			interceptor.RecoverPanicAndReportLatencyUnaryInterceptor,
		)),
	}
	grpcServer := grpc.NewServer(opts...)
	proto_gens.Register{{ServiceNameInCamelCase}}Server(grpcServer, service.Get{{ServiceNameInCamelCase}}Impl())
	reflection.Register(grpcServer)

	logger.GetGlobalLogger().Infof("grpc service is running @\x1b[1;31m%s\x1b[0m",
		config.GetConfig().ServiceHttpEndpoint)
	go func() {
		if err := grpcServer.Serve(l); err != nil {
			logger.GetGlobalLogger().
				WithError(err).Error("Failed to serve grpc service.")
		}
	}()

	<-stopCh
	grpcServer.GracefulStop()
	logger.GetGlobalLogger().Info("Stop grpc service.")
}
