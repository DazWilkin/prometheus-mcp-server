package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/DazWilkin/prometheus-mcp-server/errors"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
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
func Err(method, msg string, err error, logger *slog.Logger) (*mcp.CallToolResult, *errors.ErrToolHandler) {
	logger.Error(msg, "err", err)

	// Increment Prometheus error metric
	errorx.With(prometheus.Labels{
		"tool": method,
	}).Inc()

	return mcp.NewToolResultError(msg), errors.NewErrToolHandler(msg, err)
}

// Tools is a method that returns the MCP server tools implemeneted by Client
// For every tool defined in this method, there should be a corresponding handler method
func (x *Client) Tools() []server.ServerTool {
	method := "tools"
	logger := x.logger.With("method", method)
	logger.Debug("Entered")
	defer logger.Debug("Exited")

	tools := []server.ServerTool{
		{
			Tool: mcp.NewTool(
				"alertmanagers",
				mcp.WithDescription("Prometheus Alertmanagers"),
			),
			Handler: x.Alertmanagers,
		},
		{
			Tool: mcp.NewTool(
				"alerts",
				mcp.WithDescription("Prometheus Alerts"),
			),
			Handler: x.Alerts,
		},
		{
			Tool: mcp.NewTool(
				"exemplars",
				mcp.WithDescription("Prometheus Exemplars"),
				mcp.WithString("query",
					mcp.Required(),
					mcp.Description("Prometheus expression query string"),
				),
				mcp.WithString("start",
					mcp.Required(),
					mcp.Description("Start timestamp (RFC-3339)"),
				),
				mcp.WithString("end",
					mcp.Required(),
					mcp.Description("End timestamp (RFC-3339)"),
				),
			),
			Handler: x.Exemplars,
		},
		{
			Tool: mcp.NewTool(
				"metrics",
				mcp.WithDescription("Prometheus Metrics"),
			),
			Handler: x.Metrics,
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
					mcp.Description("Evaluation timestamp (RFC-3339)"),
				),
				mcp.WithString("timeout",
					mcp.Description("Evaluation timeout"),
				),
				mcp.WithNumber("limit",
					mcp.Description("Maximum number of returned series"),
				),
			),
			Handler: x.Query,
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
					mcp.Description("Start timestamp (RFC-3339)"),
				),
				mcp.WithString("end",
					mcp.Required(),
					mcp.Description("End timestamp (RFC-3339)"),
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
			Handler: x.QueryRange,
		},
		{
			Tool: mcp.NewTool("rules",
				mcp.WithDescription("Prometheus Rules"),
			),
			Handler: x.Rules,
		},
		{
			Tool: mcp.NewTool(
				"series",
				mcp.WithDescription("Prometheus Series"),
				// https://prometheus.io/docs/prometheus/latest/querying/api/#finding-series-by-label-matchers
				// The query string parameter is "match[]"
				mcp.WithArray("match[]",
					// Need to define the type of the array
					// https://github.com/mark3labs/mcp-go/blob/607df92aa3eac5bbb2f1e1edb2a5fc1da2601c94/examples/typed_tools/main.go#L49
					mcp.Items(map[string]any{"type": "string"}),
					mcp.Required(),
					mcp.Description("Repeated series selector argument that selects the series"),
				),
				mcp.WithString("start",
					mcp.Required(),
					mcp.Description("Start timestamp (RFC-3339)"),
				),
				mcp.WithString("end",
					mcp.Required(),
					mcp.Description("End timestamp (RFC-3339)"),
				),
				mcp.WithNumber("limit",
					mcp.Description("Maximum number of returned series"),
				),
			),
			Handler: x.Series,
		},
		{
			Tool: mcp.NewTool(
				"status_tsdb",
				mcp.WithDescription("Prometheus Status: TSDB"),
			),
			Handler: x.StatusTSDB,
		},
		{
			Tool: mcp.NewTool(
				"targets",
				mcp.WithDescription("Prometheus Targets"),
			),
			Handler: x.Targets,
		},
	}
	return tools
}

// Alertmanagers is a method that queries Prometheus for a list of Alertmanagers
func (x *Client) Alertmanagers(ctx context.Context, rqst mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	method := "Alertmanagers"
	logger := x.logger.With("method", method)
	logger.Debug("Entered")
	defer logger.Debug("Exited")

	// Increment Prometheus total metric
	totalx.With(prometheus.Labels{
		"tool": method,
	}).Inc()

	// Invoke Prometheus Alertmanagers method
	result, err := x.v1api.AlertManagers(ctx)
	if err != nil {
		msg := "unable to retrieve alertmanagers"
		return Err(method, msg, err, logger)
	}

	logger.Info("Alertmanagers retrieved",
		"active", len(result.Active),
		"dropped", len(result.Dropped),
	)

	b, err := json.Marshal(result)
	if err != nil {
		msg := "unable to marshal alertmanagers"
		return Err(method, msg, err, logger)
	}

	return mcp.NewToolResultText(string(b)), nil
}

// Alerts ia a method that queries Prometheus for a list of Alerts
func (x *Client) Alerts(ctx context.Context, rqst mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	method := "Alerts"
	logger := x.logger.With("method", method)
	logger.Debug("Entered")
	defer logger.Debug("Exited")

	// Increment Prometheus total metric
	totalx.With(prometheus.Labels{
		"tool": method,
	}).Inc()

	// Invoke Prometheus Alerts method
	result, err := x.v1api.Alerts(ctx)
	if err != nil {
		msg := "unable to retrieve alerts"
		return Err(method, msg, err, logger)
	}

	logger.Info("Alerts retrieved",
		"alerts", len(result.Alerts),
	)

	b, err := json.Marshal(result)
	if err != nil {
		msg := "unable to marshal alerts"
		return Err(method, msg, err, logger)
	}

	return mcp.NewToolResultText(string(b)), nil
}

// Exemplars is a method that queries Prometheus for a list of Exemplars
func (x *Client) Exemplars(ctx context.Context, rqst mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	method := "Exemplars"
	logger := x.logger.With("method", method)
	logger.Debug("Entered")
	defer logger.Debug("Exited")

	// Increment Prometheus total metric
	totalx.With(prometheus.Labels{
		"tool": method,
	}).Inc()

	// Tool provides arguments; retrieve these
	// required: query, start, end
	args := rqst.GetArguments()
	// Required
	query := args["query"].(string)
	startTime, err := extractTimestamp(args["start"], logger)
	if err != nil {
		msg := "unable to extract 'start' parameter"
		return Err(method, msg, err, logger)
	}
	endTime, err := extractTimestamp(args["end"], logger)
	if err != nil {
		msg := "unable to extract 'end' parameter"
		return Err(method, msg, err, logger)
	}

	// Invoke Prometheus Exemplars method
	results, err := x.v1api.QueryExemplars(ctx, query, startTime, endTime)
	if err != nil {
		msg := "unable to retrieve exemplars"
		return Err(method, msg, err, logger)
	}

	logger.Info("Exemplars retrieved",
		"exemplars", len(results),
	)
	b, err := json.Marshal(results)
	if err != nil {
		msg := "unable to marshal exemplars"
		return Err(method, msg, err, logger)
	}

	return mcp.NewToolResultText(string(b)), nil
}

// Metrics is a method that queries Prometheus for a list of Metrics
func (x *Client) Metrics(ctx context.Context, rqst mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	method := "Metrics"
	logger := x.logger.With("method", method)
	logger.Debug("Entered")
	defer logger.Debug("Exited")

	// Increment Prometheus total metric
	totalx.With(prometheus.Labels{
		"tool": method,
	}).Inc()

	// Tool provides no arguments (neither required nor optional)
	// Define parameters
	label := "__name__"
	var matches []string
	startTime := time.Time{}
	endTime := time.Time{}

	// Invoke Prometheus LabelValues method
	labelvalues, warnings, err := x.v1api.LabelValues(ctx, label, matches, startTime, endTime)
	if err != nil {
		msg := "unable to retrieve metrics"
		return Err(method, msg, err, logger)
	}

	logger.Info("Metrics retrieved",
		"metrics", len(labelvalues),
	)

	// If there are warnings, log them
	if len(warnings) != 0 {
		logger.Info("Warnings", "warnings", warnings)
	}

	b, err := json.Marshal(labelvalues)
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
	logger.Debug("Entered")
	defer logger.Debug("Exited")

	// Increment Prometheus total metric
	totalx.With(prometheus.Labels{
		"tool": method,
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

	// If there are warnings, log them
	if len(warnings) != 0 {
		logger.Info("Warnings", "warnings", warnings)
	}

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
	logger.Debug("Entered")
	defer logger.Debug("Exited")

	// Increment Prometheus total metric
	totalx.With(prometheus.Labels{
		"tool": method,
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

	// If there are warnings, log them
	if len(warnings) != 0 {
		logger.Info("Warnings", "warnings", warnings)
	}

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
	logger.Debug("Entered")
	defer logger.Debug("Exited")

	// Increment Prometheus total metric
	totalx.With(prometheus.Labels{
		"tool": method,
	}).Inc()

	// Invoke Prometheus Rules method
	result, err := x.v1api.Rules(ctx)
	if err != nil {
		msg := "unable to retrieve rules"
		return Err(method, msg, err, logger)
	}

	logger.Info("Rules retrieved",
		"rules", len(result.Groups),
	)

	b, err := json.Marshal(result)
	if err != nil {
		msg := "unable to marshal targets"
		return Err(method, msg, err, logger)
	}

	return mcp.NewToolResultText(string(b)), nil
}

// Series is a method that queries Prometheus for a list of Series
func (x *Client) Series(ctx context.Context, rqst mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	method := "Series"
	logger := x.logger.With("method", method)
	logger.Debug("Entered")
	defer logger.Debug("Exited")

	// Increment Prometheus total metric
	totalx.With(prometheus.Labels{
		"tool": method,
	}).Inc()

	args := rqst.GetArguments()

	// Required
	matches, err := extractMatches(args["match[]"], logger)
	if err != nil {
		msg := "unable to extract repeated 'match[]' parameters"
		return Err(method, msg, err, logger)
	}

	startTime, err := extractTimestamp(args["start"], logger)
	if err != nil {
		msg := "unable to extract 'start' parameter"
		return Err(method, msg, err, logger)
	}

	endTime, err := extractTimestamp(args["end"], logger)
	if err != nil {
		msg := "unable to extract 'end' parameter"
		return Err(method, msg, err, logger)
	}

	// Optional
	// Optional for Prometheus API method: timeout,limit
	opts, err := extractOptions(args, logger)
	if err != nil {
		msg := "unable to extract optional arguments"
		return Err(method, msg, err, logger)
	}

	// Invoke Prometheus Series method
	results, warnings, err := x.v1api.Series(ctx, matches, startTime, endTime, opts...)
	if err != nil {
		msg := "unable to retrieve rules"
		return Err(method, msg, err, logger)
	}

	logger.Info("Series retrieved",
		"series", len(results),
	)

	// If there are warnings, log them
	if len(warnings) != 0 {
		logger.Info("Warnings", "warnings", warnings)
	}

	b, err := json.Marshal(results)
	if err != nil {
		msg := "unable to marshal targets"
		return Err(method, msg, err, logger)
	}
	return mcp.NewToolResultText(string(b)), nil
}

// StatusTSDB is a method that queries Prometheus for the status of its time-series database
func (x *Client) StatusTSDB(ctx context.Context, rqst mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	method := "StatusTSDB"
	logger := x.logger.With("method", method)
	logger.Debug("Entered")
	defer logger.Debug("Exited")

	// Increment Prometheus total metric
	totalx.With(prometheus.Labels{
		"tool": method,
	}).Inc()

	// Invoke Prometheus Status TSDB method
	tsdb, err := x.v1api.TSDB(ctx)
	if err != nil {
		msg := "unable to retrieve TSDB status"
		return Err(method, msg, err, logger)
	}

	b, err := json.Marshal(tsdb)
	if err != nil {
		msg := "unable to marshal TSDB status"
		return Err(method, msg, err, logger)
	}

	return mcp.NewToolResultText(string(b)), nil
}

// Targets is a method that queries Prometheus for a list of Targets
func (x *Client) Targets(ctx context.Context, rqst mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	method := "Targets"
	logger := x.logger.With("method", method)
	logger.Debug("Entered")
	defer logger.Debug("Exited")

	// Increment Prometheus total metric
	totalx.With(prometheus.Labels{
		"tool": method,
	}).Inc()

	// Invoke Prometheus Targets method
	result, err := x.v1api.Targets(ctx)
	if err != nil {
		msg := "unable to retrieve targets"
		return Err(method, msg, err, logger)
	}

	logger.Info("Targets retrieved",
		"active", len(result.Active),
		"dropped", len(result.Dropped),
	)

	b, err := json.Marshal(result)
	if err != nil {
		msg := "unable to marshal targets"
		return Err(method, msg, err, logger)
	}

	return mcp.NewToolResultText(string(b)), nil
}
