// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher

import (
	"github.com/joshdk/krf/resources"
)

// NotMatcher wraps the given Matcher instances and inverts its value.
func NotMatcher(matcher Matcher) Matcher {
	return notMatcher{matcher}
}

// NotMatcher wraps a single Matcher instance and inverts its value.
type notMatcher struct {
	matcher Matcher
}

func (m notMatcher) Matches(item resources.Resource) bool {
	return !m.matcher.Matches(item)
}

// AllMatcher wraps a sequence of Matcher instances and returns true if each of
// the wrapped Matcher instances returns true.
type AllMatcher struct {
	matchers []Matcher
}

// Append adds the given Matcher to the sequence of wrapped Matcher instances.
func (m *AllMatcher) Append(matcher Matcher) {
	if matcher == nil {
		return
	}

	m.matchers = append(m.matchers, matcher)
}

// Matches returns true if each of the wrapped Matcher instances returns true.
func (m *AllMatcher) Matches(item resources.Resource) bool {
	for _, matcher := range m.matchers {
		if !matcher.Matches(item) {
			return false
		}
	}

	return true
}

// AnyMatcher wraps a sequence of Matcher instances and returns true if any of
// the wrapped Matcher instances returns true. An empty AnyMatcher will also
// return true.
type AnyMatcher struct {
	matchers []Matcher
}

// Append adds the given Matcher to the sequence of wrapped Matcher instances.
func (m *AnyMatcher) Append(matcher Matcher) {
	if matcher == nil {
		return
	}

	m.matchers = append(m.matchers, matcher)
}

// Matches returns true if any of the wrapped Matcher instances returns true.
// An empty AnyMatcher will also return true.
func (m *AnyMatcher) Matches(item resources.Resource) bool {
	if len(m.matchers) == 0 {
		return true
	}

	for _, matcher := range m.matchers {
		if matcher.Matches(item) {
			return true
		}
	}

	return false
}
