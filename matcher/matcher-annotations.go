// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher

import (
	"errors"

	"github.com/gobwas/glob"

	"github.com/joshdk/krf/resources"
)

// NewAnnotationMatcher matches resources.Resource instances that contain the
// given annotation, and optionally matching a value for that annotation.
func NewAnnotationMatcher(selector string) (Matcher, error) {
	key, value := splitSelector(selector)

	if key == "" {
		return nil, errors.New("empty annotation matcher")
	}

	keyGlob, err := glob.Compile(key)
	if err != nil {
		return nil, err
	}

	m := annotationMatcher{keyGlob: keyGlob}

	// No annotation value was given to match against, so we'll only be
	// checking for the existence of the given annotation key.
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

type annotationMatcher struct {
	keyGlob   glob.Glob
	valueGlob glob.Glob
}

func (m annotationMatcher) Matches(item resources.Resource) bool {
	for key, value := range item.GetAnnotations() {
		if !m.keyGlob.Match(key) {
			continue
		}

		// Return the existence of a matching annotation key, since annotation
		// value matching was not requested.
		if m.valueGlob == nil {
			return true
		}

		if m.valueGlob.Match(value) {
			return true
		}
	}

	return false
}
