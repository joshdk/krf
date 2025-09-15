// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher_test

import (
	"testing"

	"github.com/joshdk/krf/matcher"
)

func TestAPIVersionMatcherMatcher(t *testing.T) {
	t.Parallel()

	testMatcher(t, []spec{
		{
			title:   "full apiversion",
			matcher: must(matcher.NewAPIVersionMatcher("v1")),
			matches: []string{
				"ConfigMap/my-configmap",
				"Pod/test-pod",
				"Service/my-service",
			},
		},
		{
			title:   "wildcard prefix",
			matcher: must(matcher.NewAPIVersionMatcher("*/v1")),
			matches: []string{
				"ClusterRoleBinding/read-secrets-global",
				"Deployment/nginx-deployment",
			},
		},
		{
			title:   "wildcard suffix",
			matcher: must(matcher.NewAPIVersionMatcher("rbac*")),
			matches: []string{
				"ClusterRoleBinding/read-secrets-global",
			},
		},
	})
}
