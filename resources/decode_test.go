// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package resources_test

import (
	"os"
	"testing"

	"github.com/joshdk/krf/resources"
)

func TestDecode(t *testing.T) { //nolint:funlen
	tests := map[string]struct {
		source   string
		setup    func() func()
		expected []resourceCheck
	}{
		"yaml": {
			source: "testdata/multiple.yaml",
			expected: []resourceCheck{
				{
					name:     "Deployment/nginx-deployment",
					filename: "testdata/multiple.yaml",
				},
				{
					name:     "Service/my-service",
					filename: "testdata/multiple.yaml",
				},
				{
					name:     "ConfigMap/myconfigmap",
					filename: "testdata/multiple.yaml",
				},
			},
		},

		"json": {
			source: "testdata/multiple.json",
			expected: []resourceCheck{
				{
					name:     "Deployment/nginx-deployment",
					filename: "testdata/multiple.json",
				},
				{
					name:     "Service/my-service",
					filename: "testdata/multiple.json",
				},
				{
					name:     "ConfigMap/myconfigmap",
					filename: "testdata/multiple.json",
				},
			},
		},

		"list": {
			source: "testdata/list.yaml",
			expected: []resourceCheck{
				{
					name:     "Deployment/nginx-deployment",
					filename: "testdata/list.yaml",
				},
				{
					name:     "Service/my-service",
					filename: "testdata/list.yaml",
				},
				{
					name:     "ConfigMap/myconfigmap",
					filename: "testdata/list.yaml",
				},
			},
		},

		"directory": {
			source: "testdata/directory",
			expected: []resourceCheck{
				{
					name:     "Deployment/nginx-deployment",
					filename: "testdata/directory/multiple.yaml",
				},
				{
					name:     "Service/my-service",
					filename: "testdata/directory/multiple.yaml",
				},
				{
					name:     "ConfigMap/myconfigmap",
					filename: "testdata/directory/subdirectory/single.yaml",
				},
			},
		},

		"stdin": {
			source: "-",
			setup: func() func() {
				os.Stdin, _ = os.Open("testdata/multiple.yaml")

				return func() {
					os.Stdin.Close() //nolint:errcheck
				}
			},
			expected: []resourceCheck{
				{name: "Deployment/nginx-deployment"},
				{name: "Service/my-service"},
				{name: "ConfigMap/myconfigmap"},
			},
		},
	}

	t.Parallel()

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if test.setup != nil {
				cleanup := test.setup()
				defer cleanup()
			}

			var items []resources.Resource

			err := resources.Decode(test.source, func(item resources.Resource) {
				items = append(items, item)
			})
			if err != nil {
				t.Fatal(err)
			}

			requireResources(t, items, test.expected)
		})
	}
}

type resourceCheck struct {
	name     string
	filename string
}

func requireResources(t *testing.T, items []resources.Resource, checks []resourceCheck) {
	t.Helper()

	if len(items) != len(checks) {
		t.Fatalf("expected %d items, got %d", len(checks), len(items))
	}

	for i := range items {
		item := items[i]
		check := checks[i]
		name := item.GetKind() + "/" + item.GetName()

		switch {
		case check.name != name:
			t.Fatalf("expected name %s, got %s", check.name, name)
		case check.filename != item.GetFilename():
			t.Fatalf("expected filename %s, got %s", check.filename, item.GetFilename())
		}
	}
}
