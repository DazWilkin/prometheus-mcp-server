package handlers

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/DazWilkin/prometheus-mcp-server/config"
	"github.com/DazWilkin/prometheus-mcp-server/testdata"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/prometheus/client_golang/api"
)

// Expect Prometheus server to be running
const (
	p string = "http://localhost:9090"
)

func TestClientTools(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Define MCPServer
	opts := []server.ServerOption{}
	s := server.NewMCPServer(
		"MockPrometheusMCP",
		"0.0.1",
		opts...,
	)

	// Create
	// TODO(dazwilkin): Naming?
	// TODO(dazwilkin): {} suggests refactoring to a function
	{
		config := config.Config{
			Prometheus: p,
		}

		// Create Prometheus API client
		apiClient, err := api.NewClient(api.Config{
			Address: config.Prometheus,
		})
		if err != nil {
			t.Fatalf("unable to create Prometheus API client: %+q", err)
		}

		client := NewClient(apiClient, logger)
		tools := client.Tools()
		s.AddTools(tools...)
	}

	// Test server requires MCPServer
	httptest := server.NewTestStreamableHTTPServer(s)
	defer httptest.Close()

	baseURL := httptest.URL
	t.Logf("Test server URL: %s", baseURL)

	// Create Streamable HTTP MCP client
	client, err := client.NewStreamableHttpClient(baseURL)
	if err != nil {
		t.Errorf("expected success: %+v", err)
	}

	ctx := context.Background()

	// Start the MCP client
	t.Log("Start MCP client")
	if err := client.Start(ctx); err != nil {
		t.Errorf("expected success: %+v", err)
	}

	// Initialize the MCP client
	t.Log("Initialize MCP client")
	{
		resp, err := client.Initialize(ctx, mcp.InitializeRequest{})
		if err != nil {
			t.Errorf("expected success: %+v", err)
		}

		t.Logf("Initialize: %+v", resp)
	}

	// Ping the MCP client
	t.Log("Ping MCP client")
	if err := client.Ping(ctx); err != nil {
		t.Errorf("expected success: %+v", err)
	}

	// ListTools
	t.Log("Use client to list the server's tools")
	{
		rqst := mcp.ListToolsRequest{}
		resp, err := client.ListTools(ctx, rqst)
		if err != nil {
			t.Errorf("expected success: %+v", err)
		}

		t.Logf("Response: %+v", resp)
	}

	// CallTools
	// testdata maps tool names to a map of tests
	for tool, test := range testdata.ClientToolsTests {
		// test maps a test name to a map of tool params
		for name, args := range test {
			t.Logf("[%s] tool: %s; args: %+v",
				name,
				tool,
				args,
			)
			t.Run(name, func(t *testing.T) {
				rqst := mcp.CallToolRequest{
					Params: mcp.CallToolParams{
						Name:      tool,
						Arguments: args,
					},
				}
				resp, err := client.CallTool(ctx, rqst)
				if err != nil {
					t.Errorf("expected success: %+v", err)
				}

				t.Logf("Response: %+v", resp)
			})
		}
	}
}
