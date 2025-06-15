package main

import (
	"flag"
	"fmt"
	"log/slog"
)

// Config is a type that represent the app's configuration
type Config struct {
	Prometheus string
	Server     Server
	Metric     Metric
}

// NewConfig is a function that creates a new Config
func NewConfig(logger *slog.Logger) (*Config, error) {
	// MCP config
	// If server.addr=="", MCP will be configured to use stdio not HTTP
	serverAddr := flag.String("server.addr", ":7777", "Endpoint on which MCP tools are published")
	serverPath := flag.String("server.path", "/mcp", "Path on which MCP tools are served")

	// Metrics config
	// If metric.addr=="", Prometheus metrics will **not** be exported
	metricAddr := flag.String("metric.addr", ":8080", "Endpoint on which metrics are published")
	metricPath := flag.String("metric.path", "/metrics", "Path on which metrics are served")

	// Prometheus server
	prometheus := flag.String("prometheus", "http://localhost:9090", "Endpoint of Prometheus server")

	flag.Parse()

	if *prometheus == "" {
		msg := "Flag '--prometheus' is required"
		err := NewErrConfig(msg, nil)
		logger.Error(msg, "err", err)
		return nil, err
	}

	return &Config{
		Prometheus: *prometheus,
		Server: Server{
			Addr: *serverAddr,
			Path: *serverPath,
		},
		Metric: Metric{
			Addr: *metricAddr,
			Path: *metricPath,
		},
	}, nil
}

// Server represents the MCP server's configuration
// TODO(dazwilkin): Possibly unify with Metric type?
type Server struct {
	Addr string
	Path string
}

// GoString is a method that generates a Go string
func (m Server) GoString() string {
	return fmt.Sprintf("MCP{Addr: %q, Path: %q}", m.Addr, m.Path)
}

// String is a method that generates a string
func (m Server) String() string {
	return fmt.Sprintf("%s/%s", m.Addr, m.Path)
}

// Metric is a type that represents the Prometheus metrics exporter configuration
// TODO(dazwilkin): Possibly unify with MCP type?
type Metric struct {
	Addr string
	Path string
}

// GoString is a method that returns a Go string
func (m Metric) GoString() string {
	return fmt.Sprintf("Metrics{Addr: %q, Path: %q}", m.Addr, m.Path)
}

// String is a method that returns a string
func (m Metric) String() string {
	return fmt.Sprintf("%s/%s", m.Addr, m.Path)
}
