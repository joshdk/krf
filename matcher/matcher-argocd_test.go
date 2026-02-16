// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher_test

import (
	"testing"

	"github.com/joshdk/krf/matcher"
)

func TestArgoCDMatcher(t *testing.T) {
	t.Parallel()

	testMatcher(t, []spec{
		{
			title:   "literal app name",
			matcher: must(matcher.NewArgoCDMatcher("app")),
			matches: []string{
				"Service/my-service",
			},
		},
		{
			title:   "app name with wildcard suffix",
			matcher: must(matcher.NewArgoCDMatcher("app*")),
			matches: []string{
				"Deployment/nginx-deployment",
				"Service/my-service",
			},
		},
		{
			title:   "wildcard",
			matcher: must(matcher.NewArgoCDMatcher("*")),
			matches: []string{
				"Deployment/nginx-deployment",
				"Service/my-service",
			},
		},
		{
			title:   "wildcard with suffix",
			matcher: must(matcher.NewArgoCDMatcher("*-dev")),
			matches: []string{
				"Deployment/nginx-deployment",
			},
		},
	})
}
