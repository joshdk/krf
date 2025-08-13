// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher

import (
	"github.com/joshdk/krf/resources"
)

// NewNamespaceMatcher matches resources.Resource instances based on the given
// namespace.
func NewNamespaceMatcher(namespace string) Matcher {
	return namespaceMatcher{namespace: namespace}
}

type namespaceMatcher struct {
	namespace string
}

func (m namespaceMatcher) Matches(item resources.Resource) bool {
	return m.namespace == item.GetNamespace()
}
