package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

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

	m := c.Metric

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
	function := "run"
	logger = logger.With("function", function)

	serverOpts := []server.ServerOption{
		// server.WithToolCapabilities(true),
		// server.WithResourceCapabilities(true, true),
	}
	logger.Info("ServerOptions", "opts", serverOpts)
	s := server.NewMCPServer(
		"PrometheusMCP",
		"0.0.1",
		serverOpts...,
	)

	// Create Prometheus Client proxy
	// TODO(dazwilkin): Naming?
	// TODO(dazwilkin): {} suggests refactoring to a function
	{
		// Create Prometheus API client
		apiClient, err := api.NewClient(api.Config{
			Address: c.Prometheus,
		})
		if err != nil {
			logger.Error("unable to create Prometheus API client", "err", err)
			os.Exit(1)
		}

		client := NewClient(apiClient, logger)
		s.AddTools(client.Tools()...)
	}

	// Create Prometheus Meta proxy
	// TODO(dazwilkin): Naming?
	// TODO(dazwilkin): {} suggests refactoring to a function
	{
		meta := NewMeta(c.Prometheus, logger)
		s.AddTools(meta.Tools()...)
	}

	stdioOpts := []server.StdioOption{}
	logger.Info("StdioOptions", "opts", stdioOpts)

	// Either
	// Check only --server.addr since --server.path is optional (default: /mcp)
	if c.Server.Addr == "" {
		logger.Info("Configuring Server to use stdio",
			"server.addr", c.Server.Addr,
			"server.path", c.Server.Path,
		)
		return server.ServeStdio(s, stdioOpts...)
	}

	// Or
	logger.Info("Configuring Server to use HTTP streaming",
		"server.addr", c.Server.Addr,
		"server.path", c.Server.Path,
	)
	streamOpts := []server.StreamableHTTPOption{
		server.WithEndpointPath(c.Server.Path), // Default endpoint path
		server.WithHTTPContextFunc(func(ctx context.Context, r *http.Request) context.Context {
			logger := logger.With("function", "WithHttpContextFunc")
			logger.Info("Entered")
			defer logger.Info("Exited")

			// Does nothing
			return ctx
		}),
		server.WithStateLess(true),
	}
	return server.NewStreamableHTTPServer(s, streamOpts...).Start(c.Server.Addr)
}

func main() {
	logger := getLogger()

	config, err := NewConfig(logger)
	if err != nil {
		msg := "unable to create new config"
		logger.Error(msg, "err", err)
		os.Exit(1)
	}

	// If configured, start Prometheus metrics exporter in Go routine
	// Check only --metric.addr since --metric.path is optional (default: /metrics)
	if config.Metric.Addr != "" {
		logger.Info("Starting Prometheus metrics exporter",
			"metric.addr", config.Metric.Addr,
			"metric.path", config.Metric.Path,
		)
		go metrics(config, logger)
	}

	// Create|Start MCP server
	logger.Info("Starting Prometheus MCP server")
	if err := run(config, logger); err != nil {
		msg := "unable to server"
		logger.Error(msg, "err", err)
	}
}
