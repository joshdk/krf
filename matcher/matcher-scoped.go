// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher

import (
	"github.com/joshdk/krf/resolver"
	"github.com/joshdk/krf/resources"
)

// NewClusterScopedMatcher matches resources.Resource instances that are known
// (via the resolver) to be cluster-scoped. As a caveat, this matcher will not
// match unresolved resources even if they are cluster-scoped in reality.
func NewClusterScopedMatcher() Matcher {
	return clusterScopedMatcher{}
}

type clusterScopedMatcher struct{}

func (m clusterScopedMatcher) Matches(item resources.Resource) bool {
	if item.GetNamespace() != "" {
		// This resource explicitly set a namespace, and so cannot be global.
		return false
	}

	if resource, found := resolver.LookupKind(item.GetKind()); found {
		// This resource kind was registered, so use the authoritative namespaced value.
		return !resource.Namespaced
	}

	// Reject the resource as there was no concrete conclusion.
	return false
}

// NewNamespaceScopedMatcher matches resources.Resource instances that are
// known (via the resolver) to be namespace-scoped. As a caveat, this matcher
// will not match unresolved resources even if they are namespace-scoped in
// reality.
func NewNamespaceScopedMatcher() Matcher {
	return namespaceScopedMatcher{}
}

type namespaceScopedMatcher struct{}

func (m namespaceScopedMatcher) Matches(item resources.Resource) bool {
	if item.GetNamespace() != "" {
		// This resource explicitly set a namespace, and so must be namespaced.
		return true
	}

	if resource, found := resolver.LookupKind(item.GetKind()); found {
		// This resource kind was registered, so use the authoritative namespaced value.
		return resource.Namespaced
	}

	// Reject the resource as there was no concrete conclusion.
	return false
}
