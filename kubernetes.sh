#!/usr/bin/env bash

set -x

source .env.test

jsonnet \
--ext-str "GHCR_USERNAME=${GHCR_USERNAME}" \
--ext-str "GHCR_EMAIL=${GHCR_EMAIL}" \
--ext-str "GHCR_TOKEN=${GHCR_TOKEN}" \
--ext-str "SERVER_HOST=${SERVER_HOST}" \
--ext-str "SERVER_PORT=${SERVER_PORT}" \
--ext-str "METRIC_HOST=${METRIC_HOST}" \
--ext-str "METRIC_PORT=${METRIC_PORT}" \
--ext-str "PROMETHEUS_URL=${PROMETHEUS_URL}" \
./kubernetes.jsonnet