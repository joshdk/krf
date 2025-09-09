// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher

import (
	"strings"

	"github.com/joshdk/krf/resources"
)

// NewPatchMatcher matches resources.Resource instances that were decoded from
// a kustomize patch file (the convention being a filename with the suffix
// `.patch.yaml`).
//
// For example, a resource decoded from the file
// `kustomize/environments/production/deployment.patch.yaml` would be matched.
func NewPatchMatcher() Matcher {
	return patchMatcher{}
}

type patchMatcher struct{}

func (m patchMatcher) Matches(item resources.Resource) bool {
	return strings.HasSuffix(item.GetFilename(), ".patch.yaml")
}
