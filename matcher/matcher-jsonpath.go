// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/PaesslerAG/gval"
	"github.com/PaesslerAG/jsonpath"
	"github.com/gobwas/glob"

	"github.com/joshdk/krf/resources"
)

// NewJsonpathMatcher matches resources.Resource instances that contain the given
// jsonpath, and optionally matching a value at that path.
func NewJsonpathMatcher(selector string) (Matcher, error) {
	key, value := splitSelector(selector)

	if key == "" {
		return nil, errors.New("empty jsonpath matcher")
	}

	// Do not require that the user include a '$' prefix on the jsonpath.
	if !strings.HasPrefix(key, "$") {
		key = "$" + key
	}

	keyJsonpath, err := jsonpath.New(key)
	if err != nil {
		return nil, err
	}

	m := jsonpathMatcher{keyJsonpath: keyJsonpath}

	// No jsonpath value was given to match against, so we'll only be
	// checking for the existence of the given jsonpath...path.
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

type jsonpathMatcher struct {
	keyJsonpath gval.Evaluable
	valueGlob   glob.Glob
}

func (m jsonpathMatcher) Matches(item resources.Resource) bool {
	value, err := m.keyJsonpath(context.Background(), item.Object)
	if err != nil {
		return false
	}

	// Return the existence of a matching jsonpath, since value
	// matching was not requested.
	if m.valueGlob == nil {
		return value != nil
	}

	switch value.(type) {
	case string, int, bool, float64:
		// Jsonpath evaluated to a simple type. Match the value glob against
		// the stringified value.
		return m.valueGlob.Match(fmt.Sprint(value))

	case nil:
		// Jsonpath evaluated to something nil. Match the value glob against
		// the literal string "null".
		return m.valueGlob.Match("null")

	default:
		// Jsonpath evaluated to something complex (e.g. a slice or struct)
		// which isn't possible to match against.
		return false
	}
}
