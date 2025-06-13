# Prometheus MCP Server

An MCP server for Prometheus

## Prometheus

MCP server requires a Prometheus server

```bash
PORT="9090"
VERS="v3.4.1"

podman run \
--interactive --tty --rm \
--name=prometheus \
--publish=${PORT}:9090/tcp \
quay.io/prometheus/prometheus:${VERS}
```

## MCP

```bash
PORT="9090"

go run ./cmd/server --prometheus=:${PORT}
```

### `tools/list`

```JSON
{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}
```
Yields:
```JSON
{
    "jsonrpc":"2.0",
    "id":1,
    "result":{
        "tools":[
            {
                "annotations":{
                    "readOnlyHint":false,
                    "destructiveHint":true,
                    "idempotentHint":false,
                    "openWorldHint":true
                },
                "description":"Prometheus Metrics",
                "inputSchema":{
                    "properties":{},
                    "type":"object"
                },
                "name":"metrics"
            },{
                "annotations":{
                    "readOnlyHint":false,
                    "destructiveHint":true,
                    "idempotentHint":false,
                    "openWorldHint":true
                },
                "description":"Prometheus Targets",
                "inputSchema":{
                    "properties":{},
                    "type":"object"
                },
                "name":"targets"
            }
        ]
    }
}
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