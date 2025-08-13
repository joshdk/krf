// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher

import (
	"github.com/joshdk/krf/resources"
)

// NewNameMatcher matches resources.Resource instances based on the given name.
func NewNameMatcher(name string) Matcher {
	return nameMatcher{name: name}
}

type nameMatcher struct {
	name string
}

func (m nameMatcher) Matches(item resources.Resource) bool {
	return m.name == item.GetName()
}
