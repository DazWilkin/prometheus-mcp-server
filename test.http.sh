#!/usr/bin/env bash

set -x

source .env.test

# `tools/list`
(
    JSON='{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}'
    curl \
    --request POST \
    --header "Content-Type: application/json" \
    --data "${JSON}" \
    http://${SERVER_ADDR}/${SERVER_PATH}
)

# `tools/call` (Alerts)
(
    JSON='{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"alerts","arguments":{}}}'
    curl \
    --request POST \
    --header "Content-Type: application/json" \
    --data "${JSON}" \
    http://${SERVER_ADDR}/${SERVER_PATH}
)

# `tools/call` (Metrics)
(
    JSON='{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"metrics","arguments":{}}}'
    curl \
    --request POST \
    --header "Content-Type: application/json" \
    --data "${JSON}" \
    http://${SERVER_ADDR}/${SERVER_PATH}
)

# `tools/call` (Query: instant)
# TODO(dazwilkin) Add "time" (and optional "timeout","limit")
(
    # TODO(dazwilkin) Parameterize JSON to account for dates
    JSON='{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"query","arguments":{"query":"up{job=\"prometheus\"}"}}}'
    curl \
    --request POST \
    --header "Content-Type: application/json" \
    --data "${JSON}" \
    http://${SERVER_ADDR}/${SERVER_PATH}
)

# `tools/call` (Query: range)
# TODO(dazwilkin) Add "time" (and optional "timeout","limit")
(
    # TODO(dazwilkin) Parameterize JSON to account for dates
    JSON='{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"query_range","arguments":{"query":"up{job=\"prometheus\"}","start":"2025-06-13T10:00:00-07:00","end":"2025-06-13T11:00:00-07:00","step":"5m"}}}'
    curl \
    --request POST \
    --header "Content-Type: application/json" \
    --data "${JSON}" \
    http://${SERVER_ADDR}/${SERVER_PATH}
)

# `tools/call` (Rules)
(
    JSON='{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"rules","arguments":{}}}'
    curl \
    --request POST \
    --header "Content-Type: application/json" \
    --data "${JSON}" \
    http://${SERVER_ADDR}/${SERVER_PATH}
)

# `tools/call` (Targets)
(
    JSON='{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"targets","arguments":{}}}'
    curl \
    --request POST \
    --header "Content-Type: application/json" \
    --data "${JSON}" \
    http://${SERVER_ADDR}/${SERVER_PATH}
)

# Grep Prometheus metrics
(
    curl \
    --silent \
    --get \
    http://${METRIC_ADDR}/${METRIC_PATH} | awk '/^mcp_prometheus_/'
)