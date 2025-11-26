// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher_test

import (
	"testing"

	"github.com/joshdk/krf/matcher"
)

func TestRegoMatcher(t *testing.T) {
	t.Parallel()

	testMatcher(t, []spec{
		{
			title:   "simple rego module",
			matcher: must(matcher.NewRegoMatcher("testdata/simple.rego")),
			matches: []string{
				"Service/my-service",
				"Pod/test-pod",
			},
		},
		{
			title:   "complex rego module",
			matcher: must(matcher.NewRegoMatcher("testdata/complex.rego")),
			matches: []string{
				"Pod/test-pod",
			},
		},
	})
}
