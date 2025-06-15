package management

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHealthy(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	mux := http.NewServeMux()
	ts := httptest.NewServer(mux)
	defer ts.Close()

	mux.HandleFunc("/-/healthy", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	url := ts.URL

	client := NewClient(url, logger)
	got := client.Healthy()
	want := http.StatusOK

	if got != want {
		t.Errorf("got: %s (%d); want: %s (%d)",
			http.StatusText(got), got,
			http.StatusText(want), want,
		)
	}
}
func TestReady(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	mux := http.NewServeMux()
	ts := httptest.NewServer(mux)
	defer ts.Close()

	mux.HandleFunc("/-/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	url := ts.URL

	client := NewClient(url, logger)
	got := client.Ready()
	want := http.StatusOK

	if got != want {
		t.Errorf("got: %s (%d); want: %s (%d)",
			http.StatusText(got), got,
			http.StatusText(want), want,
		)
	}
}
