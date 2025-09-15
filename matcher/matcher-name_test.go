// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher_test

import (
	"testing"

	"github.com/joshdk/krf/matcher"
)

func TestNameMatcher(t *testing.T) {
	t.Parallel()

	testMatcher(t, []spec{
		{
			title:   "single full name",
			matcher: must(matcher.NewNameMatcher("nginx-deployment")),
			matches: []string{"Deployment/nginx-deployment"},
		},
		{
			title:   "implicit wildcard suffix",
			matcher: must(matcher.NewNameMatcher("nginx-")),
			matches: []string{"Deployment/nginx-deployment"},
		},
		{
			title:   "implicit wildcard prefix",
			matcher: must(matcher.NewNameMatcher("-deployment")),
			matches: []string{"Deployment/nginx-deployment"},
		},
		{
			title:   "explicit wildcard",
			matcher: must(matcher.NewNameMatcher("my*")),
			matches: []string{
				"ConfigMap/my-configmap",
				"Service/my-service",
			},
		},
		{
			title:   "wildcard",
			matcher: must(matcher.NewNameMatcher("*")),
			matches: []string{
				"ClusterRoleBinding/read-secrets-global",
				"ConfigMap/my-configmap",
				"Deployment/nginx-deployment",
				"Pod/test-pod",
				"Service/my-service",
			},
		},
	})
}
