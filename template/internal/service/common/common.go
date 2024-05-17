package common

type ContextKey string

const (
	ReqHeaderKeyRequestId string = "x-request-id"
)

const (
	ContextKeyTraceId ContextKey = "ctx-key-trace-id"
	ContextKeySpanId  ContextKey = "ctx-key-span-id"
)

const (
	LoggerKeyTraceId string = "x-trace-id"
	LoggerKeySpanId  string = "x-span-id"
)
