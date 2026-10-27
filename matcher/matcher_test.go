// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher_test

import (
	"fmt"
	"testing"
	"testing/synctest"
	"time"

	"github.com/joshdk/krf/config"
	"github.com/joshdk/krf/matcher"
	"github.com/joshdk/krf/resolver"
	"github.com/joshdk/krf/resources"
)

// testNow is an arbitrary, consistent point in time to use for testing,
// independent of the current time.
var testNow, _ = time.ParseInLocation(time.RFC3339, "2026-10-24T00:00:00Z", time.UTC)

var testResources []resources.Resource

func init() { //nolint:gochecknoinits
	// Load the builtin configuration file so that resource metadata can be
	// utilized in tests.
	cfg, err := config.Load("../config/files/configuration.yaml")
	if err != nil {
		panic(err)
	}

	resolver.Init(cfg.Resources)

	// Decode all resources in the testdata directory for (re)use across all
	// matcher tests.
	err = resources.Decode("./testdata", func(item resources.Resource) {
		testResources = append(testResources, item)
	})
	if err != nil {
		panic(err)
	}
}

// spec represents a single matcher test case.
type spec struct {
	title   string
	matcher matcher.Matcher
	matches []string
}

// testMatcher evaluates all given specs against the shared set of decoded resources.
func testMatcher(t *testing.T, tests []spec) {
	t.Helper()

	for _, test := range tests {
		testRun(t, test, func(t *testing.T) {
			t.Helper()

			expected := make(map[string]struct{})

			for _, name := range test.matches {
				expected[name] = struct{}{}
			}

			actual := make(map[string]struct{})

			for _, item := range testResources {
				if test.matcher.Matches(item) {
					name := fmt.Sprintf("%s/%s", item.GetKind(), item.GetName())
					actual[name] = struct{}{}
				}
			}

			// Check if the current test case had any expected resources that
			// were not actually matched.
			for expectedName := range expected {
				if _, found := actual[expectedName]; !found {
					t.Errorf("resource %s was expected but not matched", expectedName)
				}
			}

			// Check if the current test case had any actually matched
			// resources that were not expected.
			for actualName := range actual {
				if _, found := expected[actualName]; !found {
					t.Errorf("resource %s was matched but not expected", actualName)
				}
			}
		})
	}
}

// testRun runs the given test function inside a new time bubble.
func testRun(t *testing.T, test spec, testFn func(t *testing.T)) {
	t.Helper()

	t.Run(test.title, func(t *testing.T) {
		t.Parallel()

		synctest.Test(t, func(t *testing.T) {
			time.Sleep(time.Until(testNow))

			testFn(t)
		})
	})
}

func must(m matcher.Matcher, err error) matcher.Matcher {
	if err != nil {
		panic(err)
	}

	return m
}
