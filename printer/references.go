// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package printer

import (
	"fmt"
	"io"
	"sort"

	"github.com/joshdk/krf/references"
	"github.com/joshdk/krf/resources"
)

// References prints the names of all resources referenced by each given
// resources.Resource.
func References(w io.Writer, results []resources.Resource) error {
	// Build a list of all resources referenced by each of the given resources.
	var names []string

	for _, item := range results {
		references.All(item.Unstructured, func(referenceKind, referenceName string) {
			name := fmt.Sprintf("%s/%s", referenceKind, referenceName)
			names = append(names, name)
		})
	}

	sort.Strings(names)

	// Print each unique resource name.
	var last string
	for _, name := range names {
		if name == last {
			continue
		}

		last = name

		fmt.Fprintln(w, name)
	}

	return nil
}
