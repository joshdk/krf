// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package cmd provides the command line handler for krf.
package cmd

import (
	"os"
	"strings"

	"github.com/joshdk/buildversion"
	"github.com/spf13/cobra"

	"github.com/joshdk/krf/cmd/mflag"
	"github.com/joshdk/krf/matcher"
	"github.com/joshdk/krf/printer"
	"github.com/joshdk/krf/resources"
)

// Command returns a complete command line handler for krf.
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "krf [directory|file|-]",
		Long:    "krf - kubernetes resource filter",
		Version: "-",

		SilenceUsage:  true,
		SilenceErrors: true,

		Args: cobra.MaximumNArgs(1),
	}

	// Set a custom list of examples.
	cmd.Example = strings.TrimRight(exampleText, "\n")

	// Add a custom usage footer template.
	cmd.SetUsageTemplate(cmd.UsageTemplate() + "\n" + buildversion.Template(usageTemplate))

	// Set a custom version template.
	cmd.SetVersionTemplate(buildversion.Template(versionTemplate))

	mf := mflag.NewMatcherFlags(cmd.Flags())

	// Define --name flag.
	mf.StringSliceErrorMatcher(matcher.NewNameMatcher,
		"name",
		nil,
		"include resources by name")

	// Define --not-name flag.
	mf.StringSliceErrorMatcher(matcher.NewNameMatcher,
		"not-name",
		nil,
		"exclude resources by name")

	// Define --namespace flag.
	mf.StringSliceErrorMatcher(matcher.NewNamespaceMatcher,
		"namespace",
		nil,
		"include resources by namespace")

	// Define --not-namespace flag.
	mf.StringSliceErrorMatcher(matcher.NewNamespaceMatcher,
		"not-namespace",
		nil,
		"exclude resources by namespace")

	// Define --format flag.
	output := cmd.Flags().StringP(
		"output",
		"o",
		"",
		"output format (name,table,yaml)")

	cmd.RunE = func(_ *cobra.Command, args []string) error {
		return cmdFunc(mf, *output, args)
	}

	return cmd
}

func cmdFunc(mf *mflag.FlagSet, output string, args []string) error {
	allMatchers, err := mf.Matcher()
	if err != nil {
		return err
	}

	printerFn, err := printer.ByName(output)
	if err != nil {
		return err
	}

	var source string
	if len(args) > 0 {
		source = args[0]
	}

	var results []resources.Resource

	err = resources.Decode(source, func(item resources.Resource) {
		if allMatchers.Matches(item) {
			results = append(results, item)
		}
	})
	if err != nil {
		return err
	}

	return printerFn(os.Stdout, results)
}
