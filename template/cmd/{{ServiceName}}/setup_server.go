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
	interceptors "{{RepoBase}}/{{RepoGroup}}/{{RepoName}}/internal/service/grpc_interceptors"
)

const (
	_DefaultMaxSendMsgSize         = 8 * 1024 * 1024
	_DefaultMaxRecvMsgSize         = 8 * 1024 * 1024
	_DefaultCliMinPingIntervalTime = 3 * time.Minute
	_DefaultSrvKeepaliveTime       = 5 * time.Minute
	_DefaultSrvKeepaliveTimeout    = 2 * time.Minute
)

func setupGrpcService(_ context.Context, wg *sync.WaitGroup, stopCh chan struct{}) {
	defer wg.Done()

	// Set up a tcp connection to the server.
	l, err := net.Listen("tcp", config.GetConfig().ServiceGrpcEndpoint)
	if err != nil {
		logger.GetGlobalLogger().WithError(err).Fatal("Failed to start grpc service.")
	}

	// gRPC server options, such as TLS, keepalive, etc.
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
			interceptors.RecoverPanicAndReportLatencyUnaryInterceptor,
		)),
	}
	// Create a gRPC server.
	grpcServer := grpc.NewServer(opts...)
	// Register the service.
	proto_gens.Register{{ServiceNameInCamelCase}}Server(grpcServer, service.Get{{ServiceNameInCamelCase}}Impl())
	if config.GetConfig().EnableReflection {
		reflection.Register(grpcServer)
	}

	go func() {
		// Listen on the given address and port.
		if err := grpcServer.Serve(l); err != nil {
			logger.GetGlobalLogger().
				WithError(err).Error("Failed to serve grpc service.")
		}
	}()
	logger.GetGlobalLogger().Infof("gRPC Server started, listening on %s.",
		config.GetConfig().ServiceGrpcEndpoint)
	logger.GetGlobalLogger().Infof("Started {{ServiceNameInCamelCase}} Server 🤘.")

	<-stopCh
	grpcServer.GracefulStop()
	logger.GetGlobalLogger().Warning("Stopped grpc service.")
}
