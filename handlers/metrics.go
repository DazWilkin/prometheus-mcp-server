package handlers

import (
	"github.com/DazWilkin/prometheus-mcp-server/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Counter of successful MCP tool invocations
	totalx = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "total",
			Namespace: config.Namespace,
			Subsystem: config.Subsystem,
			Help:      "Total number of successful MCP tool invocations",
		}, []string{
			"tool",
		},
	)
	// Counter of unsuccessful (error-generating) MCP tool invocations
	errorx = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "error",
			Namespace: config.Namespace,
			Subsystem: config.Subsystem,
			Help:      "Total number of unsuccessful MCP tool invocations",
		}, []string{
			"tool",
		},
	)
)
