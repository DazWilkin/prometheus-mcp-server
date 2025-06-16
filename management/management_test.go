package management

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	tests = []struct {
		name string
		// A clever way to reference a type's (Client's) methods
		// Since both the handlers a func() int, we can generalize
		handler func(*Client) int
		want    int
	}{
		{
			name:    "healthy",
			handler: (*Client).Healthy,
			want:    http.StatusOK,
		},
		{
			name:    "ready",
			handler: (*Client).Ready,
			want:    http.StatusOK,
		},
	}
)

func okHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("OK")); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func TestManagement(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

			mux := http.NewServeMux()
			ts := httptest.NewServer(mux)
			defer ts.Close()

			pattern := fmt.Sprintf("GET /-/%s", test.name)
			mux.HandleFunc(pattern, okHandler)
			url := ts.URL

			client := NewClient(url, logger)
			// The corresponding way pass the receiver (*Client)
			got := test.handler(client)
			want := test.want

			if got != want {
				t.Errorf("got: %s (%d); want: %s (%d)",
					http.StatusText(got), got,
					http.StatusText(want), want,
				)
			}
		})
	}
}
