// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package printer

import (
	"fmt"
	"io"

	"github.com/joshdk/krf/resources"
)

// JSON prints each given resources.Resource as a single-line json document.
// Intended to be used when piping to a script that could e.g. `read` each
// resource in a loop.
func JSON(w io.Writer, results []resources.Resource) error {
	for _, item := range results {
		data, err := item.MarshalJSON()
		if err != nil {
			return err
		}

		fmt.Fprint(w, string(data))
	}

	return nil
}
