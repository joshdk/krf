// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package mflag provides utility functions for pairing matcher.Matcher
// functions with a pflag.FlagSet. The provided functions exist only to have
// a uniform and unified call-site for declaring a flag with an associated
// matcher.Matcher.
package mflag

import (
	"strings"

	"github.com/spf13/pflag"

	"github.com/joshdk/krf/matcher"
)

// FlagSet pairs a pflag.FlagSet with a number of matcher.Matcher instances
// (including negative matcher.Matcher instances). A series of individual flags
// can be defined, and then an aggregate matcher can be composed based on their
// provided values at runtime.
type FlagSet struct {
	flags         *pflag.FlagSet
	matcherFns    []func() ([]matcher.Matcher, error)
	notMatcherFns []func() ([]matcher.Matcher, error)
}

// NewMatcherFlags returns a new FlagSet bound to the provided pflag.FlagSet.
func NewMatcherFlags(flags *pflag.FlagSet) *FlagSet {
	return &FlagSet{
		flags: flags,
	}
}

// Matcher returns a composed matcher.Matcher derived from each defined flag
// and their runtime value(s).
func (m *FlagSet) Matcher() (matcher.Matcher, error) {
	chain := &matcher.AllMatcher{}

	// Loop over each not-matcher and add them directly to the chain. We add
	// them first because if any of the not-matcher fail, then the whole chain
	// fails.
	for _, fn := range m.notMatcherFns {
		matchers, err := fn()
		if err != nil {
			return nil, err
		}

		for _, mm := range matchers {
			chain.Append(matcher.NotMatcher(mm))
		}
	}

	// Loop over each matcher and add them to a sub-chain. We add
	// them second because we need at least one of them to pass in order to
	// pass the sub-chain.
	for _, fn := range m.matcherFns {
		matchers, err := fn()
		if err != nil {
			return nil, err
		}

		subchain := &matcher.AnyMatcher{}
		for _, mm := range matchers {
			subchain.Append(mm)
		}

		chain.Append(subchain)
	}

	return chain, nil
}

// BoolMatcher creates a named bool flag paired with the given matcher.Matcher
// constructor.
func (m *FlagSet) BoolMatcher(callback func() matcher.Matcher, name string, value bool, usage string) {
	result := m.flags.Bool(name, value, usage)

	fn := func() ([]matcher.Matcher, error) {
		if !*result {
			return nil, nil
		}

		return []matcher.Matcher{callback()}, nil
	}

	m.add(name, fn)
}

// StringSliceMatcher creates a named string slice flag paired with the given
// matcher.Matcher constructor.
func (m *FlagSet) StringSliceMatcher(callback func(string) matcher.Matcher, name string, value []string, usage string) {
	m.StringSliceErrorMatcher(func(s string) (matcher.Matcher, error) {
		return callback(s), nil
	}, name, value, usage)
}

// StringSliceErrorMatcher creates a named string slice flag paired with the
// given matcher.Matcher constructor (which can return an error).
func (m *FlagSet) StringSliceErrorMatcher(callback func(string) (matcher.Matcher, error), name string, value []string, usage string) {
	results := m.flags.StringSlice(name, value, usage)

	fn := func() ([]matcher.Matcher, error) {
		if len(*results) == 0 {
			return nil, nil
		}

		var matchers []matcher.Matcher

		for _, result := range *results {
			mm, err := callback(result)
			if err != nil {
				return nil, err
			}

			matchers = append(matchers, mm)
		}

		return matchers, nil
	}

	m.add(name, fn)
}

func (m *FlagSet) add(name string, fn func() ([]matcher.Matcher, error)) {
	if strings.HasPrefix(name, "not-") {
		m.notMatcherFns = append(m.notMatcherFns, fn)
	} else {
		m.matcherFns = append(m.matcherFns, fn)
	}
}
