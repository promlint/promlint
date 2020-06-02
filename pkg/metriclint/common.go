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

func lintHelp(help string) (issues []string) {
	if len(help) == 0 {
		issues = append(issues, "no help text")
	}

	return
}

func CommonLint(opts prometheus.Opts) (problems []Problem) {
	issues := lintHelp(opts.Help)
	for _, issue := range issues {
		problems = append(problems, Problem{MetricName:prometheus.BuildFQName(opts.Namespace, opts.Subsystem, opts.Name), ProblemDesc:issue})
	}

	return
}