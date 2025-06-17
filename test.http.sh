#!/usr/bin/env bash

set -x

source .env.test

SERVER="${SERVER_HOST}:${SERVER_PORT}/${SERVER_PATH}"
METRIC="${METRIC_HOST}:${METRIC_PORT}/${METRIC_PATH}"

QUERY="up{job='prometheus'}"

# E.g. 2025-06-14
DATE="$(date +%Y-%m-%d)"
ZONE="-07:00"
START="${DATE}T00:00:00${ZONE}"
END="${DATE}T23:59:59${ZONE}"
STEP="1h"

TIMEOUT="15s"
LIMIT="10"

# Expects Prometheus server
# Must use group ({}) not subshell (()) to be able to terminate
{
    PROMETHEUS="http://localhost:9090"
    HEALTH="${PROMETHEUS}/-/healthy"
    CODE=$(\
      curl \
      --silent \
      --get \
      ${HEALTH} \
      --output /dev/null \
      --write-out '%{response_code}')

    if [[ "${CODE}" != "200" ]]
    then
      printf "Unable to get Prometheus Health endpoint (%s) got: %s; want: 200\n" "${HEALTH}" "${CODE}"
      exit 1
    fi
}

# Expects Prometheus MCP server
# MCP "ping"
# Must use group ({}) not subshell (()) to be able to terminate
{
    JSON='{"jsonrpc": "2.0","method": "ping","id": 1}'
    CODE=$(\
      curl \
      --silent \
      --request POST \
      --header "Content-Type: application/json" \
      --data "${JSON}" \
      http://${SERVER} \
      --output /dev/null \
      --write-out '%{response_code}')
    
    if [[ "${CODE}" != "200" ]]
    then
      printf "Unable to 'Ping' Prometheus MCP server (%s) got %s; want: 200\n" "${SERVER}" "${CODE}"
      exit 1
    fi
}

# `tools/list`
(
    JSON='{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}'
    curl \
    --request POST \
    --header "Content-Type: application/json" \
    --data "${JSON}" \
    http://${SERVER}
)

# `tools/call` (Alertmanagers)
(
    JSON='{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"alertmanagers","arguments":{}}}'
    curl \
    --request POST \
    --header "Content-Type: application/json" \
    --data "${JSON}" \
    http://${SERVER}
)

# `tools/call` (Alerts)
(
    JSON='{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"alerts","arguments":{}}}'
    curl \
    --request POST \
    --header "Content-Type: application/json" \
    --data "${JSON}" \
    http://${SERVER}
)

# `tools/call` (Exemplars)
{
    JSON="{
      \"jsonrpc\":\"2.0\",
      \"id\":2,
      \"method\":\"tools/call\",
      \"params\":{
        \"name\":\"exemplars\",
        \"arguments\":{
          \"query\":\"${QUERY}\",
          \"start\":\"${START}\",
          \"end\":\"${END}\"
        }
      }
    }"
    echo ${JSON} | jq -r .

    curl \
    --request POST \
    --header "Content-Type: application/json" \
    --data "${JSON}" \
    http://${SERVER}
}

# `tools/call` (Metrics)
(
    JSON='{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"metrics","arguments":{}}}'
    curl \
    --request POST \
    --header "Content-Type: application/json" \
    --data "${JSON}" \
    http://${SERVER}
)

# `tools/call` (Ping)
(
    JSON='{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"ping","arguments":{}}}'
    curl \
    --request POST \
    --header "Content-Type: application/json" \
    --data "${JSON}" \
    http://${SERVER}
)

# `tools/call` (Query: instant)
# TODO(dazwilkin) Add "time" (and optional "timeout","limit")
(
    # TODO(dazwilkin) Parameterize JSON to account for dates
    JSON="{
      \"jsonrpc\":\"2.0\",
      \"id\":2,
      \"method\":\"tools/call\",
      \"params\":{
        \"name\":\"query\",
        \"arguments\":{
          \"query\":\"${QUERY}\",
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
    http://${SERVER}
)

# `tools/call` (Query: range)
# TODO(dazwilkin) Add "time" (and optional "timeout","limit")
(
    JSON="{
      \"jsonrpc\":\"2.0\",
      \"id\":2,
      \"method\":\"tools/call\",
      \"params\":{
        \"name\":\"query_range\",
        \"arguments\":{
          \"query\":\"${QUERY}\",
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
    http://${SERVER}
)

# `tools/call` (Rules)
(
    JSON='{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"rules","arguments":{}}}'
    curl \
    --request POST \
    --header "Content-Type: application/json" \
    --data "${JSON}" \
    http://${SERVER}
)

# `tools/call` (Status TSDB)
(
    JSON='{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"status_tsdb","arguments":{}}}'
    curl \
    --request POST \
    --header "Content-Type: application/json" \
    --data "${JSON}" \
    http://${SERVER}
)

# `tools/call` (Targets)
(
    JSON='{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"targets","arguments":{}}}'
    curl \
    --request POST \
    --header "Content-Type: application/json" \
    --data "${JSON}" \
    http://${SERVER}
)

# Grep Prometheus metrics
(
    curl \
    --silent \
    --get \
    http://${METRIC} | awk '/^mcp_prometheus_/'
)