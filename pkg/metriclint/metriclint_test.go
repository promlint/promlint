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

func TestLintCounter(t *testing.T) {
	tests := []struct {
		name string
		opts prometheus.CounterOpts
		expectedResult string
	}{
		{
			name: "valid counter",
			opts: prometheus.CounterOpts{
				Name: "lint_test_total",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_total:"),
		},
		{
			name: "all metric should contains help",
			opts: prometheus.CounterOpts{
				Name: "lint_test_total",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_total:%s", LintErrMsgNoHelp),
		},
		{
			name: "should use base unit",
			opts: prometheus.CounterOpts{
				Name: "lint_test_hours_total",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_hours_total:%s", fmt.Sprintf(LintErrMsgNonBaseUnit, "seconds", "hours")),
		},
		{
			name: "counter should contains total suffix",
			opts: prometheus.CounterOpts{
				Name: "lint_test_suffix",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_suffix:%s", LintErrMsgCounterShouldHaveTotalSuffix),
		},
		{
			name: "non histogram should not have le label",
			opts: prometheus.CounterOpts{
				Name: "lint_test_total",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"le": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_total:%s", LintErrMsgNonHistogramShouldNotHaveLeLabel),
		},
		{
			name: "non summary should not have quantile label",
			opts: prometheus.CounterOpts{
				Name: "lint_test_total",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"quantile": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_total:%s", LintErrMsgNonSummaryShouldNotHaveQuantileLabel),
		},
		{
			name: "should not have metric type",
			opts: prometheus.CounterOpts{
				Name: "lint_counter_total",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_counter_total:%s", fmt.Sprintf(LintErrMsgNoMetricType, "counter")),
		},
		{
			name: "should not have special chars",
			opts: prometheus.CounterOpts{
				Name: "lint_:_total",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_:_total:%s", LintErrMsgNoReservedChars),
		},
		{
			name: "name label should in snake case",
			opts: prometheus.CounterOpts{
				Name: "lint_tesT_total",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lName": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_tesT_total:%s,%s", LintErrMsgNameShouldBeSnakeCase, LintErrMsgLabelShouldBeSnakeCase),
		},
		{
			name: "should not contain abbreviated unit",
			opts: prometheus.CounterOpts{
				Name: "lint_ms_total",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_ms_total:%s", LintErrMsgNameShouldNotHaveAbbr),
		},
	}

	for _, test := range tests {
		tc := test
		t.Run(tc.name, func(t *testing.T) {
			lintResult := LintCounter(tc.opts)
			if tc.expectedResult != lintResult.String() {
				t.Errorf("expected: %s, but got: %s", tc.expectedResult, lintResult.String())
			}
		})
	}
}

func TestLintCounterVector(t *testing.T) {
	tests := []struct {
		name string
		opts prometheus.CounterOpts
		labelNames []string
		expectedResult string
	}{
		{
			name: "valid counter vector",
			opts: prometheus.CounterOpts{
				Name: "lint_test_total",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_total:"),
		},
		{
			name: "all metric should contains help",
			opts: prometheus.CounterOpts{
				Name: "lint_test_total",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_total:%s", LintErrMsgNoHelp),
		},
		{
			name: "should use base unit",
			opts: prometheus.CounterOpts{
				Name: "lint_test_hours_total",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_hours_total:%s", fmt.Sprintf(LintErrMsgNonBaseUnit, "seconds", "hours")),
		},
		{
			name: "counter should contains total suffix",
			opts: prometheus.CounterOpts{
				Name: "lint_test_suffix",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_suffix:%s", LintErrMsgCounterShouldHaveTotalSuffix),
		},
		{
			name: "non histogram should not have le label",
			opts: prometheus.CounterOpts{
				Name: "lint_test_total",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"le", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_total:%s", LintErrMsgNonHistogramShouldNotHaveLeLabel),
		},
		{
			name: "non summary should not have quantile label",
			opts: prometheus.CounterOpts{
				Name: "lint_test_total",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"quantile", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_total:%s", LintErrMsgNonSummaryShouldNotHaveQuantileLabel),
		},
		{
			name: "should not have metric type",
			opts: prometheus.CounterOpts{
				Name: "lint_counter_total",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_counter_total:%s", fmt.Sprintf(LintErrMsgNoMetricType, "counter")),
		},
		{
			name: "should not have special chars",
			opts: prometheus.CounterOpts{
				Name: "lint_:_total",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_:_total:%s", LintErrMsgNoReservedChars),
		},
		{
			name: "name label should in snake case",
			opts: prometheus.CounterOpts{
				Name: "lint_tesT_total",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lName1", "lname2"},
			expectedResult: fmt.Sprintf("lint_tesT_total:%s,%s", LintErrMsgNameShouldBeSnakeCase, LintErrMsgLabelShouldBeSnakeCase),
		},
		{
			name: "should not contain abbreviated unit",
			opts: prometheus.CounterOpts{
				Name: "lint_ms_total",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_ms_total:%s", LintErrMsgNameShouldNotHaveAbbr),
		},
	}

	for _, test := range tests {
		tc := test
		t.Run(tc.name, func(t *testing.T) {
			lintResult := LintCounterVector(tc.opts, tc.labelNames)
			if tc.expectedResult != lintResult.String() {
				t.Errorf("expected: %s, but got: %s", tc.expectedResult, lintResult.String())
			}
		})
	}
}

func TestLintGauge(t *testing.T) {
	tests := []struct {
		name string
		opts prometheus.GaugeOpts
		expectedResult string
	}{
		{
			name: "valid gauge",
			opts: prometheus.GaugeOpts{
				Name: "lint_test_numbers",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_numbers:"),
		},
		{
			name: "all metric should contains help",
			opts: prometheus.GaugeOpts{
				Name: "lint_test_numbers",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_numbers:%s", LintErrMsgNoHelp),
		},
		{
			name: "should use base unit",
			opts: prometheus.GaugeOpts{
				Name: "lint_test_hours_numbers",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_hours_numbers:%s", fmt.Sprintf(LintErrMsgNonBaseUnit, "seconds", "hours")),
		},
		{
			name: "non counter should not have total",
			opts: prometheus.GaugeOpts{
				Name: "lint_test_total",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_total:%s", LintErrMsgNonCounterShouldNotHaveTotalSuffix),
		},
		{
			name: "non histogram should not have bucket suffix",
			opts: prometheus.GaugeOpts{
				Name: "lint_test_bucket",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_bucket:%s", LintErrMsgNonHistogramShouldNotHaveBucketSuffix),
		},
		{
			name: "non histogram summary should not have count suffix",
			opts: prometheus.GaugeOpts{
				Name: "lint_test_count",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_count:%s", LintErrMsgNonHistogramSummaryShouldNotHaveCountSuffix),
		},
		{
			name: "non histogram summary should not have sum suffix",
			opts: prometheus.GaugeOpts{
				Name: "lint_test_sum",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_sum:%s", LintErrMsgMonHistogramSummaryShouldNotHaveSumSuffix),
		},
		{
			name: "non histogram should not have le label",
			opts: prometheus.GaugeOpts{
				Name: "lint_test_numbers",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"le": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_numbers:%s", LintErrMsgNonHistogramShouldNotHaveLeLabel),
		},
		{
			name: "non summary should not have quantile label",
			opts: prometheus.GaugeOpts{
				Name: "lint_test_numbers",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"quantile": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_numbers:%s", LintErrMsgNonSummaryShouldNotHaveQuantileLabel),
		},
		{
			name: "should not have metric type",
			opts: prometheus.GaugeOpts{
				Name: "lint_gauge_numbers",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_gauge_numbers:%s", fmt.Sprintf(LintErrMsgNoMetricType, "gauge")),
		},
		{
			name: "should not have special chars",
			opts: prometheus.GaugeOpts{
				Name: "lint_:_numbers",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_:_numbers:%s", LintErrMsgNoReservedChars),
		},
		{
			name: "name label should in snake case",
			opts: prometheus.GaugeOpts{
				Name: "lint_tesT_numbers",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lName": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_tesT_numbers:%s,%s", LintErrMsgNameShouldBeSnakeCase, LintErrMsgLabelShouldBeSnakeCase),
		},
		{
			name: "should not contain abbreviated unit",
			opts: prometheus.GaugeOpts{
				Name: "lint_ms_numbers",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_ms_numbers:%s", LintErrMsgNameShouldNotHaveAbbr),
		},
	}

	for _, test := range tests {
		tc := test
		t.Run(tc.name, func(t *testing.T) {
			lintResult := LintGauge(tc.opts)
			if tc.expectedResult != lintResult.String() {
				t.Errorf("expected: %s, but got: %s", tc.expectedResult, lintResult.String())
			}
		})
	}
}

func TestLintGaugeVector(t *testing.T) {
	tests := []struct {
		name string
		opts prometheus.GaugeOpts
		labelNames []string
		expectedResult string
	}{
		{
			name: "valid gauge vector",
			opts: prometheus.GaugeOpts{
				Name: "lint_test_numbers",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_numbers:"),
		},
		{
			name: "all metric should contains help",
			opts: prometheus.GaugeOpts{
				Name: "lint_test_numbers",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_numbers:%s", LintErrMsgNoHelp),
		},
		{
			name: "should use base unit",
			opts: prometheus.GaugeOpts{
				Name: "lint_test_hours_numbers",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_hours_numbers:%s", fmt.Sprintf(LintErrMsgNonBaseUnit, "seconds", "hours")),
		},
		{
			name: "non counter should not have total",
			opts: prometheus.GaugeOpts{
				Name: "lint_test_total",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_total:%s", LintErrMsgNonCounterShouldNotHaveTotalSuffix),
		},
		{
			name: "non histogram should not have bucket suffix",
			opts: prometheus.GaugeOpts{
				Name: "lint_test_bucket",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_bucket:%s", LintErrMsgNonHistogramShouldNotHaveBucketSuffix),
		},
		{
			name: "non histogram summary should not have count suffix",
			opts: prometheus.GaugeOpts{
				Name: "lint_test_count",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_count:%s", LintErrMsgNonHistogramSummaryShouldNotHaveCountSuffix),
		},
		{
			name: "non histogram summary should not have sum suffix",
			opts: prometheus.GaugeOpts{
				Name: "lint_test_sum",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_sum:%s", LintErrMsgMonHistogramSummaryShouldNotHaveSumSuffix),
		},
		{
			name: "non histogram should not have le label",
			opts: prometheus.GaugeOpts{
				Name: "lint_test_numbers",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"le", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_numbers:%s", LintErrMsgNonHistogramShouldNotHaveLeLabel),
		},
		{
			name: "non summary should not have quantile label",
			opts: prometheus.GaugeOpts{
				Name: "lint_test_numbers",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"quantile", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_numbers:%s", LintErrMsgNonSummaryShouldNotHaveQuantileLabel),
		},
		{
			name: "should not have metric type",
			opts: prometheus.GaugeOpts{
				Name: "lint_gauge_numbers",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_gauge_numbers:%s", fmt.Sprintf(LintErrMsgNoMetricType, "gauge")),
		},
		{
			name: "should not have special chars",
			opts: prometheus.GaugeOpts{
				Name: "lint_:_numbers",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_:_numbers:%s", LintErrMsgNoReservedChars),
		},
		{
			name: "name label should in snake case",
			opts: prometheus.GaugeOpts{
				Name: "lint_tesT_numbers",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lName1", "lname2"},
			expectedResult: fmt.Sprintf("lint_tesT_numbers:%s,%s", LintErrMsgNameShouldBeSnakeCase, LintErrMsgLabelShouldBeSnakeCase),
		},
		{
			name: "should not contain abbreviated unit",
			opts: prometheus.GaugeOpts{
				Name: "lint_ms_numbers",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_ms_numbers:%s", LintErrMsgNameShouldNotHaveAbbr),
		},
	}

	for _, test := range tests {
		tc := test
		t.Run(tc.name, func(t *testing.T) {
			lintResult := LintGaugeVector(tc.opts, tc.labelNames)
			if tc.expectedResult != lintResult.String() {
				t.Errorf("expected: %s, but got: %s", tc.expectedResult, lintResult.String())
			}
		})
	}
}

func TestLintHistogram(t *testing.T) {
	tests := []struct {
		name string
		opts prometheus.HistogramOpts
		expectedResult string
	}{
		{
			name: "valid histogram",
			opts: prometheus.HistogramOpts{
				Name: "lint_test_seconds",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_seconds:"),
		},
		{
			name: "all metric should contains help",
			opts: prometheus.HistogramOpts{
				Name: "lint_test_seconds",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_seconds:%s", LintErrMsgNoHelp),
		},
		{
			name: "should use base unit",
			opts: prometheus.HistogramOpts{
				Name: "lint_test_hours",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_hours:%s", fmt.Sprintf(LintErrMsgNonBaseUnit, "seconds", "hours")),
		},
		{
			name: "non counter should not have total",
			opts: prometheus.HistogramOpts{
				Name: "lint_test_total",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_total:%s", LintErrMsgNonCounterShouldNotHaveTotalSuffix),
		},
		{
			name: "non summary should not have quantile label",
			opts: prometheus.HistogramOpts{
				Name: "lint_test_seconds",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"quantile": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_seconds:%s", LintErrMsgNonSummaryShouldNotHaveQuantileLabel),
		},
		{
			name: "should not contains type name",
			opts: prometheus.HistogramOpts{
				Name: "lint_histogram_seconds",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_histogram_seconds:%s", fmt.Sprintf(LintErrMsgNoMetricType, "histogram")),
		},
		{
			name: "should not have special chars",
			opts: prometheus.HistogramOpts{
				Name: "lint_:_seconds",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_:_seconds:%s", LintErrMsgNoReservedChars),
		},
		{
			name: "name label should in snake case",
			opts: prometheus.HistogramOpts{
				Name: "lint_tesT_seconds",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lName": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_tesT_seconds:%s,%s", LintErrMsgNameShouldBeSnakeCase, LintErrMsgLabelShouldBeSnakeCase),
		},
		{
			name: "should not contain abbreviated unit",
			opts: prometheus.HistogramOpts{
				Name: "lint_ms_seconds",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_ms_seconds:%s", LintErrMsgNameShouldNotHaveAbbr),
		},
	}

	for _, test := range tests {
		tc := test
		t.Run(tc.name, func(t *testing.T) {
			lintResult := LintHistogram(tc.opts)
			if tc.expectedResult != lintResult.String() {
				t.Errorf("expected: %s, but got: %s", tc.expectedResult, lintResult.String())
			}
		})
	}
}

func TestLintHistogramVector(t *testing.T) {
	tests := []struct {
		name string
		opts prometheus.HistogramOpts
		labelNames []string
		expectedResult string
	}{
		{
			name: "valid histogram",
			opts: prometheus.HistogramOpts{
				Name: "lint_test_seconds",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_seconds:"),
		},
		{
			name: "all metric should contains help",
			opts: prometheus.HistogramOpts{
				Name: "lint_test_seconds",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_seconds:%s", LintErrMsgNoHelp),
		},
		{
			name: "should use base unit",
			opts: prometheus.HistogramOpts{
				Name: "lint_test_hours",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_hours:%s", fmt.Sprintf(LintErrMsgNonBaseUnit, "seconds", "hours")),
		},
		{
			name: "non counter should not have total",
			opts: prometheus.HistogramOpts{
				Name: "lint_test_total",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_total:%s", LintErrMsgNonCounterShouldNotHaveTotalSuffix),
		},
		{
			name: "non summary should not have quantile label",
			opts: prometheus.HistogramOpts{
				Name: "lint_test_seconds",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"quantile", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_seconds:%s", LintErrMsgNonSummaryShouldNotHaveQuantileLabel),
		},
		{
			name: "should not have metric type",
			opts: prometheus.HistogramOpts{
				Name: "lint_histogram_seconds",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_histogram_seconds:%s", fmt.Sprintf(LintErrMsgNoMetricType, "histogram")),
		},
		{
			name: "should not have special chars",
			opts: prometheus.HistogramOpts{
				Name: "lint_:_seconds",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_:_seconds:%s", LintErrMsgNoReservedChars),
		},
		{
			name: "name label should in snake case",
			opts: prometheus.HistogramOpts{
				Name: "lint_tesT_seconds",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lName1", "lname2"},
			expectedResult: fmt.Sprintf("lint_tesT_seconds:%s,%s", LintErrMsgNameShouldBeSnakeCase, LintErrMsgLabelShouldBeSnakeCase),
		},
		{
			name: "should not contain abbreviated unit",
			opts: prometheus.HistogramOpts{
				Name: "lint_ms_seconds",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_ms_seconds:%s", LintErrMsgNameShouldNotHaveAbbr),
		},
	}

	for _, test := range tests {
		tc := test
		t.Run(tc.name, func(t *testing.T) {
			lintResult := LintHistogramVector(tc.opts, tc.labelNames)
			if tc.expectedResult != lintResult.String() {
				t.Errorf("expected: %s, but got: %s", tc.expectedResult, lintResult.String())
			}
		})
	}
}

func TestLintSummary(t *testing.T) {
	tests := []struct {
		name string
		opts prometheus.SummaryOpts
		expectedResult string
	}{
		{
			name: "valid histogram",
			opts: prometheus.SummaryOpts{
				Name: "lint_test_seconds",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_seconds:"),
		},
		{
			name: "all metric should contains help",
			opts: prometheus.SummaryOpts{
				Name: "lint_test_seconds",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_seconds:%s", LintErrMsgNoHelp),
		},
		{
			name: "should use base unit",
			opts: prometheus.SummaryOpts{
				Name: "lint_test_hours",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_hours:%s", fmt.Sprintf(LintErrMsgNonBaseUnit, "seconds", "hours")),
		},
		{
			name: "non counter should not have total",
			opts: prometheus.SummaryOpts{
				Name: "lint_test_total",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_total:%s", LintErrMsgNonCounterShouldNotHaveTotalSuffix),
		},
		{
			name: "non histogram should not have bucket suffix",
			opts: prometheus.SummaryOpts{
				Name: "lint_test_bucket",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_bucket:%s", LintErrMsgNonHistogramShouldNotHaveBucketSuffix),
		},
		{
			name: "non histogram should not have le label",
			opts: prometheus.SummaryOpts{
				Name: "lint_test_seconds",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"le": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_test_seconds:%s", LintErrMsgNonHistogramShouldNotHaveLeLabel),
		},
		{
			name: "should not have metric type",
			opts: prometheus.SummaryOpts{
				Name: "lint_summary_seconds",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_summary_seconds:%s", fmt.Sprintf(LintErrMsgNoMetricType, "summary")),
		},
		{
			name: "should not have special chars",
			opts: prometheus.SummaryOpts{
				Name: "lint_:_seconds",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_:_seconds:%s", LintErrMsgNoReservedChars),
		},
		{
			name: "name label should in snake case",
			opts: prometheus.SummaryOpts{
				Name: "lint_tesT_seconds",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lName": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_tesT_seconds:%s,%s", LintErrMsgNameShouldBeSnakeCase, LintErrMsgLabelShouldBeSnakeCase),
		},
		{
			name: "should not contain abbreviated unit",
			opts: prometheus.SummaryOpts{
				Name: "lint_ms_seconds",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			expectedResult: fmt.Sprintf("lint_ms_seconds:%s", LintErrMsgNameShouldNotHaveAbbr),
		},
	}

	for _, test := range tests {
		tc := test
		t.Run(tc.name, func(t *testing.T) {
			lintResult := LintSummary(tc.opts)
			if tc.expectedResult != lintResult.String() {
				t.Errorf("expected: %s, but got: %s", tc.expectedResult, lintResult.String())
			}
		})
	}
}

func TestLintSummaryVector(t *testing.T) {
	tests := []struct {
		name string
		opts prometheus.SummaryOpts
		labelNames []string
		expectedResult string
	}{
		{
			name: "valid histogram",
			opts: prometheus.SummaryOpts{
				Name: "lint_test_seconds",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_seconds:"),
		},
		{
			name: "all metric should contains help",
			opts: prometheus.SummaryOpts{
				Name: "lint_test_seconds",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_seconds:%s", LintErrMsgNoHelp),
		},
		{
			name: "should use base unit",
			opts: prometheus.SummaryOpts{
				Name: "lint_test_hours",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_hours:%s", fmt.Sprintf(LintErrMsgNonBaseUnit, "seconds", "hours")),
		},
		{
			name: "non counter should not have total",
			opts: prometheus.SummaryOpts{
				Name: "lint_test_total",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_total:%s", LintErrMsgNonCounterShouldNotHaveTotalSuffix),
		},
		{
			name: "non histogram should not have bucket suffix",
			opts: prometheus.SummaryOpts{
				Name: "lint_test_bucket",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_bucket:%s", LintErrMsgNonHistogramShouldNotHaveBucketSuffix),
		},
		{
			name: "non histogram should not have le label",
			opts: prometheus.SummaryOpts{
				Name: "lint_test_seconds",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"le", "lname2"},
			expectedResult: fmt.Sprintf("lint_test_seconds:%s", LintErrMsgNonHistogramShouldNotHaveLeLabel),
		},
		{
			name: "should not have metric type",
			opts: prometheus.SummaryOpts{
				Name: "lint_summary_seconds",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_summary_seconds:%s", fmt.Sprintf(LintErrMsgNoMetricType, "summary")),
		},
		{
			name: "should not have special chars",
			opts: prometheus.SummaryOpts{
				Name: "lint_:_seconds",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_:_seconds:%s", LintErrMsgNoReservedChars),
		},
		{
			name: "name label should in snake case",
			opts: prometheus.SummaryOpts{
				Name: "lint_tesT_seconds",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lName1", "lname2"},
			expectedResult: fmt.Sprintf("lint_tesT_seconds:%s,%s", LintErrMsgNameShouldBeSnakeCase, LintErrMsgLabelShouldBeSnakeCase),
		},
		{
			name: "should not contain abbreviated unit",
			opts: prometheus.SummaryOpts{
				Name: "lint_ms_seconds",
				Help: "this is help message",
				ConstLabels: prometheus.Labels{
					"lname": "lvalue",
				},
			},
			labelNames: []string{"lname1", "lname2"},
			expectedResult: fmt.Sprintf("lint_ms_seconds:%s", LintErrMsgNameShouldNotHaveAbbr),
		},
	}

	for _, test := range tests {
		tc := test
		t.Run(tc.name, func(t *testing.T) {
			lintResult := LintSummaryVector(tc.opts, tc.labelNames)
			if tc.expectedResult != lintResult.String() {
				t.Errorf("expected: %s, but got: %s", tc.expectedResult, lintResult.String())
			}
		})
	}
}
