// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package printer

import (
	"io"

	"github.com/rodaine/table"

	"github.com/joshdk/krf/resources"
)

// Table prints each given resources.Resource as a row in a formatted table.
func Table(w io.Writer, items []resources.Resource) error {
	headers := []any{"Namespace", "API Version", "Kind", "Name"}

	// Check if any of the resources were decoded from a file opposed to from
	// e.g. stdin.
	for _, item := range items {
		if item.GetFilename() != "" {
			headers = append(headers, "Path")

			break
		}
	}

	tbl := table.New(headers...)
	tbl.WithHeaderSeparatorRow('─')
	tbl.WithWriter(w)

	for _, item := range items {
		tbl.AddRow(item.GetNamespace(), item.GetAPIVersion(), item.GetKind(), item.GetName(), item.GetFilename())
	}

	tbl.Print()

	return nil
}
