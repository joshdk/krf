// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package references provides functionality for searching through a resource
// for references to other resources. An example of a reference would be a Pod
// referencing a named ServiceAccount or a StatefulSet referencing a named
// PersistentVolumeClaim.
package references

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/joshdk/krf/references/path"
	"github.com/joshdk/krf/resolver"
)

// References searches the given unstructured.Unstructured for named resource
// references of the given kind.
func References(uu unstructured.Unstructured, kind string, callback func(s string) bool) bool {
	if resource, found := resolver.LookupKind(uu.GetKind()); found {
		for _, reference := range resource.References {
			// Skip reference check as this isn't the kind we're looking for.
			if kind != reference.Kind && kind != "" {
				continue
			}

			for _, spec := range reference.Paths {
				if path.Walk(uu.Object, spec, callback) {
					return true
				}
			}
		}
	}

	return false
}
