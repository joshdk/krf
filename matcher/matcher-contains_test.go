// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher_test

import (
	"testing"

	"github.com/joshdk/krf/matcher"
)

func TestContainsMatcher(t *testing.T) {
	t.Parallel()

	testMatcher(t, []spec{
		{
			title:   "image substring",
			matcher: must(matcher.NewContainsMatcher("mage: gcr.io/goog")),
			matches: []string{"Pod/test-pod"},
		},
		{
			title:   "port substring",
			matcher: must(matcher.NewContainsMatcher("ort: 80")),
			matches: []string{
				"Deployment/nginx-deployment",
				"Service/my-service",
			},
		},
	})
}
