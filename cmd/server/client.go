package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
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

// Alerts ia a method that queries Prometheus for a list of Alerts
func (x *Client) Alerts(ctx context.Context, rqst mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	method := "Alerts"
	logger := x.logger.With("method", method)
	logger.Info("Entered")
	defer logger.Info("Exited")

	alerts, err := x.v1api.Alerts(ctx)
	if err != nil {
		msg := "unable to retrieve alerts"
		logger.Error(msg, "err", err)
		return mcp.NewToolResultError(msg), NewErrPrometheusClient(msg, err)
	}

	b, err := json.Marshal(alerts)
	if err != nil {
		msg := "unable to marshal alerts"
		logger.Error(msg, "err", err)
		return mcp.NewToolResultError(msg), NewErrPrometheusClient(msg, err)
	}

	return mcp.NewToolResultText(string(b)), nil
}

// Metrics is a method that queries Prometheus for a list of Metrics
func (x *Client) Metrics(ctx context.Context, rqst mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	method := "Metrics"
	logger := x.logger.With("method", method)
	logger.Info("Entered")
	defer logger.Info("Exited")

	values, warnings, err := x.v1api.LabelValues(ctx, "__name__", nil, time.Time{}, time.Time{})
	if err != nil {
		msg := "unable to retrieve metrics"
		logger.Error(msg, "err", err)
		return mcp.NewToolResultError(msg), NewErrPrometheusClient(msg, err)
	}

	logger.Info("Warnings", "warnings", warnings)

	b, err := json.Marshal(values)
	if err != nil {
		msg := "unable to marshal metrics"
		logger.Error(msg, "err", err)
		return mcp.NewToolResultError(msg), NewErrPrometheusClient(msg, err)
	}

	return mcp.NewToolResultText(string(b)), nil
}

// Query is a method that queries Prometheus with PromQL and returns an instant query
func (x *Client) Query(ctx context.Context, rqst mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	method := "Query"
	logger := x.logger.With("method", method)
	logger.Info("Entered")
	defer logger.Info("Exited")

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
		return mcp.NewToolResultError(msg), err
	}

	// Optional
	// Optional for Prometheus API method: timeout,limit
	opts, err := extractOptions(args, logger)
	if err != nil {
		msg := "unable to extract optional arguments"
		return mcp.NewToolResultError(msg), err
	}

	value, warnings, err := x.v1api.Query(ctx, query, ts, opts...)
	if err != nil {
		msg := "unable to retrieve query results"
		logger.Error(msg, "err", err)
		return mcp.NewToolResultError(msg), NewErrPrometheusClient(msg, err)
	}

	logger.Info("Warnings", "warnings", warnings)

	b, err := json.Marshal(value)
	if err != nil {
		msg := "unable to marshal query results"
		logger.Error(msg, "err", err)
		return mcp.NewToolResultError(msg), NewErrPrometheusClient(msg, err)
	}

	return mcp.NewToolResultText(string(b)), nil
}

// QueryRange is a method queries Promethues with PromQL and returns a range query
func (x *Client) QueryRange(ctx context.Context, rqst mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	method := "QueryRange"
	logger := x.logger.With("method", method)
	logger.Info("Entered")
	defer logger.Info("Exited")

	args := rqst.GetArguments()

	// Required
	query := args["query"].(string)

	start, err := extractTimestamp(args["start"], logger)
	if err != nil {
		msg := "unable to extract 'start' parameter"
		return mcp.NewToolResultError(msg), err
	}

	end, err := extractTimestamp(args["end"], logger)
	if err != nil {
		msg := "unable to extract 'end' parameter"
		return mcp.NewToolResultError(msg), err
	}

	step, err := extractDuration(args["step"], logger)
	if err != nil {
		msg := "unable to extract 'step' parameter"
		return mcp.NewToolResultError(msg), err
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
		return mcp.NewToolResultError(msg), err
	}

	value, warnings, err := x.v1api.QueryRange(ctx, query, r, opts...)
	if err != nil {
		msg := "unable to query results"
		logger.Error(msg, "err", err)
		return mcp.NewToolResultError(msg), NewErrPrometheusClient(msg, err)
	}

	logger.Info("Warnings", "warnings", warnings)

	b, err := json.Marshal(value)
	if err != nil {
		msg := "unable to marshal query results"
		logger.Error(msg, "err", err)
		return mcp.NewToolResultError(msg), NewErrPrometheusClient(msg, err)
	}

	return mcp.NewToolResultText(string(b)), nil
}

// Rules is a method that queries Prometheus for a list of Rules
func (x *Client) Rules(ctx context.Context, rqst mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	method := "Rules"
	logger := x.logger.With("method", method)
	logger.Info("Entered")
	defer logger.Info("Exited")

	rules, err := x.v1api.Rules(ctx)
	if err != nil {
		msg := "unable to retrieve rules"
		logger.Error(msg, "err", err)
		return mcp.NewToolResultError(msg), NewErrPrometheusClient(msg, err)
	}

	b, err := json.Marshal(rules)
	if err != nil {
		msg := "unable to marshal targets"
		logger.Error(msg, "err", err)
		return mcp.NewToolResultError(msg), NewErrPrometheusClient(msg, err)
	}

	return mcp.NewToolResultText(string(b)), nil
}

// Targets is a method that queries Prometheus for a list of Targets
func (x *Client) Targets(ctx context.Context, rqst mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	method := "Targets"
	logger := x.logger.With("method", method)
	logger.Info("Entered")
	defer logger.Info("Exited")

	targets, err := x.v1api.Targets(ctx)
	if err != nil {
		msg := "unable to retrieve targets"
		logger.Error(msg, "err", err)
		return mcp.NewToolResultError(msg), NewErrPrometheusClient(msg, err)
	}

	b, err := json.Marshal(targets)
	if err != nil {
		msg := "unable to marshal targets"
		logger.Error(msg, "err", err)
		return mcp.NewToolResultError(msg), NewErrPrometheusClient(msg, err)
	}

	return mcp.NewToolResultText(string(b)), nil
}
