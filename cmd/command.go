// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package cmd provides the command line handler for krf.
package cmd

import (
	"fmt"

	"github.com/joshdk/buildversion"
	"github.com/spf13/cobra"

	"github.com/joshdk/krf/resources"
)

// Command returns a complete command line handler for krf.
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "krf [directory|file|-]",
		Long: "krf - kubernetes resources filter",

		SilenceUsage:  true,
		SilenceErrors: true,

		Args: cobra.MaximumNArgs(1),

		RunE: func(_ *cobra.Command, args []string) error {
			return cmdFunc(args)
		},

		Version: "-",
	}

	// Add a custom usage footer template.
	cmd.SetUsageTemplate(cmd.UsageTemplate() + "\n" + buildversion.Template(usageTemplate))

	// Set a custom version template.
	cmd.SetVersionTemplate(buildversion.Template(versionTemplate))

	return cmd
}

func cmdFunc(args []string) error {
	var results []resources.Resource

	var source string
	if len(args) > 0 {
		source = args[0]
	}

	err := resources.Decode(source, func(item resources.Resource) {
		results = append(results, item)
	})
	if err != nil {
		return err
	}

	for _, result := range results {
		fmt.Printf("%s/%s\n", result.GetKind(), result.GetName())
	}

	return nil
}
