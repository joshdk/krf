// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package path

// split will split the given path into segments separated by '/' characters. A
// '\' can be used as an escape character.
func split(path string) []string {
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	segments := []string{""}
	escaped := false

	for _, char := range path {
		switch {
		case escaped:
			segments[len(segments)-1] += string(char)
			escaped = false

		case char == '\\':
			escaped = true

		case char == '/':
			segments = append(segments, "")

		default:
			segments[len(segments)-1] += string(char)
		}
	}

	return segments
}
