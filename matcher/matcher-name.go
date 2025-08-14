// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher

import (
	"github.com/gobwas/glob"

	"github.com/joshdk/krf/resources"
)

// NewNameMatcher matches resources.Resource instances based on the given name.
func NewNameMatcher(name string) (Matcher, error) {
	nameGlob, err := asGlob(name)
	if err != nil {
		return nil, err
	}

	return nameMatcher{nameGlob: nameGlob}, nil
}

type nameMatcher struct {
	nameGlob glob.Glob
}

func (m nameMatcher) Matches(item resources.Resource) bool {
	return m.nameGlob.Match(item.GetName())
}
