// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher

import (
	"strings"

	"github.com/gobwas/glob"

	"github.com/joshdk/krf/resolver"
	"github.com/joshdk/krf/resources"
)

// NewKindMatcher matches resources.Resource instances based on the given kind
// which might be the kind verbatim, a glob, or a kind alias.
//
// For example, a resource of kind `Service` would be matched by the input
// "Service", "svc", or "sv*".
func NewKindMatcher(kind string) (Matcher, error) {
	kindGlob, err := glob.Compile(strings.ToLower(kind))
	if err != nil {
		return nil, err
	}

	return kindMatcher{kindGlob: kindGlob}, nil
}

type kindMatcher struct {
	kindGlob glob.Glob
}

func (m kindMatcher) Matches(item resources.Resource) bool {
	// Initially, try to match the current kind.
	if m.kindGlob.Match(strings.ToLower(item.GetKind())) {
		return true
	}

	if resolved, found := resolver.LookupKind(item.GetKind()); found {
		// Otherwise try to match against any aliases for the current kind.
		for _, alias := range resolved.Aliases {
			if m.kindGlob.Match(alias) {
				return true
			}
		}
	}

	return false
}
