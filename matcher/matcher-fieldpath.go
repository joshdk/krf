// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher

import (
	"fmt"

	"github.com/gobwas/glob"
	"sigs.k8s.io/kustomize/kyaml/yaml"

	"github.com/joshdk/krf/resources"
)

// NewFieldPathMatcher matches resources.Resource instances that contain the given
// Kustomize-style fieldpath, and optionally matching a target value at that
// path.
func NewFieldPathMatcher(selector string) (Matcher, error) {
	path, value := splitSelector(selector)

	m := fieldPathMatcher{fieldPath: path}

	// No target value was given to match against, so we'll only be checking
	// for the existence of the given fieldpath...path.
	if value == "" {
		return m, nil
	}

	valueGlob, err := glob.Compile(value)
	if err != nil {
		return nil, err
	}

	m.valueGlob = valueGlob

	return m, nil
}

type fieldPathMatcher struct {
	fieldPath string
	valueGlob glob.Glob
}

func (m fieldPathMatcher) Matches(item resources.Resource) bool {
	node, err := yaml.FromMap(item.Object)
	if err != nil {
		panic(err)
	}

	value, err := node.GetFieldValue(m.fieldPath)
	if err != nil {
		return false
	}

	// Return the existence of a matching fieldpath, since value matching was
	// not requested.
	if m.valueGlob == nil {
		return true
	}

	switch value.(type) {
	case string, int, bool, float64:
		// Fieldpath evaluated to a simple type. Match the value glob against
		// the stringified value.
		if m.valueGlob.Match(fmt.Sprint(value)) {
			return true
		}
	}

	return false
}
