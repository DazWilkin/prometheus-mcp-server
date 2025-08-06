#!/usr/bin/env bash

# set -x

source .env.test

# Override NAME
NAME="inspector"

# Generate 32-byte hex token for MCP_PROXY_AUTH_TOKEN
# Either
TOKEN=$(openssl rand -hex 32)
# Or
# TOKEN=$(xxd -l 32 -p /dev/urandom | tr -d '\n')

go tool jsonnet \
--ext-str "NAME=${NAME}" \
--ext-str "TOKEN=${TOKEN}" \
--ext-str "TAILNET=${TAILNET}" \
./inspector.jsonnet