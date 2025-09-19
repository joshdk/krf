// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package references

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/joshdk/krf/references/path"
	"github.com/joshdk/krf/resolver"
)

// All iterates over all named resource references in the given
// unstructured.Unstructured.
func All(uu unstructured.Unstructured, callback func(kind, name string)) {
	resource, found := resolver.LookupKind(uu.GetKind())
	if !found {
		return
	}

	for _, reference := range resource.References {
		for _, spec := range reference.Paths {
			path.Walk(uu.Object, spec, func(s string) bool {
				callback(reference.Kind, s)

				return false
			})
		}
	}
}
