// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher

import (
	"strings"

	"github.com/joshdk/krf/resources"
)

// NewPathMatcher matches resources.Resource instances based on the given file
// path substring.
//
// For example, a resource decoded from the file
// `kustomize/environments/production/deployment.yaml` would be matched by the
// input `environments/production` or `deployment.yaml`.
func NewPathMatcher(path string) Matcher {
	return pathMatcher{path: path}
}

type pathMatcher struct {
	path string
}

func (m pathMatcher) Matches(item resources.Resource) bool {
	return strings.Contains(item.GetFilename(), m.path)
}
