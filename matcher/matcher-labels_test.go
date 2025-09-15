// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher_test

import (
	"testing"

	"github.com/joshdk/krf/matcher"
)

func TestLabelMatcher(t *testing.T) {
	t.Parallel()

	testMatcher(t, []spec{
		{
			title:   "full label and value",
			matcher: must(matcher.NewLabelMatcher("app=myapp")),
			matches: []string{"Service/my-service"},
		},
		{
			title:   "label with wildcard value",
			matcher: must(matcher.NewLabelMatcher("app=*")),
			matches: []string{
				"Deployment/nginx-deployment",
				"Service/my-service",
			},
		},
		{
			title:   "existence of label",
			matcher: must(matcher.NewLabelMatcher("app")),
			matches: []string{
				"Deployment/nginx-deployment",
				"Service/my-service",
			},
		},
		{
			title:   "wildcard label",
			matcher: must(matcher.NewLabelMatcher("app.kubernetes.io/*")),
			matches: []string{"ConfigMap/my-configmap"},
		},
	})
}
