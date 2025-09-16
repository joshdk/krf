// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher_test

import (
	"testing"

	"github.com/joshdk/krf/matcher"
)

func TestJsonpathMatcher(t *testing.T) {
	t.Parallel()

	testMatcher(t, []spec{
		{
			title:   "compare value",
			matcher: must(matcher.NewJsonpathMatcher(".spec.ports[0].port=80")),
			matches: []string{"Service/my-service"},
		},
		{
			title:   "compare null",
			matcher: must(matcher.NewJsonpathMatcher(".spec.securityContext=null")),
			matches: []string{"Pod/test-pod"},
		},
		{
			title:   "nonexistent path",
			matcher: must(matcher.NewJsonpathMatcher(".foo.bar")),
			matches: []string{},
		},
		{
			title:   "existence of path",
			matcher: must(matcher.NewJsonpathMatcher(".data.username")),
			matches: []string{"ConfigMap/my-configmap"},
		},
		{
			title:   "wildcard value",
			matcher: must(matcher.NewJsonpathMatcher(".data.username=*-admin")),
			matches: []string{"ConfigMap/my-configmap"},
		},
	})
}
