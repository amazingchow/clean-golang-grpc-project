FROM golang:1.21.10 as builder
WORKDIR /build
ADD . /build/
RUN make build

FROM alpine:3.18
WORKDIR /app
RUN apk add libc6-compat
COPY --from=builder /build/{{ServiceName}} /app
RUN mkdir -p /app/config /app/logs /app/persistent /app/locks /app/shares
EXPOSE {{ServiceGrpcPort}}
EXPOSE {{ServiceHttpPort}}
CMD [ \
    "/app/{{ServiceName}}", \
    "-conf", \
    "/config/{{ServiceName}}.json" \
]
