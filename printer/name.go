// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package printer provides functions for rendering resources.Resource
// collections into various formats for outputting to the terminal or command
// pipeline.
package printer

import (
	"fmt"
	"io"

	"github.com/joshdk/krf/resources"
)

// Name prints each given resources.Resource in "name" format similar to
// kubectl.
func Name(w io.Writer, results []resources.Resource) error {
	for _, item := range results {
		if _, err := fmt.Fprintf(w, "%s/%s\n", item.GetKind(), item.GetName()); err != nil {
			return err
		}
	}

	return nil
}
