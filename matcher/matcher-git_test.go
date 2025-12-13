// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher_test

import (
	"testing"

	"github.com/go-git/go-git/v6"

	"github.com/joshdk/krf/matcher"
)

func TestGitMatcher(t *testing.T) { //nolint:funlen
	t.Parallel()

	gitStatusRoot := git.Status{
		"testdata/service.yaml": &git.FileStatus{
			Staging:  git.Unmodified,
			Worktree: git.Modified,
		},
		"testdata/subdir/deployment.yaml": &git.FileStatus{
			Staging:  git.Unmodified,
			Worktree: git.Modified,
		},
	}

	gitStatusSubdir := git.Status{
		"subdir/testdata/service.yaml": &git.FileStatus{
			Staging:  git.Modified,
			Worktree: git.Unmodified,
		},
		"subdir/testdata/subdir/deployment.yaml": &git.FileStatus{
			Staging:  git.Unmodified,
			Worktree: git.Modified,
		},
		"subdir/testdata/subdir/subsubdir/configmap.yaml": &git.FileStatus{
			Staging:  git.Unmodified,
			Worktree: git.Unmodified,
		},
	}

	testMatcher(t, []spec{
		{
			title:   "at git repo root",
			matcher: must(matcher.NewTestGitMatcher("/root/repo", "/root/repo", gitStatusRoot, "m")),
			matches: []string{
				"Deployment/nginx-deployment",
				"Service/my-service",
			},
		},
		{
			title:   "at git repo subdir",
			matcher: must(matcher.NewTestGitMatcher("/root/repo/subdir", "/root/repo", gitStatusSubdir, "M")),
			matches: []string{
				"Deployment/nginx-deployment",
				"Service/my-service",
			},
		},
		{
			title:   "file added",
			matcher: must(matcher.NewTestGitMatcher("/root/repo/subdir", "/root/repo", gitStatusSubdir, "added")),
			matches: []string{
				"Service/my-service",
			},
		},
		{
			title:   "file unmodified",
			matcher: must(matcher.NewTestGitMatcher("/root/repo/subdir", "/root/repo", gitStatusSubdir, "unmodified")),
			matches: []string{
				"ConfigMap/my-configmap",
			},
		},
		{
			title:   "file untracked",
			matcher: must(matcher.NewTestGitMatcher("/root/repo/subdir", "/root/repo", gitStatusSubdir, "?")),
			matches: []string{
				"ClusterRoleBinding/read-secrets-global",
				"Pod/test-pod",
			},
		},
	})
}
