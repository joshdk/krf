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

// NewLabelMatcher matches resources.Resource instances that contain the given
// label, and optionally matching a value for that label.
func NewLabelMatcher(selector string) (Matcher, error) {
	key, value := splitSelector(selector)

	if key == "" {
		return nil, errors.New("empty label matcher")
	}

	keyGlob, err := glob.Compile(key)
	if err != nil {
		return nil, err
	}

	m := labelMatcher{keyGlob: keyGlob}

	// No label value was given to match against, so we'll only be checking for
	// the existence of the given label key.
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

type labelMatcher struct {
	keyGlob   glob.Glob
	valueGlob glob.Glob
}

func (m labelMatcher) Matches(item resources.Resource) bool {
	for key, value := range item.GetLabels() {
		if !m.keyGlob.Match(key) {
			continue
		}

		// Return the existence of a matching label key, since label value
		// matching was not requested.
		if m.valueGlob == nil {
			return true
		}

		if m.valueGlob.Match(value) {
			return true
		}
	}

	return false
}
