package management

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// Client is a type that represents a Prometheus Management API
type Client struct {
	client     *http.Client
	prometheus string
	logger     *slog.Logger
}

// NewClient is a function that creates a new ManagementAPI
func NewClient(prometheus string, logger *slog.Logger) *Client {
	// Need an HTTP client
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	return &Client{
		client:     client,
		prometheus: prometheus,
		logger:     logger,
	}
}

// Do is a function that invokes Prometheus Management API methods
func (x *Client) Do(method string) int {
	logger := x.logger.With("method", method)
	logger.Info("Entered")
	defer logger.Info("Exited")

	url := fmt.Sprintf("%s/-/%s", x.prometheus, method)

	resp, err := x.client.Get(url)
	if err != nil {
		return http.StatusInternalServerError
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			msg := "unable to close response body"
			logger.Error(msg, "err", err)
		}
	}()

	return resp.StatusCode
}

// Healthy is a method that represents the Prometheus Management API Health check
func (x *Client) Healthy() int {
	method := "healthy"
	return x.Do(method)
}

// Ready is a method that represents the Prometheus Management API Readiness check
func (x *Client) Ready() int {
	method := "ready"
	return x.Do(method)
}
