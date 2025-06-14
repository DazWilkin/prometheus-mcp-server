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

# `tools/call` (Alertmanagers)
(
    JSON='{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"alertmanagers","arguments":{}}}'
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

# `tools/call` (Ping)
(
    JSON='{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"ping","arguments":{}}}'
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
    # E.g. 2025-06-14
    DATE="$(date +%Y-%m-%d)"
    ZONE="-07:00"
    START="${DATE}T00:00:00${ZONE}"
    END="${DATE}T23:59:59${ZONE}"
    STEP="1h"
    JSON="{
      \"jsonrpc\":\"2.0\",
      \"id\":2,
      \"method\":\"tools/call\",
      \"params\":{
        \"name\":\"query_range\",
        \"arguments\":{
          \"query\":\"up{job='prometheus'}\",
          \"start\":\"${START}\",
          \"end\":\"${END}\",
          \"step\":\"${STEP}\"
        }
      }
    }"
    echo ${JSON} | jq -r .

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

# `tools/call` (Status TSDB)
(
    JSON='{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"status_tsdb","arguments":{}}}'
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