// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package printer

import (
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"

	"github.com/rodaine/table"

	"github.com/joshdk/krf/resources"
)

// Selector prints each given resources.Resource with the resource name and a
// matching label selector.
func Selector(w io.Writer, items []resources.Resource) error {
	tbl := table.New("Name", "Selector")
	tbl.WithPrintHeaders(false)
	tbl.WithWriter(w)

	for _, item := range items {
		var (
			labels    = item.GetLabels()
			selectors []string
		)

		// Sort all labels by key name.
		for _, key := range slices.Sorted(maps.Keys(labels)) {
			selectors = append(selectors, fmt.Sprintf("%s=%s", key, labels[key]))
		}

		tbl.AddRow(
			fmt.Sprintf("%s/%s", item.GetKind(), item.GetName()),
			strings.Join(selectors, ","),
		)
	}

	tbl.Print()

	return nil
}
