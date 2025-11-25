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
	"slices"
	"strings"

	"github.com/joshdk/buildversion"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

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

	// Define --annotation flag.
	mf.StringSliceErrorMatcher(matcher.NewAnnotationMatcher,
		"annotation",
		nil,
		"include resources by annotation")

	// Define --not-annotation flag.
	mf.StringSliceErrorMatcher(matcher.NewAnnotationMatcher,
		"not-annotation",
		nil,
		"exclude resources by annotation")

	// Define --apiversion flag.
	mf.StringSliceErrorMatcher(matcher.NewAPIVersionMatcher,
		"apiversion",
		nil,
		"include resources by api version")

	// Define --not-apiversion flag.
	mf.StringSliceErrorMatcher(matcher.NewAPIVersionMatcher,
		"not-apiversion",
		nil,
		"exclude resources by api version")

	// Define --cluster-scoped flag.
	mf.BoolMatcher(matcher.NewClusterScopedMatcher,
		"cluster-scoped",
		false,
		"include resources that are cluster-scoped")

	// Define --contains flag.
	mf.StringSliceMatcher(matcher.NewContainsMatcher,
		"contains",
		nil,
		"include resources by substring contents")

	// Define --not-contains flag.
	mf.StringSliceMatcher(matcher.NewContainsMatcher,
		"not-contains",
		nil,
		"exclude resources by substring contents")

	// Define --diff flag.
	mf.StringErrorMatcher(matcher.NewDiffMatcher,
		"diff",
		"",
		"include resources which differ from those in a file")

	// Define --not-diff flag.
	mf.StringErrorMatcher(matcher.NewDiffMatcher,
		"not-diff",
		"",
		"exclude resources which differ from those in a file")

	// Define --exec flag.
	mf.StringErrorMatcher(matcher.NewExecMatcher,
		"exec",
		"",
		"include resources by executing a script")

	// Define --not-exec flag.
	mf.StringErrorMatcher(matcher.NewExecMatcher,
		"not-exec",
		"",
		"exclude resources by executing a script")

	// Define --fieldpath flag.
	mf.StringSliceErrorMatcher(matcher.NewFieldPathMatcher,
		"fieldpath",
		nil,
		"include resources by kustomize fieldpath")

	// Define --not-fieldpath flag.
	mf.StringSliceErrorMatcher(matcher.NewFieldPathMatcher,
		"not-fieldpath",
		nil,
		"exclude resources by kustomize fieldpath")

	// Define --jsonpath flag.
	mf.StringSliceErrorMatcher(matcher.NewJsonpathMatcher,
		"jsonpath",
		nil,
		"include resources by jsonpath")

	// Define --not-jsonpath flag.
	mf.StringSliceErrorMatcher(matcher.NewJsonpathMatcher,
		"not-jsonpath",
		nil,
		"exclude resources by jsonpath")

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

	// Define --label flag.
	mf.StringSliceErrorMatcher(matcher.NewLabelMatcher,
		"label",
		nil,
		"include resources by label")

	// Define --not-label flag.
	mf.StringSliceErrorMatcher(matcher.NewLabelMatcher,
		"not-label",
		nil,
		"exclude resources by label")

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

	// Define --patch flag.
	mf.BoolMatcher(matcher.NewPatchMatcher,
		"patch",
		false,
		"include resources from patch files")

	// Define --not-patch flag.
	mf.BoolMatcher(matcher.NewPatchMatcher,
		"not-patch",
		false,
		"exclude resources from patch files")

	// Define --path flag.
	mf.StringSliceMatcher(matcher.NewPathMatcher,
		"path",
		nil,
		"include resources by file path")

	// Define --not-path flag.
	mf.StringSliceMatcher(matcher.NewPathMatcher,
		"not-path",
		nil,
		"exclude resources by file path")

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

	// Define --no-simplify flag.
	noSimplify := cmd.Flags().Bool(
		"no-simplify",
		false,
		"skip simplifying resource properties",
	)

	// Define --output flag.
	output := cmd.Flags().StringP(
		"output",
		"o",
		"",
		"output format (json,name,path,references,table,yaml)")

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

		if !*noSimplify {
			simplifyResources(results)
		}

		sortResources(results)

		return state.printerFn(os.Stdout, results)
	}

	return cmd
}

// simplifyResources removes a number of properties (specifically properties
// that are automatically sey by Kubernetes after a resource is admitted) from
// each item in the given resources.Resource list.
func simplifyResources(items []resources.Resource) {
	for _, item := range items {
		unstructured.RemoveNestedField(item.Object, "metadata", "annotations", "kubectl.kubernetes.io/last-applied-configuration")
		unstructured.RemoveNestedField(item.Object, "metadata", "creationTimestamp")
		unstructured.RemoveNestedField(item.Object, "metadata", "generation")
		unstructured.RemoveNestedField(item.Object, "metadata", "managedFields")
		unstructured.RemoveNestedField(item.Object, "metadata", "ownerReferences")
		unstructured.RemoveNestedField(item.Object, "metadata", "resourceVersion")
		unstructured.RemoveNestedField(item.Object, "metadata", "uid")
		unstructured.RemoveNestedField(item.Object, "ownerReferences")
		unstructured.RemoveNestedField(item.Object, "status")
	}
}

// sortResources sorts the given resources.Resource list by filename, then
// namespace, name, and finally kind.
func sortResources(items []resources.Resource) {
	slices.SortFunc(items, func(a, b resources.Resource) int {
		if cmp := strings.Compare(a.GetFilename(), b.GetFilename()); cmp != 0 {
			return cmp
		}

		if cmp := strings.Compare(a.GetNamespace(), b.GetNamespace()); cmp != 0 {
			return cmp
		}

		if cmp := strings.Compare(a.GetName(), b.GetName()); cmp != 0 {
			return cmp
		}

		if cmp := strings.Compare(a.GetKind(), b.GetKind()); cmp != 0 {
			return cmp
		}

		return 0
	})
}
