// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gobwas/glob"
	"k8s.io/client-go/util/jsonpath"

	"github.com/joshdk/krf/resources"
)

// NewJsonpathMatcher matches resources.Resource instances that contain the given
// jsonpath, and optionally matching a value at that path.
func NewJsonpathMatcher(selector string) (Matcher, error) {
	key, value := splitSelector(selector)

	keyJsonpath, err := newJsonpath(key)
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
	keyJsonpath *jsonpath.JSONPath
	valueGlob   glob.Glob
}

func (m jsonpathMatcher) Matches(item resources.Resource) bool {
	results, err := m.keyJsonpath.FindResults(item.Object)
	if err != nil {
		return false
	}

	if len(results) == 0 || len(results[0]) == 0 {
		return false
	}

	// Return the existence of a matching jsonpath, since value
	// matching was not requested.
	if m.valueGlob == nil {
		return true
	}

	// Loop over the (potentially) multiple matches and compare the values at
	// each match.
	for _, value := range results[0] {
		switch value.Interface().(type) {
		case string, int, bool, float64:
			// Jsonpath evaluated to a simple type. Match the value glob against
			// the stringified value.
			if m.valueGlob.Match(fmt.Sprint(value)) {
				return true
			}

		case nil:
			// Jsonpath evaluated to something nil. Match the value glob against
			// the literal string "null".
			if m.valueGlob.Match("null") {
				return true
			}
		}
	}

	return false
}

func newJsonpath(spec string) (*jsonpath.JSONPath, error) {
	// Do not require that the user include the surrounding '{...}' on the
	// jsonpath.
	if !strings.HasPrefix(spec, "{") {
		spec = "{" + spec
	}

	if !strings.HasSuffix(spec, "}") {
		spec += "}"
	}

	if spec == "{}" {
		return nil, errors.New("jsonpath was empty")
	}

	// Prevent the usage of multiple jsonpath expressions like "{.foo} {.bar}".
	if jp, err := jsonpath.Parse("matcher", spec); err != nil {
		return nil, err
	} else if len(jp.Root.Nodes) != 1 {
		return nil, errors.New("jsonpath contained multiple expressions")
	}

	path := jsonpath.New("matcher").AllowMissingKeys(false)

	return path, path.Parse(spec)
}
