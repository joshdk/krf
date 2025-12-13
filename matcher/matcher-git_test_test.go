// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher //nolint:testpackage

import (
	"github.com/go-git/go-git/v6"
)

// NewTestGitMatcher is a test-only helper function which bridges access to the
// internal newGitMatcher constructor.
func NewTestGitMatcher(currentDir string, gitDir string, gitStatus git.Status, pattern string) (Matcher, error) {
	return newGitMatcher(
		func() (string, error) {
			return currentDir, nil
		},
		func() (string, git.Status, error) {
			return gitDir, gitStatus, nil
		},
		pattern,
	)
}
