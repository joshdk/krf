// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher_test

import (
	"testing"

	"github.com/joshdk/krf/matcher"
)

func TestExecMatcher(t *testing.T) {
	t.Parallel()

	testMatcher(t, []spec{
		{
			title:   "always false",
			matcher: must(matcher.NewExecMatcher(`false`)),
			matches: []string{},
		},
		{
			title:   "always true",
			matcher: must(matcher.NewExecMatcher(`true`)),
			matches: []string{
				"ClusterRoleBinding/read-secrets-global",
				"ConfigMap/my-configmap",
				"Deployment/nginx-deployment",
				"Pod/test-pod",
				"Service/my-service",
			},
		},
		{
			title:   "jq example",
			matcher: must(matcher.NewExecMatcher(`jq --exit-status '.spec.replicas == 3'`)),
			matches: []string{"Deployment/nginx-deployment"},
		},
		{
			title:   "custom script",
			matcher: must(matcher.NewExecMatcher(`./testdata/matcher.sh v1 Service my-service custom-app`)),
			matches: []string{"Service/my-service"},
		},
	})
}
