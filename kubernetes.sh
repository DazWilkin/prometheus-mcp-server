#!/usr/bin/env bash

# set -x

source .env.test

go tool jsonnet \
--ext-str "NAME=${NAME}" \
--ext-str "GHCR_USERNAME=${GHCR_USERNAME}" \
--ext-str "GHCR_EMAIL=${GHCR_EMAIL}" \
--ext-str "GHCR_TOKEN=${GHCR_TOKEN}" \
--ext-str "SERVER_PORT=${SERVER_PORT}" \
--ext-str "METRIC_PORT=${METRIC_PORT}" \
./kubernetes.jsonnet