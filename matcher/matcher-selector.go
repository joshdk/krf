// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher

import (
	"k8s.io/apimachinery/pkg/labels"

	"github.com/joshdk/krf/resources"
)

// NewSelectorMatcher matches resources.Resource instances that match the given
// Kubernetes label selector.
func NewSelectorMatcher(selector string) (Matcher, error) {
	sel, err := labels.Parse(selector)
	if err != nil {
		return nil, err
	}

	return selectorMatcher{
		selector: sel,
	}, nil
}

type selectorMatcher struct {
	selector labels.Selector
}

func (m selectorMatcher) Matches(item resources.Resource) bool {
	return m.selector.Matches(labels.Set(item.GetLabels()))
}
