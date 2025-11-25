// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher

import (
	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/joshdk/krf/resources"
)

// NewDiffMatcher matches resources.Resource instances that differ from
// companion resources decoded from the given filename. An example usage would
// be to save the output of `kustomize build` to a file, make refactors to the
// kustomize code, then rerun `kustomize build` to check if resources were
// changed.
func NewDiffMatcher(filename string) (Matcher, error) {
	var originals []unstructured.Unstructured

	// Decode resources from the given file to be used later for matching.
	if err := resources.Decode(filename, func(item resources.Resource) {
		originals = append(originals, item.Unstructured)
	}); err != nil {
		return nil, err
	}

	return diffMatcher{originals: originals}, nil
}

type diffMatcher struct {
	originals []unstructured.Unstructured
}

func (m diffMatcher) Matches(item resources.Resource) bool {
	for _, original := range m.originals {
		switch {
		// Check if the resource under investigation is the counterpart of any
		// original resources.
		case item.GetAPIVersion() != original.GetAPIVersion():
			continue
		case item.GetKind() != original.GetKind():
			continue
		case item.GetName() != original.GetName():
			continue
		case item.GetNamespace() != original.GetNamespace():
			continue
		}

		// Diff the resource under investigation against its original
		// counterpart.
		return !cmp.Equal(item.Object, original.Object)
	}

	// The resource under investigation did not have an original counterpart,
	// and has therefore differed.
	return true
}
