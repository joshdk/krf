// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher_test

import (
	"testing"

	"github.com/joshdk/krf/matcher"
)

func TestNewerThanMatcher(t *testing.T) {
	t.Parallel()

	testMatcher(t, []spec{
		{
			title:   "created after right now",
			matcher: must(matcher.NewCreatedAfterMatcher("0")),
			matches: []string{},
		},
		{
			title:   "created after 24 hours and 1 minute ago",
			matcher: must(matcher.NewCreatedAfterMatcher("24h1m")),
			matches: []string{
				"Service/my-service", // 24h old
			},
		},
		{
			title:   "created after 48 hours and 1 minute ago",
			matcher: must(matcher.NewCreatedAfterMatcher("48h1m")),
			matches: []string{
				"Pod/test-pod",       // 48h old
				"Service/my-service", // 24h old
			},
		},
		{
			title:   "created after a specific timestamp",
			matcher: must(matcher.NewCreatedAfterMatcher("2026-10-21T23:59:00Z")),
			matches: []string{
				"Pod/test-pod",       // 48h old
				"Service/my-service", // 24h old
			},
		},
		{
			title:   "created after 72 hours and 1 minute ago",
			matcher: must(matcher.NewCreatedAfterMatcher("72h1m")),
			matches: []string{
				"Deployment/nginx-deployment", // 72h old
				"Pod/test-pod",                // 48h old
				"Service/my-service",          // 24h old
			},
		},
	})
}

func TestOlderThanMatcher(t *testing.T) {
	t.Parallel()

	testMatcher(t, []spec{
		{
			title:   "created before 24 hours less 1 minute ago",
			matcher: must(matcher.NewCreatedBeforeMatcher("23h59m")),
			matches: []string{
				"Deployment/nginx-deployment", // 72h old
				"Pod/test-pod",                // 48h old
				"Service/my-service",          // 24h old
			},
		},
		{
			title:   "created before 48 hours less 1 minute ago",
			matcher: must(matcher.NewCreatedBeforeMatcher("47h59m")),
			matches: []string{
				"Deployment/nginx-deployment", // 72h old
				"Pod/test-pod",                // 48h old
			},
		},
		{
			title:   "created before unix time",
			matcher: must(matcher.NewCreatedBeforeMatcher("1792627140")),
			matches: []string{
				"Deployment/nginx-deployment", // 72h old
			},
		},
		{
			title:   "created before 72 hours less 1 minute ago",
			matcher: must(matcher.NewCreatedBeforeMatcher("71h59m")),
			matches: []string{
				"Deployment/nginx-deployment", // 72h old
			},
		},
		{
			title:   "created before many hours ago",
			matcher: must(matcher.NewCreatedBeforeMatcher("100h")),
		},

		{
			title: "clamp between a created before and created after time",
			matcher: func() matcher.Matcher {
				am := &matcher.AllMatcher{}
				am.Append(must(matcher.NewCreatedBeforeMatcher("47h59m")))
				am.Append(must(matcher.NewCreatedAfterMatcher("48h1m")))

				return am
			}(),
			matches: []string{
				"Pod/test-pod", // 48h old
			},
		},
	})
}
