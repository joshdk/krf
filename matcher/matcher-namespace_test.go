// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher_test

import (
	"testing"

	"github.com/joshdk/krf/matcher"
)

func TestNamespaceMatcher(t *testing.T) {
	t.Parallel()

	testMatcher(t, []spec{
		{
			title:   "single full namespace",
			matcher: must(matcher.NewNamespaceMatcher("custom-app")),
			matches: []string{
				"ConfigMap/my-configmap",
				"Service/my-service",
			},
		},
		{
			title:   "implicit wildcard suffix",
			matcher: must(matcher.NewNamespaceMatcher("custom-app-")),
			matches: []string{
				"Pod/test-pod",
			},
		},
		{
			title:   "wildcard",
			matcher: must(matcher.NewNamespaceMatcher("*")),
			matches: []string{
				"ConfigMap/my-configmap",
				"Deployment/nginx-deployment",
				"Pod/test-pod",
				"Service/my-service",
			},
		},
	})
}
