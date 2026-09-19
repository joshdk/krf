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

// Count prints the total number of resources.Resource which were matched.
func Count(w io.Writer, results []resources.Resource) error {
	fmt.Fprintln(w, len(results))

	return nil
}
