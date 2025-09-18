// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher_test

import (
	"testing"

	"github.com/joshdk/krf/matcher"
)

func TestDiffMatcher(t *testing.T) {
	t.Parallel()

	testMatcher(t, []spec{
		{
			title:   "diff resources",
			matcher: must(matcher.NewDiffMatcher("testdata/resources.txt")),
			matches: []string{
				"Deployment/nginx-deployment",
				"Pod/test-pod",
			},
		},
	})
}
