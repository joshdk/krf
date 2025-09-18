// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package matcher

import (
	"strings"

	"sigs.k8s.io/yaml"

	"github.com/joshdk/krf/resources"
)

// NewContainsMatcher matches resources.Resource instances that contain the
// given substring after being re-marshalled back into YAML. This also applies
// to resources that were originally decoded from JSON files. Additionally,
// these re-marshalled YAML documents no longer contain any of the original
// comments or other formatting.
func NewContainsMatcher(substring string) Matcher {
	return containsMatcher{substring: substring}
}

type containsMatcher struct {
	substring string
}

func (m containsMatcher) Matches(item resources.Resource) bool {
	body, err := yaml.Marshal(item.Object)
	if err != nil {
		return false
	}

	return strings.Contains(string(body), m.substring)
}
