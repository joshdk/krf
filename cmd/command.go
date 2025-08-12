// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package cmd provides the command line handler for krf.
package cmd

import (
	"github.com/spf13/cobra"
)

// Command returns a complete command line handler for krf.
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "krf [directory|file|-]",
		Long: "krf - kubernetes resources filter",

		SilenceUsage:  true,
		SilenceErrors: true,

		Args: cobra.MaximumNArgs(1),

		RunE: func(*cobra.Command, []string) error {
			return nil
		},
	}

	return cmd
}
