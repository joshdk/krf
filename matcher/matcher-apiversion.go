// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher

import (
	"github.com/gobwas/glob"

	"github.com/joshdk/krf/resources"
)

// NewAPIVersionMatcher matches resources.Resource instances based on the given
// apiversion glob.
func NewAPIVersionMatcher(apiVersion string) (Matcher, error) {
	apiVersionGlob, err := glob.Compile(apiVersion)
	if err != nil {
		return nil, err
	}

	return apiVersionMatcher{apiVersionGlob: apiVersionGlob}, nil
}

type apiVersionMatcher struct {
	apiVersionGlob glob.Glob
}

func (m apiVersionMatcher) Matches(item resources.Resource) bool {
	return m.apiVersionGlob.Match(item.GetAPIVersion())
}
