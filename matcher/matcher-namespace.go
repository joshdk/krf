// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher

import (
	"github.com/gobwas/glob"

	"github.com/joshdk/krf/resources"
)

// NewNamespaceMatcher matches resources.Resource instances based on the given
// namespace.
func NewNamespaceMatcher(namespace string) (Matcher, error) {
	namespaceGlob, err := asGlob(namespace)
	if err != nil {
		return nil, err
	}

	return namespaceMatcher{namespaceGlob: namespaceGlob}, nil
}

type namespaceMatcher struct {
	namespaceGlob glob.Glob
}

func (m namespaceMatcher) Matches(item resources.Resource) bool {
	// Prevent wildcard matches against empty namespace values.
	if item.GetNamespace() == "" {
		return false
	}

	return m.namespaceGlob.Match(item.GetNamespace())
}
