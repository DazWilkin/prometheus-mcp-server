// Not currently used
package handlers

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

func TestAlertmanagers(t *testing.T) {
	t.Skip("Test not implemented but covered by tools tests")
}

// TestAlerts tests Alerts
// https://prometheus.io/docs/prometheus/latest/querying/api/#alerts
func TestAlerts(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	want := string(testdata.JsonAlertsResult)

	mux.HandleFunc("/api/v1/alerts", func(w http.ResponseWriter, r *http.Request) {

		data := want
		resp := fmt.Sprintf(`{"data":%s,"status":"success"}`, data)

		w.Header().Set("Content-Type", "application/json")
		// if err := json.NewEncoder(w).Encode(s); err != nil {
		if _, err := w.Write([]byte(resp)); err != nil {
			msg := "error encoding JSON"
			t.Logf("%s: %+q", msg, err)
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
	resp, err := c.Alerts(context.Background(), rqst)
	if err != nil {
		t.Errorf("unable to invoke Alerts method")
	}

	t.Logf("Response: %+v", resp)

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

// TestExemplars tests Exemplars
func TestExemplars(t *testing.T) {
	t.Skip("Test not implemented but covered by tools tests")
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

		data := want
		resp := fmt.Sprintf(`{"data":%s,"status":"success"}`, data)

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(resp)); err != nil {
			msg := "error encoding JSON"
			t.Logf("%s: %+q", msg, err)
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
	resp, err := c.Metrics(context.Background(), rqst)
	if err != nil {
		t.Errorf("unable to invoke Metrics method")
	}

	t.Logf("Response: %+v", resp)

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

	// JsonModelVector is the correct type to match "got"
	// However, it's not value of "data" returned by the handler
	// See the handler's response construction for more
	want := string(testdata.JsonModelVector)
	t.Logf("Want: %s", want)

	mux.HandleFunc("/api/v1/query", func(w http.ResponseWriter, r *http.Request) {
		// JSON-RPC method arguments are the request body
		// Example "query=up%7Bjob%3D%22prometheus%22%7D&time=1749772800"
		b, err := io.ReadAll(r.Body)
		if err != nil {
			msg := "error reading request body"
			t.Logf("%s: %+q", msg, err)
			http.Error(w, msg, http.StatusInternalServerError)
		}
		defer func() {
			if err := r.Body.Close(); err != nil {
				msg := "unable to close request body"
				t.Logf("%s: %+q", msg, err)
				http.Error(w, msg, http.StatusInternalServerError)
			}
		}()

		// Parse the URL-encoded QueryString
		querystring := string(b)
		var values url.Values
		if values, err = url.ParseQuery(querystring); err != nil {
			msg := "unable to parse query string from body"
			t.Logf("%s: %+q", msg, err)
			http.Error(w, msg, http.StatusInternalServerError)
		}

		t.Logf("Values: %+v", values)

		// Retrieve parameters from the decoded map
		query := values.Get("query")
		ts := values.Get("time")

		t.Logf("Arguments: {query: %s, time: %s}", query, ts)

		// Construction of this handler's response differs to the other tests
		// In the other tests "data" is the value of "want"
		// But, in this case, "model.Value" differs from "model.Vector"
		// https://pkg.go.dev/github.com/prometheus/common/model#Value
		// https://pkg.go.dev/github.com/prometheus/common/model#Vector
		// I know "data" is correct by querying the Prometheus API directly
		// http://localhost:9090/api/v1/query?query=up
		// testdata.JsonModelValue uses the type (!) testdata.ModelValue
		// This type exist solely to implement Prometheus' model.Value interface
		// To be JSON marshaled into the correct value by this handler
		data := testdata.JsonModelValue
		resp := fmt.Sprintf(`{"data":%s,"status":"success"}`, data)
		t.Logf("Response: %s", string(resp))

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(resp)); err != nil {
			msg := "error encoding JSON"
			t.Logf("%s: %+q", msg, err)
			http.Error(w, msg, http.StatusInternalServerError)
		}
	})

	address := server.URL

	apiClient, err := api.NewClient(api.Config{
		Address: address,
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
	resp, err := c.Query(context.Background(), rqst)
	if err != nil {
		t.Errorf("unable to invoke Query method")
	}

	t.Logf("Response: %+v", resp)

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

func TestQueryRange(t *testing.T) {
	t.Skip("Test not implemented but covered by tools tests")
}

func TestRules(t *testing.T) {
	t.Skip("Test not implemented but covered by tools tests")
}

func TestSeries(t *testing.T) {
	t.Skip("Test not implemented but covered by tools tests")
}

func TestStatusTSDB(t *testing.T) {
	t.Skip("Test not implemented but covered by tools tests")
}

func TestTargets(t *testing.T) {
	t.Skip("Test not implemented but covered by tools tests")
}
