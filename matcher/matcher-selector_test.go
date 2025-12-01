// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher_test

import (
	"testing"

	"github.com/joshdk/krf/matcher"
)

func TestSelectorMatcher(t *testing.T) {
	t.Parallel()

	testMatcher(t, []spec{
		{
			title:   "existence of a label",
			matcher: must(matcher.NewSelectorMatcher("app")),
			matches: []string{
				"Deployment/nginx-deployment",
				"Service/my-service",
			},
		},
		{
			title:   "existence of multiple labels",
			matcher: must(matcher.NewSelectorMatcher("app,component")),
			matches: []string{
				"Deployment/nginx-deployment",
				"Service/my-service",
			},
		},
		{
			title:   "label with a value",
			matcher: must(matcher.NewSelectorMatcher("app=myapp")),
			matches: []string{"Service/my-service"},
		},
		{
			title:   "multiple labels with values",
			matcher: must(matcher.NewSelectorMatcher("app=myapp,component=api")),
			matches: []string{"Service/my-service"},
		},
	})
}
