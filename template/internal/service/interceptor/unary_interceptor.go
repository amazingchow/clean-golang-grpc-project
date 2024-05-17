package interceptor

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"{{RepoBase}}/{{RepoGroup}}/{{RepoName}}/internal/common/logger"
	"{{RepoBase}}/{{RepoGroup}}/{{RepoName}}/internal/service/common"
)

func RecoverPanicAndReportLatencyUnaryInterceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
	resp interface{}, err error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = status.Error(codes.InvalidArgument, "request header must be provided")
		return
	}

	traceId := uuid.New().String()
	spanId := uuid.New().String()
	if v, ok := md[common.ReqHeaderKeyRequestId]; ok {
		traceId = v[0]
	}
	ctx = context.WithValue(ctx, common.ContextKeyTraceId, traceId)
	ctx = context.WithValue(ctx, common.ContextKeySpanId, spanId)

	start := time.Now()
	defer func() {
		if e := recover(); e != nil {
			logger.GetGlobalLogger().
				WithField(common.LoggerKeyTraceId, traceId).
				WithField(common.LoggerKeySpanId, spanId).
				WithField("stack", string(debug.Stack())).
				WithError(e.(error)).
				Error("recover from internal server panic")
			err = status.Error(codes.Internal, "recover from internal server panic")
		}

		end := time.Now()
		logger.GetGlobalLogger().
			WithField(common.LoggerKeyTraceId, traceId).
			WithField(common.LoggerKeySpanId, spanId).
			WithField("latency", end.Sub(start).Milliseconds()).
			Debug("request latency")
	}()

	return handler(ctx, req)
}
