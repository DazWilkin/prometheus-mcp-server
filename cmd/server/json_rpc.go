package main

const (
	example string = `{
	"jsonrpc":"2.0",
	"id":2,
	"method":"tools/call",
	"params":{
		"name":"query",
		"arguments":{
			"query":"up{job=\"prometheus\"}",
			"time":"2025-06-13T10:10:00-07:00"
		}
	}
}`
)

type Arguments map[string]any

type Message struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
}
type Params struct {
	Name      string    `json:"name"`
	Arguments Arguments `json:"arguments"`
}
