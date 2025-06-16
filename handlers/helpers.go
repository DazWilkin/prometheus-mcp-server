package handlers

import (
	"log/slog"
	"time"

	"github.com/DazWilkin/prometheus-mcp-server/errors"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

// extractOptions is a function that extracts optional parameters (timeout|limit) from the arguments
func extractOptions(args map[string]any, logger *slog.Logger) ([]v1.Option, error) {
	opts := []v1.Option{}

	if s, ok := args["timeout"].(string); ok {
		timeout, err := time.ParseDuration(s)
		if err != nil {
			msg := "unable to parse timeout"
			logger.Error(msg, "err", err)
			return opts, errors.NewErrToolHandler(msg, err)
		}
		opts = append(opts, v1.WithTimeout(timeout))
	}

	if s, ok := args["limit"]; ok {
		limit, ok := s.(uint64)
		if ok {
			opts = append(opts, v1.WithLimit(limit))
		} else {
			// The value is optional so, log but continue...
			msg := "unable to parse limit"
			logger.Info(msg)
		}
	}

	return opts, nil
}

// extractDuration is a function that extracts a time.Duration from an argument
func extractDuration(x any, logger *slog.Logger) (time.Duration, error) {
	var d time.Duration

	if s, ok := x.(string); ok {
		var err error
		d, err = time.ParseDuration(s)
		if err != nil {
			msg := "unable to parse duration"
			logger.Error(msg, "err", err)
			return d, errors.NewErrToolHandler(msg, err)
		}
	}

	return d, nil
}

// extractTimestamp is a function that extracts a time.Time from an argument
func extractTimestamp(x any, logger *slog.Logger) (time.Time, error) {
	var t time.Time

	if s, ok := x.(string); ok {
		var err error
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			msg := "unable to parse time"
			logger.Error(msg, "err", err)
			return t, errors.NewErrToolHandler(msg, err)
		}
	}

	return t, nil
}

// extractMatches is a function that extracts recurring "match[]" from an argument
func extractMatches(x any, logger *slog.Logger) ([]string, error) {
	if mm, ok := x.([]any); !ok {
		msg := "unable to extract repeated 'match[]' parameters"
		logger.Info(msg)
		return nil, errors.NewErrToolHandler(msg, nil)
	} else {
		matches := make([]string, len(mm))
		for i, v := range mm {
			if matches[i], ok = v.(string); !ok {
				msg := "unable to convert a 'match[] parameter"
				logger.Info(msg, "match[]", v)
				return nil, errors.NewErrToolHandler(msg, nil)
			}
		}
		return matches, nil
	}
}
