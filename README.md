# Prometheus MCP Server

[![Go Report Card](https://goreportcard.com/badge/github.com/DazWilkin/prometheus-mcp-server)](https://goreportcard.com/report/github.com/DazWilkin/prometheus-mcp-server)
[![Go Reference](https://pkg.go.dev/badge/github.com/DazWilkin/prometheus-mcp-server.svg)](https://pkg.go.dev/github.com/DazWilkin/prometheus-mcp-server)
[![build](https://github.com/DazWilkin/prometheus-mcp-server/actions/workflows/build.yml/badge.svg)](https://github.com/DazWilkin/prometheus-mcp-server/actions/workflows/build.yml)

An MCP server for [Prometheus](https://prometheus.io)

Very much a work in progress: **not tested** in a MCP client host

+ Implements MCP `stdio` and HTTP streamable
+ Implements Prometheus [HTTP API](https://prometheus.io/docs/prometheus/latest/querying/api/) methods:
  + [Instant queries](https://prometheus.io/docs/prometheus/latest/querying/api/#instant-queries)
  + [Range queries](https://prometheus.io/docs/prometheus/latest/querying/api/#range-queries)
  + [List Alertmanagers](https://prometheus.io/docs/prometheus/latest/querying/api/#alertmanagers)
  + [List Alerts](https://prometheus.io/docs/prometheus/latest/querying/api/#alerts)
  + [List Exemplars](https://prometheus.io/docs/prometheus/latest/querying/api/#querying-exemplars)
  + [List Metrics](https://prometheus.io/docs/prometheus/latest/querying/api/#querying-label-values)
  + [List Rules](https://prometheus.io/docs/prometheus/latest/querying/api/#rules)
  + [List Series](https://prometheus.io/docs/prometheus/latest/querying/api/#finding-series-by-label-matchers)
  + [List Status TSDB](https://prometheus.io/docs/prometheus/latest/querying/api/#tsdb-stats)
  + [List Targets](https://prometheus.io/docs/prometheus/latest/querying/api/#targets)
+ Implements [Prometheus Management API](https://prometheus.io/docs/prometheus/latest/management_api/)
  + [Health check](https://prometheus.io/docs/prometheus/latest/management_api/) 

## Limitations

A non-exhaustive list:

+ Prometheus API: Only partially implements

## Visual Studio Code GitHub Copilot agent

Functions as an MCP host and can be configured to interact with `prometheus-mcp-server`.

See [Use MCP servers in VS Code (Preview)](https://code.visualstudio.com/docs/copilot/chat/mcp-servers) to configure this.

`.vscode/mcp.json`:
```JSON
{
    "servers": {
        "prometheus-mcp-server": {
            "url": "http://localhost:7777/mcp"
        }
    }
}
```

Then open "Chat" and ensure it's configured for "Agent" (bottom left hand corner of the chat window).

The tools icon should list the `prometheus-mcp-server` tools and you can:

```console
#ping
```
```console
#metrics
```
```console
#query up{job="prometheus"}
```

> **NOTE** There's an issue with `metrics`. The server responds correctly and chat responds with "Here is a list of all available Prometheus metrics currently exposed by your server. If you need details or want to query a specific metric, let me know which one you're interested in!" but the metrics aren't listed. You must ask the agent to "give me the metrics filtered by ..." for example "filtered by promhttp_metric_handler_requests_total".

For further ways to interact with the MCP server, try "CTRL-SHIFT-P" and e.g. `MCP: List Servers`, select `prometheus-mcp-server` and then select one of the commands.

## Prometheus

MCP server requires a Prometheus server

```bash
PROM="9090"
VERS="v3.4.1"

podman run \
--interactive --tty --rm \
--name=prometheus \
--publish=${PROM}:9090/tcp \
quay.io/prometheus/prometheus:${VERS}
```

## MCP

### `stdio`

Configured if `--server.addr==""`

Pipe the `stdout` through `jq`:

```bash
 # Upstream Prometheus server
PROMETHEUS_URL="http://localhost:9090"

go run \
./cmd/server \
--prometheus="${PROMETHEUS_URL}" \
| jq -r .
```
See [`tools/list`](#toolslist) for container example.

### HTTP streamable

Configured if `--server.addr!=""` defaults to `:7777` 

`--server.path` defaults to `/mcp`

Currently configured to be stateless because I'm unsure how to provide session IDs.

```bash
# Prometheus MCP server
SERVER_HOST="0.0.0.0"
SERVER_PORT="7777"
SERVER_ADDR="${SERVER_HOST}:${SERVER_PORT}"
SERVER_PATH="/mcp"

# Prometheus MCP metrics exporter
METRIC_HOST="0.0.0.0"
METRIC_PORT="8080"
METRIC_ADDR="${METRIC_HOST}:${METRIC_PORT}"
METRIC_PATH="/metrics"

 # Upstream Prometheus server
PROMETHEUS_URL="http://localhost:9090"

go run \
./cmd/server \
--server.addr="${SERVER_ADDR}" \
--server.path="${SERVER_PATH}" \
--metric.addr="${METRIC_ADDR}" \
--metric.path="${METRIC_PATH}" \
--prometheus="${PROMETHEUS_URL}"
```
Or:
```bash
IMAGE="ghcr.io/dazwilkin/prometheus-mcp-server:9c2f0987302f23f533c31e2cc8eafe02fa477232"

# Prometheus MCP server
SERVER_HOST="0.0.0.0"
SERVER_PORT="7777"
SERVER_ADDR="${SERVER_HOST}:${SERVER_PORT}"
SERVER_PATH="/mcp"

# Prometheus MCP metrics exporter
METRIC_HOST="0.0.0.0"
METRIC_PORT="8080"
METRIC_ADDR="${METRIC_HOST}:${METRIC_PORT}"
METRIC_PATH="/metrics"

 # Upstream Prometheus server
PROMETHEUS_URL="http://localhost:9090"

# Uses --net=host to access upstream Prometheus
# --publish= provided for documentation
podman run \
--interactive --tty --rm \
--net=host \
--publish="${SERVER_PORT}:${SERVER_PORT}/tcp" \
--publish="${METRIC_PORT}:${METRIC_PORT}/tcp" \
${IMAGE} \
--server.addr="${SERVER_ADDR}" \
--server.path="${SERVER_PATH}" \
--metric.addr="${METRIC_ADDR}" \
--metric.path="${METRIC_PATH}" \
--prometheus="${PROMETHEUS_URL}"
```

### Prometheus metrics exporter

Configured if `--metric.addr!=""` defaults to `:8080`

`--metric.path` defaults to `/metrics`

## API

For `stdio` copy-paste examples below into server's stdin.

See [`test.stdio.sh`](./test.stdio.sh)

For HTTP streamable:

```bash
curl \
--request POST \
--header "Content-Type: application/json"
--data '{json}' \
http://{server.addr}/{server.path}
```

See [`test.http.sh`](./test.http.sh)

### `tools/list`

```JSON
{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}
```
Yields [`tools.list.json`](./tools.list.json)

You may also pipe MCP (JSON-RPC) messages into the `prometheus-mcp-server` container:

```bash
IMAGE="ghcr.io/dazwilkin/prometheus-mcp-server:9c2f0987302f23f533c31e2cc8eafe02fa477232"

 # Upstream Prometheus server
PROMETHEUS_URL="http://localhost:9090"

JSON='{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}'

echo ${JSON} \
| podman run \
  --interactive --rm \
  --net=host \
  --name=prometheus-mcp-server \
  ${IMAGE} \
  --server.addr="" \
  --metric.addr="" \
  --prometheus="${PROMETHEUS_URL}" \
| jq -r .
```

### `tools/call`

### `alerts`

```JSON
{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"alerts","arguments":{}}}
```
Yields:
```JSON
{"jsonrpc":"2.0","id":2,"result":{"content":[{"type":"text","text":"{\"alerts\":[]}"}]}}
```

#### `metrics`

```JSON
{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"metrics","arguments":{}}}
```
Yields:
```JSON
{
    "jsonrpc":"2.0",
    "id":2,
    "result":{
        "content":[
            {
                "type":"text",
                "text":"[\"go_gc_cycles_automatic_gc_cycles_total\",,...,\"up\"]"
            }
        ]
    }
}
```

#### `query`

```JSON
{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"query","arguments":{"query":"up{job=\"prometheus\"}"}}}
```
Yields:
```JSON
{"jsonrpc":"2.0","id":2,"result":{"content":[{"type":"text","text":"[{\"metric\":{\"__name__\":\"up\",\"app\":\"prometheus\",\"instance\":\"localhost:9090\",\"job\":\"prometheus\"},\"value\":[1749834703.085,\"1\"]}]"}]}}
```

With `time`:

```JSON
{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"query","arguments":{"query":"up{job=\"prometheus\"}","time":"2025-06-13T10:10:00-07:00"}}}
```
Yields:
```JSON
{"jsonrpc":"2.0","id":2,"result":{"content":[{"type":"text","text":"[{\"metric\":{\"__name__\":\"up\",\"app\":\"prometheus\",\"instance\":\"localhost:9090\",\"job\":\"prometheus\"},\"value\":[1749834600,\"1\"]}]"}]}}
```
With `time`, `timeout`
```JSON
{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"query","arguments":{"query":"up{job=\"prometheus\"}","time":"2025-06-13T10:10:00-07:00","timeout":"10s"}}}
```

#### `query_range`

```JSON
{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"query_range","arguments":{"query":"up{job=\"prometheus\"}","start":"2025-06-13T10:00:00-07:00","end":"2025-06-13T11:00:00-07:00","step":"5m"}}}
```

#### `rules`

```JSON
{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"rules","arguments":{}}}
```
Yields:
```JSON
{"jsonrpc":"2.0","id":3,"result":{"content":[{"type":"text","text":"{\"groups\":[]}"}]}}
```

#### `series`

```JSON
{"jsonrpc": "2.0","id": 2,"method":"tools/call","params":{"name":"series","arguments":{"match[]":["up{}","up{job=\"prometheus\"}"]}}}
```
Yields:
```JSON
{"jsonrpc":"2.0","id":2,"result":{"content":[{"type":"text","text":"[{\"__name__\":\"up\",\"app\":\"prometheus\",\"instance\":\"localhost:9090\",\"job\":\"prometheus\"}]"}]}}
```

#### `targets`

```JSON
{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"targets","arguments":{}}}
```
Yields:
```JSON
{
    "jsonrpc":"2.0",
    "id":2,
    "result":{
        "content":[
            {
                "type":"text",
                "text":"{\"activeTargets\":[{\"discoveredLabels\":{\"__address__\":\"localhost:9090\",\"__metrics_path__\":\"/metrics\",\"__scheme__\":\"http\",\"__scrape_interval__\":\"15s\",\"__scrape_timeout__\":\"10s\",\"app\":\"prometheus\",\"job\":\"prometheus\"},\"labels\":{\"app\":\"prometheus\",\"instance\":\"localhost:9090\",\"job\":\"prometheus\"},\"scrapePool\":\"prometheus\",\"scrapeUrl\":\"http://localhost:9090/metrics\",\"globalUrl\":\"http://6dd53f16a42c:9090/metrics\",\"lastError\":\"\",\"lastScrape\":\"2025-06-13T16:14:29.097537931Z\",\"lastScrapeDuration\":0.003549414,\"health\":\"up\"}],\"droppedTargets\":[]}"
            }
        ]
    }
}
```

## Exporter

The MCP server exports Prometheus metrics

The metrics are prefix `mcp_prometheus_`

|Name|Type|Description|
|----|----|-----------|
|`build`|Counter|A metric with a constant '1' value labels by build|start time, git commit, OS and Go versions|
|`total`|Counter|Total number of successful MCP tool invocations|
|`error`|Counter|Total number of unsuccessful MCP tool invocations|

## Sigstore
`prometheus-mcp-server` container images are being signed by Sigstore and may be verified:

```bash
go tool cosign verify \
--key=./cosign.pub \
ghcr.io/dazwilkin/prometheus-mcp-server:9c2f0987302f23f533c31e2cc8eafe02fa477232
```

> **Note**

`cosign.pub` may be downloaded [here](./cosign.pub)

`cosign` is included as a `go.mod` tool.

<hr/>
<br/>
<a href="https://www.buymeacoffee.com/dazwilkin" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/default-orange.png" alt="Buy Me A Coffee" height="41" width="174"></a>
