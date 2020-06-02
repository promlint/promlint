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

import "github.com/prometheus/client_golang/prometheus"

// A Problem is an issue detected by lint.
type Problem struct {
	// The name of the metric indicated by this Problem.
	MetricName string

	// A description of the issue for this Problem.
	ProblemDesc string
}

func LintCounter(counterOpts prometheus.CounterOpts) (problems []Problem) {
	subProblems := CommonLint(prometheus.Opts(counterOpts))
	if len(subProblems) > 0 {
		problems = append(problems, subProblems...)
	}

	return
}