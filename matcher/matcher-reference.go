// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher

import (
	"strings"

	"github.com/gobwas/glob"

	"github.com/joshdk/krf/references"
	"github.com/joshdk/krf/resolver"
	"github.com/joshdk/krf/resources"
)

// NewReferenceMatcher matches resources.Resource instances that reference the
// given named resource. A name glob like "my-secret-*" can be used to match
// any reference to a resource with that name, or can include a kind like
// "cm/my-configmap-*" to match any reference to a resource with that kind and
// name.
func NewReferenceMatcher(reference string) (Matcher, error) {
	var kind, name string

	// Split a string like "cm/my-configmap" or "my-configmap" into the kind
	// and name.
	parts := strings.SplitN(reference, "/", 2)
	if len(parts) == 2 {
		kind, name = parts[0], parts[1]
	} else {
		kind, name = "", parts[0]
	}

	nameGlob, err := asGlob(name)
	if err != nil {
		return nil, err
	}

	var kinds []string

	if kind != "" {
		for _, resource := range resolver.LookupAlias(kind) {
			kinds = append(kinds, resource.Kind)
		}

		if len(kinds) == 0 {
			kinds = []string{kind}
		}
	}

	return referenceMatcher{kinds, nameGlob}, nil
}

type referenceMatcher struct {
	kinds    []string
	nameGlob glob.Glob
}

func (m referenceMatcher) Matches(item resources.Resource) bool {
	// References to any kind.
	if len(m.kinds) == 0 {
		return references.References(item.Unstructured, "", func(s string) bool {
			return m.nameGlob.Match(s)
		})
	}

	// References to only the requested kind(s).
	for _, kind := range m.kinds {
		if references.References(item.Unstructured, kind, m.nameGlob.Match) {
			return true
		}
	}

	return false
}
