// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package path provides functionality for walking through an object using a
// path in search of specific values.
package path

// Walk traverses through the given object to find strings located by the given
// path. Each located string is passed to the given callback function. If the
// callback function ever returns true, then traversal stops immediately.
func Walk(object any, path string, callback func(string) bool) bool {
	return walk(object, split(path), callback)
}

func walk(object any, paths []string, callback func(string) bool) bool {
	switch current := object.(type) {
	case map[string]any:
		// Since there are no more path segments, it isn't possible to Walk to
		// any other strings.
		if len(paths) == 0 {
			return false
		}

		// Descend into the map using the next path segment.
		if next, found := current[paths[0]]; found {
			return walk(next, paths[1:], callback)
		}

	case []any:
		// Descend into each list entry.
		for _, next := range current {
			if walk(next, paths, callback) {
				return true
			}
		}

	case string:
		// There are no more path segments, so we have found our target string.
		if len(paths) == 0 {
			return callback(current)
		}
	}

	return false
}
