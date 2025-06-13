// Not currently used
package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/DazWilkin/prometheus-mcp-server/testdata"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/prometheus/client_golang/api"
)

// TestAlerts tests Alerts
// https://prometheus.io/docs/prometheus/latest/querying/api/#alerts
func TestAlerts(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	want := string(testdata.JsonAlertsResult)

	mux.HandleFunc("/api/v1/alerts", func(w http.ResponseWriter, r *http.Request) {
		logger := logger.With("handler", "alerts")
		logger.Info("Entered")
		defer logger.Info("Exited")

		s := fmt.Sprintf(`{"data":%s,"status":"success"}`, want)

		w.Header().Set("Content-Type", "application/json")
		// if err := json.NewEncoder(w).Encode(s); err != nil {
		if _, err := w.Write([]byte(s)); err != nil {
			msg := "error encoding JSON"
			logger.Error(msg, "err", err)
			http.Error(w, msg, http.StatusInternalServerError)
		}
	})

	mockPrometheus := server.URL

	apiClient, err := api.NewClient(api.Config{
		Address: mockPrometheus,
	})
	if err != nil {
		t.Errorf("unable to create Prometheus API client")
	}

	c := NewClient(apiClient, logger)

	rqst := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
		Params: mcp.CallToolParams{
			Name:      "Alerts",
			Arguments: map[string]any{},
		},
	}
	resp, err := c.Alerts(context.TODO(), rqst)
	if err != nil {
		t.Errorf("unable to invoke Alerts method")
	}

	logger.Info("Response", "resp", resp)

	if len(resp.Content) == 0 {
		t.Errorf("expected content")
	}

	content := resp.Content[0].(mcp.TextContent)

	if content.Type != "text" {
		t.Errorf("expected text content")
	}

	got := content.Text
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

// TestMetrics test Metrics
// https://prometheus.io/docs/prometheus/latest/querying/api/#querying-label-values
func TestMetrics(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	want := string(testdata.JsonLabelValues)

	mux.HandleFunc("/api/v1/label/__name__/values", func(w http.ResponseWriter, r *http.Request) {
		logger := logger.With("handler", "metrics")
		logger.Info("Entered")
		defer logger.Info("Exited")

		s := fmt.Sprintf(`{"data":%s,"status":"success"}`, want)

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(s)); err != nil {
			msg := "error encoding JSON"
			logger.Error(msg, "err", err)
			http.Error(w, msg, http.StatusInternalServerError)
		}
	})

	mockPrometheus := server.URL

	apiClient, err := api.NewClient(api.Config{
		Address: mockPrometheus,
	})
	if err != nil {
		t.Errorf("unable to create Prometheus API client")
	}

	c := NewClient(apiClient, logger)

	rqst := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
		Params: mcp.CallToolParams{
			Name:      "Metrics",
			Arguments: map[string]any{},
		},
	}
	resp, err := c.Metrics(context.TODO(), rqst)
	if err != nil {
		t.Errorf("unable to invoke Metrics method")
	}

	logger.Info("Response", "resp", resp)

	if len(resp.Content) == 0 {
		t.Errorf("expected content")
	}

	content := resp.Content[0].(mcp.TextContent)
	if content.Type != "text" {
		t.Errorf("expected text content")
	}

	got := content.Text
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

// TestQuery tests Query
// https://prometheus.io/docs/prometheus/latest/querying/api/#instant-queries
func TestQuery(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	b, err := testdata.ModelVector.Type().MarshalJSON()
	if err != nil {
		msg := "unable to marshal Prometheus type"
		logger.Error(msg, "err", err)
		t.Error(msg)
	}

	want := string(b)
	logger.Info("Output", "want", want)

	mux.HandleFunc("/api/v1/query", func(w http.ResponseWriter, r *http.Request) {
		logger := logger.With("handler", "query")
		logger.Info("Entered")
		defer logger.Info("Exited")

		// JSON-RPC method arguments are the request body
		// Example "query=up%7Bjob%3D%22prometheus%22%7D&time=1749772800"
		b, err := io.ReadAll(r.Body)
		if err != nil {
			msg := "error reading request body"
			logger.Error(msg, "err", err)
			http.Error(w, msg, http.StatusInternalServerError)
		}
		defer func() {
			if err := r.Body.Close(); err != nil {
				msg := "unable to close request body"
				logger.Error(msg, "err", err)
				http.Error(w, msg, http.StatusInternalServerError)
			}
		}()

		// Parse the URL-encoded QueryString
		querystring := string(b)
		var values url.Values
		if values, err = url.ParseQuery(querystring); err != nil {
			msg := "unable to parse query string from body"
			logger.Error(msg, "err", err)
			http.Error(w, msg, http.StatusInternalServerError)
		}

		logger.Info("Values", "values", values)

		// Retrieve parameters from the decoded map
		query := values.Get("query")
		ts := values.Get("time")

		logger.Info("Arguments", "query", query, "time", ts)

		s := fmt.Sprintf(`{"data":%s,"status":"success"}`, want)

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(s)); err != nil {
			msg := "error encoding JSON"
			logger.Error(msg, "err", err)
			http.Error(w, msg, http.StatusInternalServerError)
		}
	})

	mockPrometheus := server.URL

	apiClient, err := api.NewClient(api.Config{
		Address: mockPrometheus,
	})
	if err != nil {
		t.Errorf("unable to create Prometheus API client")
	}

	c := NewClient(apiClient, logger)

	rqst := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
		Params: mcp.CallToolParams{
			Name: "Query",
			Arguments: map[string]any{
				"query":   "up{job=\"prometheus\"}",
				"time":    testdata.Time,
				"timeout": testdata.Duration,
				"limit":   testdata.Limit,
			},
		},
	}
	resp, err := c.Query(context.TODO(), rqst)
	if err != nil {
		t.Errorf("unable to invoke Metrics method")
	}

	logger.Info("Response", "resp", resp)

	if len(resp.Content) == 0 {
		t.Errorf("expected content")
	}

	content := resp.Content[0].(mcp.TextContent)
	if content.Type != "text" {
		t.Errorf("expected text content")
	}

	got := content.Text
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}
