#!/usr/bin/env bash

# set -x

source .env.test

SERVER="https://${NAME}.${TAILNET}"

# Generalizes invoking a JSON-RPC method
# Expects: METHOD, NAME, ARGUMENTS
# METHOD: list, call
# NAME: alertmanagers, alerts, exemplars, metrics, ping, query, query_range, rules, status_tsdb, targets
# ARGUMENTS: JSON object (e.g. {"query":"up","time":"2024-01-01T00:00:00Z"})
# Returns: JSON-RPC response
# Uses test.http.jsonnet to construct JSON-RPC request
# Uses: SERVER
jsonrpc() {
  local METHOD="${1}"
  local NAME="${2}"
  local ARGUMENTS="${3}"

  # Prefix method with "tools/"
  printf "\nTesting MCP tools/%s '%s'\n" "${METHOD}" "${NAME}"

  # Construct request
  local DATA
  DATA=$(\
    go tool jsonnet \
    --ext-str "METHOD=${METHOD}" \
    --ext-str "NAME=${NAME}" \
    --ext-str "ARGUMENTS=${ARGUMENTS}" \
    ./test.http.jsonnet)

  echo "${DATA}" | jq -r .

  # Send request to Prometheus MCP server
  local RESPONSE
  RESPONSE=$(
    curl \
    --silent \
    --show-error \
    --request POST \
    --header "Content-Type: application/json" \
    --data "${DATA}" \
    "${SERVER}/mcp")

    echo "${RESPONSE}" | jq -r .
}

# Expects Prometheus server
# Must use group ({}) not subshell (()) to be able to terminate
{
  HEALTH="${PROMETHEUS_URL}/-/healthy"
  CODE=$(\
    curl \
    --silent \
    --get \
    "${HEALTH}" \
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
    "${SERVER}/mcp" \
    --output /dev/null \
    --write-out '%{response_code}')
  
  if [[ "${CODE}" != "200" ]]
  then
    printf "Unable to 'Ping' Prometheus MCP server (%s) got %s; want: 200\n" "${SERVER}" "${CODE}"
    exit 1
  fi
}

# `tools/list`
list() {
  local METHOD="list"
  local NAME=""
  local ARGUMENTS=""
  jsonrpc "${METHOD}" "${NAME}" "${ARGUMENTS}"
}

# `tools/call`(Foo)
call() {
  local METHOD="call"
  local NAME="${1}"
  local ARGUMENTS="${2}"
  jsonrpc "${METHOD}" "${NAME}" "${ARGUMENTS}"
}

# # `tools/call` (Alertmanagers)
# call_alertmanagers() {
#   local METHOD="call"
#   local NAME="alertmanagers"
#   local ARGUMENTS=""
#   jsonrpc "${METHOD}" "${NAME}" "${ARGUMENTS}"
# }

# # `tools/call` (Alerts)
# call_alerts() {
#   local METHOD="call"
#   local NAME="alerts"
#   local ARGUMENTS=""
#   jsonrpc "${METHOD}" "${NAME}" "${ARGUMENTS}"
# }

# # `tools/call` (Exemplars)
# call_exemplars() {
#   local METHOD="call"
#   local NAME="exemplars"
#   local ARGUMENTS
#   ARGUMENTS=$(\
#     go tool jsonnet \
#     --ext-str "QUERY=${QUERY}" \
#     --ext-str "START=${START}" \
#     --ext-str "END=${END}" \
#     ./testjson/exemplars.jsonnet \
#     | jq -r .)
#   jsonrpc "${METHOD}" "${NAME}" "${ARGUMENTS}"
# }

# # `tools/call` (Metrics)
# call_metrics() {
#   METHOD="call"
#   NAME="metrics"
#   ARGUMENTS=""
#   jsonrpc "${METHOD}" "${NAME}" "${ARGUMENTS}"
# }

# # `tools/call` (Ping)
# call_ping() {
#   local METHOD="call"
#   local NAME="ping"
#   local ARGUMENTS=""
#   jsonrpc "${METHOD}" "${NAME}" "${ARGUMENTS}"  
# }

# # `tools/call` (Query: instant)
# call_query() {
#   local METHOD="call"
#   local NAME="query"
#   local ARGUMENTS
#   ARGUMENTS=$(\
#     go tool jsonnet \
#     --ext-str "QUERY=${QUERY}" \
#     --ext-str "START=${START}" \
#     --ext-str "END=${END}" \
#     --ext-str "STEP=${STEP}" \
#     ./testjson/query.jsonnet \
#     | jq -r .)
#   jsonrpc "${METHOD}" "${NAME}" "${ARGUMENTS}"
# }

# # `tools/call` (Query: range)
# # TODO(dazwilkin) Add "time" (and optional "timeout","limit")
# call_query_range() {
#   local METHOD="call"
#   local NAME="query_range"
#   local ARGUMENTS
#   ARGUMENTS=$(\
#   go tool jsonnet \
#   --ext-str "QUERY=${QUERY}" \
#   --ext-str "START=${START}" \
#   --ext-str "END=${END}" \
#   --ext-str "STEP=${STEP}" \
#   ./testjson/query_range.jsonnet \
#   | jq -r .)
#   jsonrpc "${METHOD}" "${NAME}" "${ARGUMENTS}"
# }

# # `tools/call` (Rules)
# call_rules() {
#   local METHOD="call"
#   local NAME="rules"
#   local ARGUMENTS=""
#   jsonrpc "${METHOD}" "${NAME}" "${ARGUMENTS}"
# }

# # `tools/call` (Status TSDB)
# call_status_tsdb() {
#   local METHOD="call"
#   local NAME="status_tsdb"
#   local ARGUMENTS=""
#   jsonrpc "${METHOD}" "${NAME}" "${ARGUMENTS}"
# }

# # `tools/call` (Targets)
# call_targets() {
#   local METHOD="call"
#   local NAME="targets"
#   local ARGUMENTS=""
#   jsonrpc "${METHOD}" "${NAME}" "${ARGUMENTS}"
# }

# Use Jsonnet to generate the JSON arguments for those tools that need them
JSON_EXEMPLARS=$(\
  go tool jsonnet \
  --ext-str "QUERY=${QUERY}" \
  --ext-str "START=${START}" \
  --ext-str "END=${END}" \
  ./testjson/exemplars.jsonnet \
  | jq -r .)

JSON_QUERY_RANGE=$(\
  go tool jsonnet \
  --ext-str "QUERY=${QUERY}" \
  --ext-str "START=${START}" \
  --ext-str "END=${END}" \
  --ext-str "STEP=${STEP}" \
  ./testjson/query_range.jsonnet \
  | jq -r .)

# Invoke the JSON-RPCs
list
call "alertmanagers" ""
call "alerts" ""
call "exemplars" "${JSON_EXEMPLARS}"
call "metrics" ""
call "ping" ""
call "query" ""
call "query_range" "${JSON_QUERY_RANGE}"
call "rules" ""
call "status_tsdb" ""
call "targets" ""

# Grep Prometheus metrics
(
  FILTER='
   .data.result[]
   |"\(.metric.__name__){tool=\(.metric.tool)} \(.value[1])"'
  curl \
  --silent \
  --show-error \
  --get \
  --data-urlencode "query={__name__=~\"mcp_prometheus_(total|error)\"}" \
  ${PROMETHEUS_URL}/api/v1/query \
  | jq -r "${FILTER}"
)