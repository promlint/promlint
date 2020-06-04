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
	"regexp"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

// Units and their possible prefixes recognized by this library.  More can be
// added over time as needed.
var (
	// map a unit to the appropriate base unit.
	units = map[string]string{
		// Base units.
		"amperes": "amperes",
		"bytes":   "bytes",
		"celsius": "celsius", // Celsius is more common in practice than Kelvin.
		"grams":   "grams",
		"joules":  "joules",
		"meters":  "meters", // Both American and international spelling permitted.
		"metres":  "metres",
		"seconds": "seconds",
		"volts":   "volts",

		// Non base units.
		// Time.
		"minutes": "seconds",
		"hours":   "seconds",
		"days":    "seconds",
		"weeks":   "seconds",
		// Temperature.
		"kelvin":     "celsius",
		"kelvins":    "celsius",
		"fahrenheit": "celsius",
		"rankine":    "celsius",
		// Length.
		"inches": "meters",
		"yards":  "meters",
		"miles":  "meters",
		// Bytes.
		"bits": "bytes",
		// Energy.
		"calories": "joules",
		// Mass.
		"pounds": "grams",
		"ounces": "grams",
	}

	unitPrefixes = []string{
		"pico",
		"nano",
		"micro",
		"milli",
		"centi",
		"deci",
		"deca",
		"hecto",
		"kilo",
		"kibi",
		"mega",
		"mibi",
		"giga",
		"gibi",
		"tera",
		"tebi",
		"peta",
		"pebi",
	}

	// Common abbreviations that we'd like to discourage.
	unitAbbreviations = []string{
		"s",
		"ms",
		"us",
		"ns",
		"sec",
		"b",
		"kb",
		"mb",
		"gb",
		"tb",
		"pb",
		"m",
		"h",
		"d",
	}
)

var camelCase = regexp.MustCompile(`[a-z][A-Z]`)

const (
	NameSuffixSum = "_sum"

	LabelLe = "le"
	LabelQuantile = "quantile"
)

const (
	LintErrMsgNoHelp = "no help text"
	LintErrMsgNonBaseUnit = `use base unit "%s" instead of "%s"`
)

func lintHelp(help string) (issues []string) {
	if len(help) == 0 {
		issues = append(issues, "no help text")
	}

	return
}

func hasTotalSuffix(name string) bool {
	return strings.HasSuffix(name, "_total")
}

func lintCounterContainsTotal(name string) (issues []string) {
	if !hasTotalSuffix(name) {
		issues = append(issues, `counter metrics should have "_total" suffix`)
	}

	return issues
}

func lintNonCounterNoTotal(name string) (issues []string) {
	if hasTotalSuffix(name) {
		issues = append(issues, `counter metrics should not have "_total" suffix`)
	}

	return issues
}

// metricUnits attempts to detect known unit types used as part of a metric name,
// e.g. "foo_bytes_total" or "bar_baz_milligrams".
func getMetricUnit(m string) (unit string, base string, ok bool) {
	ss := strings.Split(m, "_")

	for unit, base := range units {
		// Also check for "no prefix".
		for _, p := range append(unitPrefixes, "") {
			for _, s := range ss {
				// Attempt to explicitly match a known unit with a known prefix,
				// as some words may look like "units" when matching suffix.
				//
				// As an example, "thermometers" should not match "meters", but
				// "kilometers" should.
				if s == p+unit {
					return p + unit, base, true
				}
			}
		}
	}

	return "", "", false
}

func lintMetricUnit(name string) (issues []string) {
	unit, base, ok := getMetricUnit(name)
	if !ok {
		// No known units detected.
		return nil
	}

	// Unit is already a base unit.
	if unit == base {
		return nil
	}

	issues = append(issues, fmt.Sprintf("use base unit %q instead of %q", base, unit))

	return issues
}

func hasBucketSuffix(name string) bool{
	return strings.HasSuffix(name, "_bucket")
}

func lintNonHistogramNoBucket(name string) (issues []string) {
	if hasBucketSuffix(name) {
		issues = append(issues, `non-histogram metrics should not have "_bucket" suffix`)
	}

	return issues
}

func hasCountSuffix(name string) bool {
	return strings.HasSuffix(name, "_count")
}

func lintNonHistogramSummaryNoCount(name string) (issues []string) {
	if hasCountSuffix(name) {
		issues = append(issues, `non-histogram and non-summary metrics should not have "_count" suffix`)
	}

	return issues
}

func lintNonHistogramSummaryNoSum(name string) (issues []string) {
	if strings.HasSuffix(name, NameSuffixSum) {
		issues = append(issues, `non-histogram and non-summary metrics should not have "_sum" suffix`)
	}

	return issues
}

func lintNonHistogramNoLabelLe(constLabels prometheus.Labels, labelNames []string) (issues []string) {
	for ln, _ := range constLabels {
		if ln == LabelLe {
			issues = append(issues, `non-histogram metrics should not have "le" label`)
		}
	}

	for _, ln := range labelNames {
		if ln == LabelLe {
			issues = append(issues, `non-histogram metrics should not have "le" label`)
		}
	}

	return issues
}

func lintNonSummaryNoLabelQuantile(constLabels prometheus.Labels, labelNames []string) (issues []string) {
	for ln, _ := range constLabels {
		if ln == LabelQuantile {
			issues = append(issues, `non-summary metrics should not have "quantile" label`)
		}
	}

	return issues
}

func lintNoMetricTypeInName(name string) (issues []string) {
	n := strings.ToLower(name)

	for i, t := range dto.MetricType_name {
		if i == int32(dto.MetricType_UNTYPED) {
			continue
		}

		typename := strings.ToLower(t)
		if strings.Contains(n, "_"+typename+"_") || strings.HasSuffix(n, "_"+typename) {
			issues = append(issues, fmt.Sprintf(`metric name should not include type '%s'`, typename))
		}
	}

	return issues
}

func lintReservedChars(name string) (issues []string) {
	if strings.Contains(name, ":") {
		issues = append(issues, "metric names should not contain ':'")
	}

	return issues
}

func lintNameCamelCase(name string) (issues []string) {
	if camelCase.FindString(name) != "" {
		issues = append(issues, "metric names should be written in 'snake_case' not 'camelCase'")
	}

	return issues
}

// TODO(RainbowMango): Should check label value?
func lintLabelNameCamelCase(constLabels prometheus.Labels, labelNames []string) (issues []string) {
	for ln, _ := range constLabels {
		if camelCase.FindString(ln) != "" {
			issues = append(issues, "label names should be written in 'snake_case' not 'camelCase'")
		}
	}

	for _, ln := range labelNames {
		if camelCase.FindString(ln) != "" {
			issues = append(issues, "label names should be written in 'snake_case' not 'camelCase'")
		}
	}

	return issues
}

// lintUnitAbbreviations detects abbreviated units in the metric name.
// TODO(RainbowMango): It'd be better to return which abbreviated unit contains in name. Check with promlint guys.
func lintUnitAbbreviations(name string) (issues []string) {
	n := strings.ToLower(name)
	for _, s := range unitAbbreviations {
		if strings.Contains(n, "_"+s+"_") || strings.HasSuffix(n, "_"+s) {
			issues = append(issues, "metric names should not contain abbreviated units")
		}
	}

	return issues
}

func commonLint(opts interface{}) (issues []string) {
	switch opts.(type) {
	case prometheus.Opts: // prometheus.CounterOpts and prometheus.GaugeOpts share the type.
		counterGagueOpts := opts.(prometheus.Opts)
		fqName := prometheus.BuildFQName(counterGagueOpts.Namespace, counterGagueOpts.Subsystem, counterGagueOpts.Name)

		issues = append(issues, lintHelp(counterGagueOpts.Help)...) // metrics should contains help.
		issues = append(issues, lintMetricUnit(fqName)...) // name should use standard units.

	case prometheus.CounterOpts:
		counterOpts := opts.(prometheus.CounterOpts)
		issues = append(issues, lintNoMetricTypeInName(counterOpts.Name)...)
		issues = append(issues, lintReservedChars(counterOpts.Name)...)
		issues = append(issues, lintNameCamelCase(counterOpts.Name)...)
		issues = append(issues, lintUnitAbbreviations(counterOpts.Name)...)
		issues = append(issues, lintNonHistogramNoBucket(counterOpts.Name)...)
		issues = append(issues, lintNonHistogramSummaryNoCount(counterOpts.Name)...)
		issues = append(issues, lintNonHistogramSummaryNoSum(counterOpts.Name)...)
		issues = append(issues, lintNonHistogramNoLabelLe(counterOpts.ConstLabels, nil)...)
		issues = append(issues, lintNonSummaryNoLabelQuantile(counterOpts.ConstLabels, nil)...)
		issues = append(issues, lintLabelNameCamelCase(counterOpts.ConstLabels, nil)...)
	case prometheus.GaugeOpts:
		gaugeOpts := opts.(prometheus.GaugeOpts)
		issues = append(issues, lintHelp(gaugeOpts.Help)...)
		issues = append(issues, lintNoMetricTypeInName(gaugeOpts.Name)...)
		issues = append(issues, lintReservedChars(gaugeOpts.Name)...)
		issues = append(issues, lintMetricUnit(gaugeOpts.Name)...)
		issues = append(issues, lintNameCamelCase(gaugeOpts.Name)...)
		issues = append(issues, lintUnitAbbreviations(gaugeOpts.Name)...)
		issues = append(issues, lintNonHistogramNoBucket(gaugeOpts.Name)...)
		issues = append(issues, lintNonHistogramSummaryNoCount(gaugeOpts.Name)...)
		issues = append(issues, lintNonHistogramSummaryNoSum(gaugeOpts.Name)...)
		issues = append(issues, lintNonHistogramNoLabelLe(gaugeOpts.ConstLabels, nil)...)
		issues = append(issues, lintNonSummaryNoLabelQuantile(gaugeOpts.ConstLabels, nil)...)
		issues = append(issues, lintLabelNameCamelCase(gaugeOpts.ConstLabels, nil)...)
	case prometheus.HistogramOpts:
		histogramOpts := opts.(prometheus.HistogramOpts)
		fqName := prometheus.BuildFQName(histogramOpts.Namespace, histogramOpts.Subsystem, histogramOpts.Name)

		issues = append(issues, lintHelp(histogramOpts.Help)...)
		issues = append(issues, lintMetricUnit(fqName)...)

		// TODO: delete me if no items below
		issues = append(issues, lintNoMetricTypeInName(fqName)...)
		issues = append(issues, lintReservedChars(fqName)...)
		issues = append(issues, lintNameCamelCase(fqName)...)
		issues = append(issues, lintUnitAbbreviations(fqName)...)

		issues = append(issues, lintNonSummaryNoLabelQuantile(histogramOpts.ConstLabels, nil)...)
		issues = append(issues, lintLabelNameCamelCase(histogramOpts.ConstLabels, nil)...)
	case prometheus.SummaryOpts:
		summaryOpts := opts.(prometheus.SummaryOpts)
		fqName := prometheus.BuildFQName(summaryOpts.Namespace, summaryOpts.Subsystem, summaryOpts.Name)

		issues = append(issues, lintHelp(summaryOpts.Help)...)
		issues = append(issues, lintMetricUnit(fqName)...)

		// TODO: delete me if no items below
		issues = append(issues, lintNoMetricTypeInName(fqName)...)
		issues = append(issues, lintReservedChars(fqName)...)
		issues = append(issues, lintNonHistogramNoBucket(fqName)...)
		issues = append(issues, lintNameCamelCase(fqName)...)
		issues = append(issues, lintUnitAbbreviations(fqName)...)
		issues = append(issues, lintNonHistogramNoLabelLe(summaryOpts.ConstLabels, nil)...)
		issues = append(issues, lintLabelNameCamelCase(summaryOpts.ConstLabels, nil)...)
	default:
		panic(fmt.Sprintf("unknow metric type: %T", opts))
	}

	return issues
}