package testdata

type Tool = string
type ID = string

type ToolsTests = map[Tool]Tests
type Tests = map[ID]Params
type Params = map[string]any

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
		"metrics": {
			// No additional params
			"": {},
		},
		"query": {
			"basic": {
				"query": "up{}",
			},
			"+time": {
				"query": "up{}",
				"time":  "2025-06-14T13:00:00-07:00",
			},
			"+time+timeout": {
				"query":   "up{}",
				"time":    "2025-06-14T13:00:00-07:00",
				"timeout": "15s",
			},
			"+time+timeout+limit": {
				"query":   "up{}",
				"time":    "2025-06-14T13:00:00-07:00",
				"timeout": "15s",
				"limit":   10,
			},
		},
		"query_range": {
			"basic": {
				"query": "up{}",
				"start": "2025-06-14T00:00:00-07:00",
				"end":   "2025-06-14T23:59:59-07:00",
				"step":  "1h",
			},
		},
		"rules": {
			// No additional params
			"": {},
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
