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
	"os"

	"golang.org/x/term"

	"github.com/joshdk/krf/resources"
)

// ByName returns a printer function from the given name. If no name is given
// and the program output is being redirected or piped to a consumer process,
// hen default to the YAML printer. Additionally, if no name is given and the
// program output is being sent directly to the terminal, then instead default
// to the Table printer.
func ByName(name string) (func(io.Writer, []resources.Resource) error, error) {
	switch name {
	case "":
		// Is program output being redirected or piped to a consumer process?
		if !term.IsTerminal(int(os.Stdout.Fd())) {
			return YAML, nil
		}

		// Default for when output is directly to a terminal.
		return Table, nil

	case "name":
		return Name, nil

	case "references":
		return References, nil

	case "table":
		return Table, nil

	case "yaml":
		return YAML, nil

	default:
		return nil, fmt.Errorf("unknown printer name: %s", name)
	}
}
