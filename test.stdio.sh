#!/usr/bin/env bash

set -x

source .env.test

# `tools/list`
(
    JSON='{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}'
    echo ${JSON} \
    | go run ./cmd/server \
      --metric.addr="" \
      --server.addr="" \
    | jq -r .
)

# `tools/call` (Alerts)
(
    JSON='{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"alerts","arguments":{}}}'
    echo ${JSON} \
    | go run ./cmd/server \
      --metric.addr="" \
      --server.addr="" \
    | jq -r .    
)

# `tools/call` (Metrics)
(
    JSON='{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"metrics","arguments":{}}}'
        echo ${JSON} \
    | go run ./cmd/server \
      --metric.addr="" \
      --server.addr="" \
    | jq -r .
)