// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher_test

import (
	"testing"

	"github.com/joshdk/krf/matcher"
)

func TestAnnotationMatcher(t *testing.T) {
	t.Parallel()

	testMatcher(t, []spec{
		{
			title:   "full annotation and value",
			matcher: must(matcher.NewAnnotationMatcher("kubernetes.io/description=proxy")),
			matches: []string{"Deployment/nginx-deployment"},
		},
		{
			title:   "annotation with wildcard value",
			matcher: must(matcher.NewAnnotationMatcher("kubernetes.io/description=*")),
			matches: []string{
				"ConfigMap/my-configmap",
				"Deployment/nginx-deployment",
				"Pod/test-pod",
			},
		},
		{
			title:   "existence of annotation",
			matcher: must(matcher.NewAnnotationMatcher("kubernetes.io/description")),
			matches: []string{
				"ConfigMap/my-configmap",
				"Deployment/nginx-deployment",
				"Pod/test-pod",
			},
		},
		{
			title:   "wildcard annotation",
			matcher: must(matcher.NewAnnotationMatcher("kubernetes.io/*")),
			matches: []string{
				"ConfigMap/my-configmap",
				"Deployment/nginx-deployment",
				"Pod/test-pod",
			},
		},
	})
}
