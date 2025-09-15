// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher_test

import (
	"testing"

	"github.com/joshdk/krf/matcher"
)

func TestKindMatcher(t *testing.T) {
	t.Parallel()

	testMatcher(t, []spec{
		{
			title:   "canonical kind",
			matcher: must(matcher.NewKindMatcher("Deployment")),
			matches: []string{"Deployment/nginx-deployment"},
		},
		{
			title:   "plural",
			matcher: must(matcher.NewKindMatcher("deployments")),
			matches: []string{"Deployment/nginx-deployment"},
		},
		{
			title:   "kind alias",
			matcher: must(matcher.NewKindMatcher("deploy")),
			matches: []string{"Deployment/nginx-deployment"},
		},
		{
			title:   "ignore capitalization",
			matcher: must(matcher.NewKindMatcher("DEPL*")),
			matches: []string{"Deployment/nginx-deployment"},
		},
	})
}
