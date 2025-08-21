// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package cmd provides the command line handler for krf.
package cmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/joshdk/buildversion"
	"github.com/spf13/cobra"

	"github.com/joshdk/krf/cmd/mflag"
	"github.com/joshdk/krf/config"
	"github.com/joshdk/krf/matcher"
	"github.com/joshdk/krf/printer"
	"github.com/joshdk/krf/resolver"
	"github.com/joshdk/krf/resources"
)

// Command returns a complete command line handler for krf.
func Command() *cobra.Command { //nolint:funlen
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

	// Define --cluster-scoped flag.
	mf.BoolMatcher(matcher.NewClusterScopedMatcher,
		"cluster-scoped",
		false,
		"include resources that are cluster-scoped")

	// Define --kind flag.
	mf.StringSliceErrorMatcher(matcher.NewKindMatcher,
		"kind",
		nil,
		"include resources by kind")

	// Define --not-kind flag.
	mf.StringSliceErrorMatcher(matcher.NewKindMatcher,
		"not-kind",
		nil,
		"exclude resources by kind")

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

	// Define --namespace-scoped flag.
	mf.BoolMatcher(matcher.NewNamespaceScopedMatcher,
		"namespace-scoped",
		false,
		"include resources that are namespace-scoped")

	// Define --references flag.
	mf.StringSliceErrorMatcher(matcher.NewReferenceMatcher,
		"references",
		nil,
		"include resources that reference resource")

	// Define --not-references flag.
	mf.StringSliceErrorMatcher(matcher.NewReferenceMatcher,
		"not-references",
		nil,
		"exclude resources that reference resource")

	// Define --config flag.
	cfgfile := cmd.Flags().String(
		"config",
		"~/.config/krf/configuration.yaml",
		"path to config file")

	// Define --output flag.
	output := cmd.Flags().StringP(
		"output",
		"o",
		"",
		"output format (name,table,yaml)")

	var state struct {
		allMatchers matcher.Matcher
		printerFn   func(io.Writer, []resources.Resource) error
		source      string
	}

	cmd.PreRunE = func(_ *cobra.Command, args []string) error {
		configFilename := *cfgfile
		if strings.HasPrefix(configFilename, "~/") {
			configFilename = filepath.Join(os.Getenv("HOME"), configFilename[2:])
		}

		cfg, err := config.InitAndLoad(configFilename)
		if err != nil {
			return err
		}

		resolver.Init(cfg.Resources)

		state.printerFn, err = printer.ByName(*output)
		if err != nil {
			return err
		}

		state.allMatchers, err = mf.Matcher()
		if err != nil {
			return err
		}

		if len(args) > 0 {
			state.source = args[0]
		}

		return nil
	}

	cmd.RunE = func(*cobra.Command, []string) error {
		var results []resources.Resource

		err := resources.Decode(state.source, func(item resources.Resource) {
			if state.allMatchers.Matches(item) {
				results = append(results, item)
			}
		})
		if err != nil {
			return err
		}

		return state.printerFn(os.Stdout, results)
	}

	return cmd
}
