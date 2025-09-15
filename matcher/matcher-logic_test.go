// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher_test

import (
	"testing"

	"github.com/joshdk/krf/matcher"
)

func TestNotMatcher(t *testing.T) {
	t.Parallel()

	testMatcher(t, []spec{
		{
			title: "not matcher",
			matcher: matcher.NotMatcher(
				must(matcher.NewAPIVersionMatcher("v1")),
			),
			matches: []string{
				"ClusterRoleBinding/read-secrets-global",
				"Deployment/nginx-deployment",
			},
		},
	})
}

func TestAnyMatcher(t *testing.T) {
	t.Parallel()

	testMatcher(t, []spec{
		{
			title: "any matcher",
			matcher: func() matcher.Matcher {
				am := &matcher.AnyMatcher{}
				am.Append(must(matcher.NewNameMatcher("my-service")))
				am.Append(must(matcher.NewNameMatcher("my-configmap")))

				return am
			}(),
			matches: []string{
				"ConfigMap/my-configmap",
				"Service/my-service",
			},
		},
		{
			title: "any matcher partial",
			matcher: func() matcher.Matcher {
				am := &matcher.AnyMatcher{}
				am.Append(must(matcher.NewNameMatcher("my-service")))
				am.Append(must(matcher.NewNameMatcher("not-a-real-name")))

				return am
			}(),
			matches: []string{
				"Service/my-service",
			},
		},
	})
}

func TestAllMatcher(t *testing.T) {
	t.Parallel()

	testMatcher(t, []spec{
		{
			title: "all matcher",
			matcher: func() matcher.Matcher {
				am := &matcher.AllMatcher{}
				am.Append(must(matcher.NewNameMatcher("my-*")))
				am.Append(must(matcher.NewLabelMatcher("app")))

				return am
			}(),
			matches: []string{
				"Service/my-service",
			},
		},
		{
			title: "all matcher failed",
			matcher: func() matcher.Matcher {
				am := &matcher.AllMatcher{}
				am.Append(must(matcher.NewNameMatcher("test-pod")))
				am.Append(must(matcher.NewNamespaceMatcher("custom-app")))

				return am
			}(),
			matches: []string{},
		},
	})
}
