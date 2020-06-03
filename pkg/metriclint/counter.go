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

import "strings"

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