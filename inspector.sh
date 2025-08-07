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

# To avoid challenges using set vs. unset as a boolean in Jsonnet
# Redefine TLS to be either "T" (set) or "F" (unset)
# Jsonnet assigned a boolean based on whether TLS=="T" or not
# These values are then used by the Inspector's Jsonnet deployment to determine the path
# And by the Inspector's JQ validation script
TLS=$(\
  if [[ -v TLS ]]
  then
    echo "T"
  else
    echo "F"
  fi)

go tool jsonnet \
--ext-str "NAME=${INSPECTOR_NAME}" \
--ext-str "TOKEN=${PROXY_TOKEN}" \
--ext-str "TAILNET=${TAILNET}" \
--ext-str "TLS=${TLS}" \
./inspector.jsonnet