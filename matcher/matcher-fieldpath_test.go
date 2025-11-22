// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher_test

import (
	"testing"

	"github.com/joshdk/krf/matcher"
)

func TestFieldPathMatcher(t *testing.T) {
	t.Parallel()

	testMatcher(t, []spec{
		{
			title:   "compare value",
			matcher: must(matcher.NewFieldPathMatcher(".spec.ports[0].port=80")),
			matches: []string{"Service/my-service"},
		},
		{
			title:   "compare null",
			matcher: must(matcher.NewFieldPathMatcher(".spec.securityContext=null")),
			matches: []string{"Pod/test-pod"},
		},
		{
			title:   "nonexistent path",
			matcher: must(matcher.NewFieldPathMatcher(".foo.bar")),
			matches: []string{},
		},
		{
			title:   "existence of path",
			matcher: must(matcher.NewFieldPathMatcher(".data.username")),
			matches: []string{"ConfigMap/my-configmap"},
		},
		{
			title:   "wildcard value",
			matcher: must(matcher.NewFieldPathMatcher(".data.username=*-admin")),
			matches: []string{"ConfigMap/my-configmap"},
		},
		{
			title:   "check for existence of a container with a specific name",
			matcher: must(matcher.NewFieldPathMatcher(".spec.template.spec.containers.[name=nginx]")),
			matches: []string{"Deployment/nginx-deployment"},
		},
		{
			title:   "check for a specific volume mount",
			matcher: must(matcher.NewFieldPathMatcher(".spec.containers.[name=test-container].volumeMounts.[name=config-volume].mountPath=/etc/conf*")),
			matches: []string{"Pod/test-pod"},
		},
	})
}
