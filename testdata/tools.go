package testdata

import "time"

type Tool = string
type ID = string

type ToolsTests = map[Tool]Tests
type Tests = map[ID]Params
type Params = map[string]any

var (
	// Prometheus query to use with tool tests
	query string = `up{job="prometheus"}`
)
var (
	// Dates etc. to use with Prometheus queries in tool tests
	now              = time.Now()
	year, month, day = now.Date()
	location         = now.Location()

	t = now.Format(time.RFC3339)

	// For testing, create start at 00:00:00 and end at 23:59:59 of today
	start = time.Date(year, month, day, 0, 0, 0, 0, location).Format(time.RFC3339)
	end   = time.Date(year, month, day, 23, 59, 59, 0, location).Format(time.RFC3339)

	// Because there is a 24-hour window, step at 1 hour to yield <=24 results
	step = "1h"

	// Corresponds to the request timeout
	timeout = "15s"

	// Limit Prometheus query results
	limit = 10
)
var (
	// ExampleToolsTests = ToolsTests{
	// 	"rule01": {
	// 		"test01": {
	// 			"param01": "value01",
	// 		},
	// 	},
	// }
	MetaToolsTests = ToolsTests{
		"ping": {
			// No additional params
			"": {},
		},
	}
	ClientToolsTests = ToolsTests{
		"alertmanagers": {
			// No additional params
			"": {},
		},
		"alerts": {
			// No additional params
			"": {},
		},
		"exemplars": {
			"required": {
				"query": "up{}",
				"start": start,
				"end":   end,
			},
		},
		"metrics": {
			// No additional params
			"": {},
		},
		"query": {
			"required": {
				"query": query,
			},
			"+time": {
				"query": query,
				"time":  t,
			},
			"+time+timeout": {
				"query":   query,
				"time":    t,
				"timeout": timeout,
			},
			"+time+timeout+limit": {
				"query":   query,
				"time":    t,
				"timeout": timeout,
				"limit":   limit,
			},
		},
		"query_range": {
			"required": {
				"query": query,
				"start": start,
				"end":   end,
				"step":  step,
			},
			"+limit": {
				"query": query,
				"start": start,
				"end":   end,
				"step":  step,
				"limit": limit,
			},
		},
		"rules": {
			// No additional params
			"": {},
		},
		"series": {
			"required": {
				// match[] uniquely is a repeated field presented as slice of strings
				"match[]": []string{
					query,
				},
				"start": start,
				"end":   end,
			},
			"+limit": {
				"match[]": []string{
					query,
				},
				"start": start,
				"end":   end,
				"limit": limit,
			},
		},
		"status_tsdb": {
			// No additional params
			"": {},
		},
		"targets": {
			// No additional params
			"": {},
		},
	}
)
