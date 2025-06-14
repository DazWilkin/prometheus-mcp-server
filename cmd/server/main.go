package main

import (
	"log/slog"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace string = "mcp"
	subsystem string = "prometheus"
)

var (
	// BuildTime is the time that this binary was built represented as a UNIX epoch
	BuildTime string
	// GitCommit is the git commit value and is expected to be set during build
	GitCommit string
	// GoVersion is the Golang runtime version
	GoVersion = runtime.Version()
	// OSVersion is the OS version (uname --kernel-release) and is expected to be set during build
	OSVersion string
	// StartTime is the start time of the exporter represented as a UNIX epoch
	StartTime = time.Now().Unix()
)
var (
	buildx = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "build_info",
			Namespace: namespace,
			Subsystem: subsystem,
			Help:      "A metric with a constant '1' value labels by build|start time, git commit, OS and Go versions",
		}, []string{
			"build_time",
			"git_commit",
			"os_version",
			"go_version",
			"start_time",
		},
	)
	totalx = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "total",
			Namespace: namespace,
			Subsystem: subsystem,
			Help:      "Total",
		}, []string{
			"method",
		},
	)
	errorx = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "error",
			Namespace: namespace,
			Subsystem: subsystem,
			Help:      "Error",
		}, []string{
			"method",
		},
	)
)

// getLogger is a function that creates a logger
// It also logs the build info and records the build info metric
func getLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Create Prometheus 'static' counter for build config
	logger.Info("Build config",
		"build_time", BuildTime,
		"git_commit", GitCommit,
		"os_version", OSVersion,
		"go_version", GoVersion,
		"start_time", strconv.FormatInt(StartTime, 10),
	)
	buildx.With(prometheus.Labels{
		"build_time": BuildTime,
		"git_commit": GitCommit,
		"os_version": OSVersion,
		"go_version": GoVersion,
		"start_time": strconv.FormatInt(StartTime, 10),
	}).Inc()

	return logger
}

// metrics is a function that creates a Prometheus exporter
func metrics(c *Config, logger *slog.Logger) {
	function := "metrics"
	logger = logger.With("function", function)

	m := c.Metrics

	mux := http.NewServeMux()
	mux.Handle(m.Path, promhttp.Handler())
	s := &http.Server{
		Addr:    m.Addr,
		Handler: mux,
	}
	listen, err := net.Listen("tcp", m.Addr)
	if err != nil {
		msg := "unable to create listener"
		logger.Error(msg,
			"endpoint", m.Addr,
			"err", err,
		)
	}
	logger.Info("Starting Prometheus metrics exporter",
		"url", m.String(),
	)
	logger.Error("unable to serve",
		"err", s.Serve(listen),
	)
}

// run is a function that creates a Prometheus MCP server
func run(c *Config, logger *slog.Logger) error {
	apiClient, err := api.NewClient(api.Config{
		Address: c.Prometheus,
	})
	if err != nil {
		logger.Error("unable to create Prometheus API client", "err", err)
		os.Exit(1)
	}

	client := NewClient(apiClient, logger)

	serverOpts := []server.ServerOption{
		server.WithToolCapabilities(false),
	}
	logger.Info("ServerOptions", "opts", serverOpts)

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
	logger.Info("StdioOptions", "opts", stdioOpts)

	return server.ServeStdio(s, stdioOpts...)
}

func main() {
	logger := getLogger()
	config := NewConfig()

	// Start Prometheus metrics exporter in Go routine
	go metrics(config, logger)

	// Create|Start MCP server
	if err := run(config, logger); err != nil {
		logger.Error("unable to server", "err", err)
	}
}
