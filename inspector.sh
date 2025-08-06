#!/usr/bin/env bash

# set -x

source .env.test

# Generate 32-byte hex token for MCP_PROXY_AUTH_TOKEN
# Either
PROXY_TOKEN=$(openssl rand -hex 32)
# Or
# TOKEN=$(xxd -l 32 -p /dev/urandom | tr -d '\n')

# Store TOKEN in .env.test to be consumed by inspector.test.sh
sed \
--in-place \
--expression="s|PROXY_TOKEN=\"[0-9a-f]\{64\}\"|PROXY_TOKEN=\"${PROXY_TOKEN}\"|g" \
.env.test

go tool jsonnet \
--ext-str "NAME=${INSPECTOR_NAME}" \
--ext-str "TOKEN=${PROXY_TOKEN}" \
--ext-str "TAILNET=${TAILNET}" \
./inspector.jsonnet