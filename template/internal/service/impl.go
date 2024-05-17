package service

import (
	"context"

	"{{RepoBase}}/{{RepoGroup}}/{{RepoName}}/internal/common/logger"
	"{{RepoBase}}/{{RepoGroup}}/{{RepoName}}/internal/proto_gens"
	"{{RepoBase}}/{{RepoGroup}}/{{RepoName}}/internal/service/common"
)

func (impl *{{ServiceNameInCamelCase}}Impl) Ping(
	ctx context.Context, req *proto_gens.PingRequest) (
	resp *proto_gens.PongResponse, err error) {

	_logger := logger.GetGlobalLogger().
		WithField("method", "Ping").
		WithField(common.LoggerKeyTraceId, ctx.Value(common.ContextKeyTraceId).(string)).
		WithField(common.LoggerKeySpanId, ctx.Value(common.ContextKeySpanId).(string))
	_ = _logger

	resp = &proto_gens.PongResponse{}
	return
}
