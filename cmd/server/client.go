package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/client_golang/prometheus"
)

// Client is a type that represents a Prometheus client
type Client struct {
	v1api  v1.API
	logger *slog.Logger
}

// NewClient is a function that creates a new Client
func NewClient(apiClient api.Client, logger *slog.Logger) *Client {
	logger.Info("Creating new Prometheus client")
	v1api := v1.NewAPI(apiClient)
	return &Client{
		v1api:  v1api,
		logger: logger,
	}
}

// Err is a function that combines logging, metrics and returning errors
func Err(method, msg string, err error, logger *slog.Logger) (*mcp.CallToolResult, *ErrPrometheusClient) {
	logger.Error(msg, "err", err)
	errorx.With(prometheus.Labels{
		"method": method,
	}).Inc()
	return mcp.NewToolResultError(msg), NewErrPrometheusClient(msg, err)
}

// Alerts ia a method that queries Prometheus for a list of Alerts
func (x *Client) Alerts(ctx context.Context, rqst mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	method := "Alerts"
	logger := x.logger.With("method", method)
	logger.Info("Entered")
	defer logger.Info("Exited")

	// Increment Prometheus total metric
	// Increment Prometheus total metric
	// Increment Prometheus total metric
	totalx.With(prometheus.Labels{
		"method": method,
	}).Inc()

	// Invoke Prometheus Alerts method
	alerts, err := x.v1api.Alerts(ctx)
	if err != nil {
		msg := "unable to retrieve alerts"
		return Err(method, msg, err, logger)
	}

	b, err := json.Marshal(alerts)
	if err != nil {
		msg := "unable to marshal alerts"
		return Err(method, msg, err, logger)
	}

	return mcp.NewToolResultText(string(b)), nil
}

// Metrics is a method that queries Prometheus for a list of Metrics
func (x *Client) Metrics(ctx context.Context, rqst mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	method := "Metrics"
	logger := x.logger.With("method", method)
	logger.Info("Entered")
	defer logger.Info("Exited")

	// Increment Prometheus total metric
	totalx.With(prometheus.Labels{
		"method": method,
	}).Inc()

	// Invoke Prometheus LabelValues method
	values, warnings, err := x.v1api.LabelValues(ctx, "__name__", nil, time.Time{}, time.Time{})
	if err != nil {
		msg := "unable to retrieve metrics"
		return Err(method, msg, err, logger)
	}

	logger.Info("Warnings", "warnings", warnings)

	b, err := json.Marshal(values)
	if err != nil {
		msg := "unable to marshal metrics"
		return Err(method, msg, err, logger)
	}

	return mcp.NewToolResultText(string(b)), nil
}

// Query is a method that queries Prometheus with PromQL and returns an instant query
func (x *Client) Query(ctx context.Context, rqst mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	method := "Query"
	logger := x.logger.With("method", method)
	logger.Info("Entered")
	defer logger.Info("Exited")

	// Increment Prometheus total metric
	totalx.With(prometheus.Labels{
		"method": method,
	}).Inc()

	// Tool provides arguments; retrieve these
	// required: query
	// optional: time, timeout, limit
	args := rqst.GetArguments()

	// Required
	query := args["query"].(string)

	// Optional
	// Required by Prometheus API method
	ts, err := extractTimestamp(args["time"], logger)
	if err != nil {
		msg := "unable to extract 'time' parameter"
		return Err(method, msg, err, logger)
	}

	// Optional
	// Optional for Prometheus API method: timeout,limit
	opts, err := extractOptions(args, logger)
	if err != nil {
		msg := "unable to extract optional arguments"
		return Err(method, msg, err, logger)
	}

	// Invoke Prometheus Query method
	value, warnings, err := x.v1api.Query(ctx, query, ts, opts...)
	if err != nil {
		msg := "unable to retrieve query results"
		return Err(method, msg, err, logger)
	}

	logger.Info("Warnings", "warnings", warnings)

	b, err := json.Marshal(value)
	if err != nil {
		msg := "unable to marshal query results"
		return Err(method, msg, err, logger)
	}

	return mcp.NewToolResultText(string(b)), nil
}

// QueryRange is a method queries Promethues with PromQL and returns a range query
func (x *Client) QueryRange(ctx context.Context, rqst mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	method := "QueryRange"
	logger := x.logger.With("method", method)
	logger.Info("Entered")
	defer logger.Info("Exited")

	// Increment Prometheus total metric
	totalx.With(prometheus.Labels{
		"method": method,
	}).Inc()

	args := rqst.GetArguments()

	// Required
	query := args["query"].(string)

	start, err := extractTimestamp(args["start"], logger)
	if err != nil {
		msg := "unable to extract 'start' parameter"
		return Err(method, msg, err, logger)
	}

	end, err := extractTimestamp(args["end"], logger)
	if err != nil {
		msg := "unable to extract 'end' parameter"
		return Err(method, msg, err, logger)
	}

	step, err := extractDuration(args["step"], logger)
	if err != nil {
		msg := "unable to extract 'step' parameter"
		return Err(method, msg, err, logger)
	}

	// Create Range
	r := v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	}

	// Optional
	// Optional for Prometheus API method: timeout,limit
	opts, err := extractOptions(args, logger)
	if err != nil {
		msg := "unable to extract optional arguments"
		return Err(method, msg, err, logger)
	}

	// Invoke Prometheus QueryRange method
	value, warnings, err := x.v1api.QueryRange(ctx, query, r, opts...)
	if err != nil {
		msg := "unable to query results"
		return Err(method, msg, err, logger)
	}

	logger.Info("Warnings", "warnings", warnings)

	b, err := json.Marshal(value)
	if err != nil {
		msg := "unable to marshal query results"
		return Err(method, msg, err, logger)
	}

	return mcp.NewToolResultText(string(b)), nil
}

// Rules is a method that queries Prometheus for a list of Rules
func (x *Client) Rules(ctx context.Context, rqst mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	method := "Rules"
	logger := x.logger.With("method", method)
	logger.Info("Entered")
	defer logger.Info("Exited")

	// Increment Prometheus total metric
	totalx.With(prometheus.Labels{
		"method": method,
	}).Inc()

	// Invoke Prometheus Rules method
	rules, err := x.v1api.Rules(ctx)
	if err != nil {
		msg := "unable to retrieve rules"
		return Err(method, msg, err, logger)
	}

	b, err := json.Marshal(rules)
	if err != nil {
		msg := "unable to marshal targets"
		return Err(method, msg, err, logger)
	}

	return mcp.NewToolResultText(string(b)), nil
}

// Targets is a method that queries Prometheus for a list of Targets
func (x *Client) Targets(ctx context.Context, rqst mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	method := "Targets"
	logger := x.logger.With("method", method)
	logger.Info("Entered")
	defer logger.Info("Exited")

	// Increment Prometheus total metric
	totalx.With(prometheus.Labels{
		"method": method,
	}).Inc()

	// Invoke Prometheus Targets method
	targets, err := x.v1api.Targets(ctx)
	if err != nil {
		return Err(method, msg, err, logger)
	}

	b, err := json.Marshal(targets)
	if err != nil {
		msg := "unable to marshal targets"
		return Err(method, msg, err, logger)
	}

	return mcp.NewToolResultText(string(b)), nil
}
