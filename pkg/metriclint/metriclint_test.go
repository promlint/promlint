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
