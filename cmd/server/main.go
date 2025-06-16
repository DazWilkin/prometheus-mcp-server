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

	"github.com/DazWilkin/prometheus-mcp-server/config"
	"github.com/DazWilkin/prometheus-mcp-server/handlers"
	"github.com/mark3labs/mcp-go/server"

	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

// buildmetrics is a function that creates a Prometheus metric for build information
func buildmetrics() {
	// Build-related metrics
	buildx := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "build_info",
			Namespace: config.Namespace,
			Subsystem: config.Subsystem,
			Help:      "A metric with a constant '1' value labels by build|start time, git commit, OS and Go versions",
		}, []string{
			"build_time",
			"git_commit",
			"os_version",
			"go_version",
			"start_time",
		},
	)
	// Record the values
	buildx.With(prometheus.Labels{
		"build_time": BuildTime,
		"git_commit": GitCommit,
		"os_version": OSVersion,
		"go_version": GoVersion,
		"start_time": strconv.FormatInt(StartTime, 10),
	}).Inc()

}

// getLogger is a function that creates a logger
// It also logs the build info and records the build info metric
func getLogger(debug bool) *slog.Logger {

	opts := &slog.HandlerOptions{}
	if debug {
		opts.Level = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))

	// Create Prometheus 'static' counter for build config
	logger.Info("Build config",
		"build_time", BuildTime,
		"git_commit", GitCommit,
		"os_version", OSVersion,
		"go_version", GoVersion,
		"start_time", strconv.FormatInt(StartTime, 10),
	)

	return logger
}

// interceptor is a function that intercepts the HTTP request context
// The MCP server is configured to use this interceptor but it exists solely to log when it's called
// It is invoked when GitHub Copilot Agent performs MCP server restart|start|stop operations
// These actions are received as POST requests to the MCP server's endpoint path
// And the Content_Type is set to "application/json"
// And the Content-Length is non-zero (!)
// But the body is empty
func interceptor(logger *slog.Logger) func(ctx context.Context, r *http.Request) context.Context {
	return func(ctx context.Context, r *http.Request) context.Context {
		logger := logger.With("function", "interceptor")
		logger.Debug("Entered")
		defer logger.Debug("Exited")

		// Headers
		logger.Debug("Headers", "headers", r.Header)

		// Does nothing
		return ctx
	}
}

// exporter is a function that creates a Prometheus exporter
func exporter(c *config.Config, logger *slog.Logger) {
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
// The server combines:
// 1. Prometheus HTTP API (Client) tools
// 2. Prometheus Metadata (Meta) tools
func run(c *config.Config, logger *slog.Logger) error {
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

		client := handlers.NewClient(apiClient, logger)
		s.AddTools(client.Tools()...)
	}

	// Create Prometheus Meta proxy
	// TODO(dazwilkin): Naming?
	// TODO(dazwilkin): {} suggests refactoring to a function
	{
		meta := handlers.NewMeta(c.Prometheus, logger)
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
		server.WithHTTPContextFunc(interceptor(logger)),
		server.WithStateLess(true),
	}
	return server.NewStreamableHTTPServer(s, streamOpts...).Start(c.Server.Addr)
}

func main() {
	c, err := config.NewConfig()
	if err != nil {
		msg := "unable to create new config"
		slog.Error(msg, "err", err)
		os.Exit(1)
	}

	logger := getLogger(c.Debug)

	// If configured, start Prometheus metrics exporter in Go routine
	// Check only --metric.addr since --metric.path is optional (default: /metrics)
	if c.Metric.Addr != "" {
		logger.Info("Starting Prometheus metrics exporter",
			"metric.addr", c.Metric.Addr,
			"metric.path", c.Metric.Path,
		)

		// Create and report build metrics
		buildmetrics()

		// Create|Start Prometheus metrics exporter in a Go routine
		go exporter(c, logger)
	}

	// Create|Start MCP server
	logger.Info("Starting Prometheus MCP server")
	up := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "up",
			Namespace: config.Namespace,
			Subsystem: config.Subsystem,
			Help:      "1 if the MCP server is up, 0 otherwise",
		}, nil,
	)
	up.With(nil).Set(1)
	if err := run(c, logger); err != nil {
		msg := "unable to server"
		logger.Error(msg, "err", err)
		up.With(nil).Set(0)
	}
}
