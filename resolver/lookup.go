// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package resolver provides functionality for looking up resource metadata
// definitions using a resource kind or alias.
package resolver

import (
	"strings"
)

// Resource is a metadata definition for a single kind.
type Resource struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`

	Aliases []string `yaml:"aliases"`
}

// resources is a shared collection of Resource metadata definitions utilized
// by the lookup functions. Must be initialized using Init prior to calling the
// lookup functions.
var resources []Resource

// Init initializes the collection of known Resource metadata definitions.
func Init(rs []Resource) {
	resources = rs
}

// LookupAlias returns Resource metadata definitions that match the given
// alias. For example, this could resolve the string "po" to a Pod metadata
// definition.
//
// This function returns a Resource list since there might be multiple kinds
// with overlapping aliases.
func LookupAlias(name string) []Resource {
	name = strings.ToLower(name)

	var results []Resource

	for _, resource := range resources {
		for _, alias := range resource.Aliases {
			if alias == name {
				results = append(results, resource)
			}
		}
	}

	return results
}

// LookupKind returns the Resource metadata definition using the given kind.
// Meant to be used with a kind taken directly from an actual manifest.
func LookupKind(kind string) (Resource, bool) {
	for _, resource := range resources {
		if kind == resource.Kind {
			return resource, true
		}
	}

	return Resource{}, false
}
