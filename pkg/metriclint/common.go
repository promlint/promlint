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
	"strings"

	"github.com/prometheus/client_golang/prometheus"
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

func lintHelp(help string) (issues []string) {
	if len(help) == 0 {
		issues = append(issues, "no help text")
	}

	return
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

func CommonLint(opts prometheus.Opts) (issues []string) {
	issues = append(issues, lintHelp(opts.Help)...)
	issues = append(issues, lintMetricUnit(opts.Name)...)

	return
}
