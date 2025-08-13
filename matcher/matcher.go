// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package matcher provides functionality for building and composing functions
// for matching against various properties of resources.Resource objects based
// on user input.
package matcher

import "github.com/joshdk/krf/resources"

// Matcher represents the logic for matching against the properties of a
// resources.Resource based on some inputs.
type Matcher interface {
	// Matches returns true if the given resources.Resource object matches based
	// on the logic and inputs for a concrete matcher.
	Matches(item resources.Resource) bool
}
