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
	parts := strings.SplitN(selector, "=", 2)
	switch len(parts) {
	case 2:
		return parts[0], parts[1]
	case 1:
		return parts[0], ""
	default:
		return "", ""
	}
}
