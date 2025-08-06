#!/usr/bin/env bash

# set -x

source .env.test

INSPECTOR_PROXY_URL="https://${INSPECTOR_NAME}-proxy.${TAILNET}"

# Config
curl \
--silent \
--header "X-MCP-Proxy-Auth: Bearer ${PROXY_TOKEN}" \
${INSPECTOR_PROXY_URL}/config \
| jq -r .

# Status
curl \
--silent \
${INSPECTOR_PROXY_URL}/health \
| jq -r .