// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v6/osfs"
	"github.com/go-git/go-git/v6"

	"github.com/joshdk/krf/resources"
)

func gitStatus() (string, git.Status, error) {
	repository, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return "", nil, err
	}

	worktree, err := repository.Worktree()
	if err != nil {
		return "", nil, err
	}

	filesystem, ok := worktree.Filesystem.(*osfs.BoundOS)
	if !ok {
		return "", nil, errors.New("could not determine git repo directory")
	}

	status, err := worktree.StatusWithOptions(git.StatusOptions{Strategy: git.Preload})
	if err != nil {
		return "", nil, err
	}

	return filesystem.Root(), status, nil
}

// NewGitMatcher matches resources.Resource instances based on the current git
// status of the underlying resource file.
func NewGitMatcher(pattern string) (Matcher, error) {
	return newGitMatcher(os.Getwd, gitStatus, pattern)
}

func newGitMatcher(getwd func() (string, error), statusFn func() (string, git.Status, error), pattern string) (Matcher, error) {
	var statusCode git.StatusCode

	switch pattern {
	case "added", "a", "A":
		statusCode = git.Added
	case "modified", "m", "M":
		statusCode = git.Modified
	case "unmodified", "u", "U", " ":
		statusCode = git.Unmodified
	case "untracked", "?":
		statusCode = git.Untracked
	default:
		return nil, fmt.Errorf("unsupported pattern: %q", pattern)
	}

	currentDir, err := getwd()
	if err != nil {
		return nil, err
	}

	gitDir, status, err := statusFn()
	if err != nil {
		return nil, err
	}

	relRoot, err := filepath.Rel(gitDir, currentDir)
	if err != nil {
		return nil, err
	}

	return gitMatcher{
		relRoot:    relRoot,
		status:     status,
		statusCode: statusCode,
	}, nil
}

type gitMatcher struct {
	status     git.Status
	statusCode git.StatusCode
	relRoot    string
}

func (m gitMatcher) Matches(item resources.Resource) bool {
	filename := item.GetFilename()
	if filename == "" {
		return false
	}

	fullFile := filepath.Join(m.relRoot, filename)

	fileStatus, ok := m.status[fullFile]
	if !ok {
		fileStatus = &git.FileStatus{Worktree: git.Untracked, Staging: git.Untracked}
	}

	switch m.statusCode { //nolint:exhaustive
	case git.Added:
		return fileStatus.Staging == git.Modified
	case git.Modified:
		return fileStatus.Staging == git.Modified || fileStatus.Worktree == git.Modified
	case git.Unmodified:
		return fileStatus.Staging == git.Unmodified && fileStatus.Worktree == git.Unmodified
	case git.Untracked:
		return fileStatus.Staging == git.Untracked && fileStatus.Worktree == git.Untracked
	default:
		return false
	}
}
