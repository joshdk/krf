// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher_test

import (
	"testing"

	"github.com/joshdk/krf/matcher"
)

func TestReferenceMatcher(t *testing.T) {
	t.Parallel()

	testMatcher(t, []spec{
		{
			title:   "reference resource name",
			matcher: must(matcher.NewReferenceMatcher("example-config")),
			matches: []string{"Pod/test-pod"},
		},
		{
			title:   "implicit wildcard suffix",
			matcher: must(matcher.NewReferenceMatcher("example-")),
			matches: []string{
				"Pod/test-pod",
				"Deployment/nginx-deployment",
			},
		},
		{
			title:   "resource reference kind",
			matcher: must(matcher.NewReferenceMatcher("Secret/example-")),
			matches: []string{
				"Deployment/nginx-deployment",
			},
		},
		{
			title:   "resource reference kind alias",
			matcher: must(matcher.NewReferenceMatcher("cm/example-")),
			matches: []string{"Pod/test-pod"},
		},
	})
}
