#!/usr/bin/env bash

# 遇到执行出错，直接终止脚本的执行
set -o errexit

logger_print()
{
    local prefix="[$(date +%Y/%m/%d\ %H:%M:%S)]"
    echo "${prefix}$@" >&2
}

function benchmark_rpc_methods
{
    # go install github.com/bojand/ghz/cmd/ghz@latest
    ghz --insecure \
        --async \
        --rps 20 \
        --concurrency 20 \
        --total 200 \
        --timeout 30s \
        --proto ../protos/{{ServiceNameInUnderScoreCase}}.proto \
        --call {{ServiceNameInUnderScoreCase}}.{{ServiceNameInCamelCase}}.Ping \
        --metadata '{"x-request-id": "73338239da584998aca91639651334fa"}' \
        --data '{}' \
        localhost:{{ServiceGrpcPort}}
}

function run
{
    grpcurl -plaintext localhost:{{ServiceGrpcPort}} list {{ServiceNameInUnderScoreCase}}.{{ServiceNameInCamelCase}}
    benchmark_rpc_methods
}

run $@
