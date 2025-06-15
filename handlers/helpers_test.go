package handlers

import (
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestExtractOptions(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	timeout := time.Duration(15 * time.Second)
	limit := uint64(100)

	args := map[string]any{
		"timeout": timeout.String(),
		"limit":   limit,
	}

	_, err := extractOptions(args, logger)
	if err != nil {
		t.Fatal("expected success")
	}

	// Can't compare []v1.Option easily
}
func TestExtractDuration(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	want := time.Duration(5 * time.Second)

	got, err := extractDuration(want.String(), logger)
	if err != nil {
		t.Fatal("expected success")
	}

	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}

// TestExtractTimestamp tests extractTimestamp
// Ideally want to compare time.Time with time.Time
// But this won't work
// Easiest to use the desired type (RFC-3339)
func TestExtractTimestamp(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	want := time.Now().Format(time.RFC3339)

	ts, err := extractTimestamp(want, logger)
	if err != nil {
		t.Fatal("expected success")
	}

	got := ts.Format(time.RFC3339)

	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}
