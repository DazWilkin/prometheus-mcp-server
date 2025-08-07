#!/usr/bin/env bash

# set -x

source .env.test

# To avoid challenges using set vs. unset as a boolean in JQ
# Redefine TLS to be either "true" or "false"
# JQ accepts environment variable as a JSON value because of `--argjson`
if [[ -v TLS ]]
then
  echo "TLS set"
  TLS="true"
else
  echo "TLS unset"
  TLS="false"
fi

  ./inspector.sh \
  | jq \
    --argjson "TLS" "${TLS}" \
    --from-file ./inspector.validate.jq
