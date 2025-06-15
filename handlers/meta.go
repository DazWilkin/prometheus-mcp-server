package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/DazWilkin/prometheus-mcp-server/errors"
	"github.com/DazWilkin/prometheus-mcp-server/management"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/prometheus/client_golang/prometheus"
)

// Meta is a type that represents Prometheus Management API
type Meta struct {
	client *management.Client
	logger *slog.Logger
}

// NewMeta is a function that creates a new Meta
func NewMeta(prometheus string, logger *slog.Logger) *Meta {
	client := management.NewClient(prometheus, logger)
	return &Meta{
		client: client,
		logger: logger,
	}
}

// Tools is a method that returns the MCP server tools implemented by Meta
// For every tool defined in this method, there should be a corresponding handler method
func (x *Meta) Tools() []server.ServerTool {
	method := "tools"
	logger := x.logger.With("method", method)
	logger.Info("Entered")
	defer logger.Info("Exited")

	tools := []server.ServerTool{
		{
			Tool: mcp.NewTool(
				"ping",
				mcp.WithDescription("Ping the Prometheus sevrer"),
			),
			Handler: x.Ping,
		},
	}

	return tools
}

// Ping is a method that pings the Prometheus server's Management API's Readiness check
func (x *Meta) Ping(ctx context.Context, rqst mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	method := "Ping"
	logger := x.logger.With("method", method)
	logger.Info("Entered")
	defer logger.Info("Exited")

	// Increment Prometheus total metric
	totalx.With(prometheus.Labels{
		"tool": method,
	}).Inc()

	// Invoke Prometheus Management Ready method
	respCode := x.client.Ready()

	// Expect 200
	if respCode != http.StatusOK {
		msg := "serv"
		return mcp.NewToolResultError(msg), errors.NewErrToolHandler(msg, nil)
	}

	return mcp.NewToolResultText("OK"), nil
}
