// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher

import (
	"strings"

	"github.com/gobwas/glob"

	"github.com/joshdk/krf/resources"
)

// NewArgoCDMatcher matches resources.Resource instances based on the ArgoCD
// Application that the resource is being tracked by.
func NewArgoCDMatcher(selector string) (Matcher, error) {
	appGlob, err := glob.Compile(selector)
	if err != nil {
		return nil, err
	}

	m := argocdMatcher{appGlob: appGlob}

	return m, nil
}

type argocdMatcher struct {
	appGlob glob.Glob
}

func (m argocdMatcher) Matches(item resources.Resource) bool {
	value, found := item.GetAnnotations()["argocd.argoproj.io/tracking-id"]
	if !found {
		return false
	}

	// The format of the tracking annotation looks like this:
	// <application name>:<resource group>/<resource kind>:<resource namespace>/<resource name>
	// We are interested in the application name, so we look for the location
	// of the first ':' character and slice the string up to that point.
	index := strings.IndexByte(value, ':')
	if index == -1 {
		return false
	}

	app := value[:index]

	return m.appGlob.Match(app)
}
