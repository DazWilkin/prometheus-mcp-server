package testdata

import (
	"encoding/json"
	"fmt"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

const (
	Limit = 5
)

var (
	Alerts = []v1.Alert{
		{
			ActiveAt:    time.Now(),
			Annotations: model.LabelSet{},
			Labels:      model.LabelSet{},
			State:       v1.AlertStateFiring,
			Value:       "",
		},
	}
	AlertsResult = v1.AlertsResult{
		Alerts: Alerts,
	}
	LabelValues = model.LabelValues{
		model.LabelValue("foo"),
		model.LabelValue("bar"),
	}
	Duration    = time.Duration(5 * time.Second)
	ModelVector = model.Vector{
		{
			Metric: model.Metric{
				"__name___": "up",
				"app":       "prometheus",
				"instance":  "localhost:9090",
				"job":       "prometheus",
			},
			Timestamp: model.Time(Timestamp.Unix()),
			Value:     model.SampleValue(1),
		},
	}
	Time      = Timestamp.Format(time.RFC3339)
	Timestamp = time.Date(2025, time.June, 13, 0, 0, 0, 0, time.UTC)
)
var (
	JsonAlertsResult []byte = MustMarshal(AlertsResult)
	JsonLabelValues  []byte = MustMarshal(LabelValues)
	JsonModelVector  []byte = MustMarshal(ModelVector)
)

// ModelValue implements Prometheus' model.Value interface
// This implementation is necessary in order to correctly JSON marshal
// TestQuery's handler response by creating a wrapped for ModelVector
type ModelValue struct {
	ResultType model.ValueType `json:"resultType"`
	Result     model.Value     `json:"result"`
}

var (
	JsonModelValue []byte = MustMarshal(ModelValue{
		ResultType: ModelVector.Type(),
		Result:     ModelVector,
	})
)

// MustMarshal is a function that marshals a type to JSON ignoring errors
func MustMarshal(x any) []byte {
	b, err := json.Marshal(x)
	if err != nil {
		panic(fmt.Sprintf("json.Marshal failed: %v", err))
	}

	return b
}
