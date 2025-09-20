// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package printer

import (
	"fmt"
	"io"
	"sort"

	"github.com/joshdk/krf/resources"
)

// Name prints each given resources.Resource in "name" format similar to
// kubectl.
func Name(w io.Writer, results []resources.Resource) error {
	names := make([]string, len(results))
	for i, item := range results {
		name := fmt.Sprintf("%s/%s", item.GetKind(), item.GetName())
		names[i] = name
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
