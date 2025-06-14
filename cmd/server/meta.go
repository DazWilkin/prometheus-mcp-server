package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/prometheus/client_golang/prometheus"
)

type Meta struct {
	prometheus string
	logger     *slog.Logger
}

func NewMeta(prometheus string, logger *slog.Logger) *Meta {
	return &Meta{
		prometheus: prometheus,
		logger:     logger,
	}
}

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

// Ping is a method that pings the Prometheus service
func (x *Meta) Ping(ctx context.Context, rqst mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	method := "Ping"
	logger := x.logger.With("method", method)
	logger.Info("Entered")
	defer logger.Info("Exited")

	// Increment Prometheus total metric
	totalx.With(prometheus.Labels{
		"method": method,
	}).Inc()

	// Need an HTTP client
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Health endpoint
	url := fmt.Sprintf("%s/-/healthy", x.prometheus)

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			msg := "unable to close response body"
			logger.Error(msg, "err", err)
		}
	}()

	// Expect 200
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	return mcp.NewToolResultText("OK"), nil
}
