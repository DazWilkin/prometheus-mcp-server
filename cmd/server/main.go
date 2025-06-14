package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/prometheus/client_golang/api"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	prometheus := flag.String("prometheus", "http://localhost:9090", "Endpoint of Prometheus server")
	flag.Parse()

	apiClient, err := api.NewClient(api.Config{
		Address: *prometheus,
	})
	if err != nil {
		slog.Error("unable to create Prometheus API client", "err", err)
		os.Exit(1)
	}

	client := NewClient(apiClient, logger)

	serverOpts := []server.ServerOption{
		server.WithToolCapabilities(false),
	}
	slog.Info("ServerOptions", "opts", serverOpts)

	s := server.NewMCPServer(
		"Prometheus",
		"0.0.1",
		serverOpts...,
	)

	tools := []server.ServerTool{
		{
			Tool: mcp.NewTool(
				"alerts",
				mcp.WithDescription("Prometheus Alerts"),
			),
			Handler: client.Alerts,
		},
		{
			Tool: mcp.NewTool(
				"metrics",
				mcp.WithDescription("Prometheus Metrics"),
			),
			Handler: client.Metrics,
		},
		{
			Tool: mcp.NewTool(
				"query",
				mcp.WithDescription("Prometheus Query"),
				mcp.WithString("query",
					mcp.Required(),
					mcp.Description("Prometheus expression query string"),
				),
				mcp.WithString("time",
					mcp.Description("Evaluation timestamp (RFC-3339 or Unix)"),
				),
				mcp.WithString("timeout",
					mcp.Description("Evaluation timeout"),
				),
				mcp.WithNumber("limit",
					mcp.Description("Maximum number of returned series"),
				),
			),
			Handler: client.Query,
		},
		{
			Tool: mcp.NewTool(
				"query_range",
				mcp.WithDescription("Prometheus Query Range"),
				mcp.WithString("query",
					mcp.Required(),
					mcp.Description("Prometheus expression query string"),
				),
				mcp.WithString("start",
					mcp.Required(),
					mcp.Description("Start timestamp (RFC-3339 or Unix)"),
				),
				mcp.WithString("end",
					mcp.Required(),
					mcp.Description("End timestamp (RFC-3339 or Unix)"),
				),
				mcp.WithString("step",
					mcp.Required(),
					mcp.Description("Query resolution step width in duration format"),
				),
				mcp.WithString("timeout",
					mcp.Description("Evaluation timeout"),
				),
				mcp.WithNumber("limit",
					mcp.Description("Maximum number of returned series"),
				),
			),
			Handler: client.QueryRange,
		},
		{
			Tool: mcp.NewTool("rules",
				mcp.WithDescription("Prometheus Rules"),
			),
			Handler: client.Rules,
		},
		{
			Tool: mcp.NewTool(
				"targets",
				mcp.WithDescription("Prometheus Targets"),
			),
			Handler: client.Targets,
		},
	}

	s.AddTools(tools...)

	stdioOpts := []server.StdioOption{}
	slog.Info("StdioOptions", "opts", stdioOpts)

	if err := server.ServeStdio(s, stdioOpts...); err != nil {
		slog.Error("unable to server", "err", err)
	}
}
