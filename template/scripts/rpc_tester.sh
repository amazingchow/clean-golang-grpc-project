#!/usr/bin/env bash

# 遇到执行出错，直接终止脚本的执行
set -o errexit

function logger_print
{
    local prefix="[$(date +%Y/%m/%d\ %H:%M:%S)]"
    echo "${prefix}$@" >&2
}

function test_rpc_methods
{

    grpcurl \
        -rpc-header x-request-id:73338239da584998aca91639651334fa -d @ -plaintext \
        localhost:{{ServiceGrpcPort}} {{ServiceNameInUnderScoreCase}}.{{ServiceNameInCamelCase}}/Ping << EOM
{
}
EOM

}

function run
{
    # go install github.com/fullstorydev/grpcurl/cmd/grpcurl@v1.8.9 for go1.18
    # go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest for go1.19 or newer
    grpcurl -plaintext localhost:{{ServiceGrpcPort}} list {{ServiceNameInUnderScoreCase}}.{{ServiceNameInCamelCase}}
    test_rpc_methods
}

run $@
