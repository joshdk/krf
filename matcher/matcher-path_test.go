// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher_test

import (
	"testing"

	"github.com/joshdk/krf/matcher"
)

func TestPathMatcher(t *testing.T) {
	t.Parallel()

	testMatcher(t, []spec{
		{
			title:   "files under subdir",
			matcher: must(matcher.NewPathMatcher("/subdir/")),
			matches: []string{
				"ClusterRoleBinding/read-secrets-global",
				"Deployment/nginx-deployment",
				"ConfigMap/my-configmap",
			},
		},
		{
			title:   "files under subsubdir",
			matcher: must(matcher.NewPathMatcher("/subdir/subsubdir")),
			matches: []string{"ConfigMap/my-configmap"},
		},
	})
}
