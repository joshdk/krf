// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher

import (
	"strings"

	"github.com/gobwas/glob"
)

// asGlob returns a glob.Glob from the given pattern. Additionally, if the
// input pattern is prefixed/suffixed with a dash (`-`) character, then an
// additional glob wildcard is prepended/appended respectfully. This is a
// convenience as it is easier to e.g. type `--name my-pod-` than
// `--name 'my-pod-*'` in the terminal.
func asGlob(pattern string) (glob.Glob, error) {
	if strings.HasPrefix(pattern, "-") {
		pattern = "*" + pattern
	}

	if strings.HasSuffix(pattern, "-") {
		pattern += "*"
	}

	return glob.Compile(pattern)
}

// splitSelector splits the given string into two pieces around an equal sign.
func splitSelector(selector string) (string, string) {
	// Index for the last '=' and ']' characters. This is done as there are
	// some selectors like '.spec.containers.[name=main].securityContext' and
	// '.spec.containers.[name=main].image=docker.io/...' where blindly
	// splitting on the first (or last) '=' character would be incorrect.
	var (
		indexEquals       = strings.LastIndex(selector, "=")
		indexRightBracket = strings.LastIndex(selector, "]")
	)

	switch {
	case indexRightBracket < indexEquals:
		// The last '=' character exists after the last ']' character. Split on
		// that index.
		return selector[:indexEquals], selector[indexEquals+1:]
	case indexEquals == -1:
		// There is no '=' character. Return the input selector verbatim.
		return selector, ""
	default:
		// There is no '=' character *after* a ']' character. Return the input
		// selector verbatim.
		return selector, ""
	}
}
