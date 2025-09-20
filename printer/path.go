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

// Path prints the file path of each given resources.Resource.
func Path(w io.Writer, results []resources.Resource) error {
	paths := make([]string, len(results))
	for i, item := range results {
		paths[i] = item.GetFilename()
	}

	sort.Strings(paths)

	// Print each unique resource path.
	var last string
	for _, path := range paths {
		if path == last {
			continue
		}

		last = path

		fmt.Fprintln(w, path)
	}

	return nil
}
