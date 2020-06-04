/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package metriclint

import (
	"fmt"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	metricTypeCounter = "Counter"
	metricTypeCounterVec = "CounterVec"
	metricTypeGauge = "Gauge"
	metricTypeGaugeVec = "GaugeVec"
	metricTypeHistogram = "Histogram"
	metricTypeHistogramVec = "HistogramVec"
	metricTypeSummary = "Summary"
	metricTypeSummaryVec = "SummaryVec"
)

func TestNoHelpText(t *testing.T) {
	tests := []struct{
		name string
		metricType string
		opts interface{}
		labelNames []string
		expected string
	}{
		{
			name: "counter metric no help",
			metricType: metricTypeCounter,
			opts: prometheus.CounterOpts{
				Name: "metriclint_test_total",
			},
			expected: fmt.Sprintf("metriclint_test_total:%s", LintErrMsgNoHelp),
		},
		{
			name: "counter vector metric no help",
			metricType: metricTypeCounterVec,
			opts: prometheus.CounterOpts{
				Name: "metriclint_test_total",
			},
			labelNames: []string{"label1", "label2", "label3"},
			expected: fmt.Sprintf("metriclint_test_total:%s", LintErrMsgNoHelp),
		},
	}

	for _, test := range tests {
		tc := test
		t.Run(tc.name, func(t *testing.T) {
			var lintResult *LintResult

			switch tc.metricType{
			case metricTypeCounter:
				lintResult = LintCounter(tc.opts.(prometheus.CounterOpts))
			case metricTypeCounterVec:
				lintResult = LintCounterVector(tc.opts.(prometheus.CounterOpts), tc.labelNames)
			}

			if lintResult.String() != tc.expected {
				t.Errorf("expected %s, but got %s", tc.expected, lintResult.String())
			}
		})
	}
}