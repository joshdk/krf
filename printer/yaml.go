// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package printer

import (
	"io"

	"sigs.k8s.io/yaml"

	"github.com/joshdk/krf/resources"
)

// YAML prints each given resources.Resource as a stream of yaml documents.
func YAML(w io.Writer, items []resources.Resource) error {
	for i, item := range items {
		if i != 0 {
			// Print a yaml document separator only after the first resources.
			if _, err := w.Write([]byte("---\n")); err != nil {
				return err
			}
		}

		output, err := yaml.Marshal(item.Object)
		if err != nil {
			return err
		}

		if _, err := w.Write(output); err != nil {
			return err
		}
	}

	return nil
}
