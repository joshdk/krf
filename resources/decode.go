// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package resources contains functions for decoding Kubernetes resource
// objects from a number of different sources, in a number of different file
// formats.
// Specifically, these formats are supported:
// - A single resource in yaml/json format.
//   - Produced by running "kubectl get" for a single resource.
//
// - A v1.List containing multiple resources in yaml/json format.
//   - Produced by running "kubectl get" for a list of resources.
//
// - A stream of yaml/json documents.
//   - Produced by running "kustomize build".
package resources

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Resource represents a single Kubernetes resource. It holds the original
// unstructured object contents, as well as (optionally) the file filename where
// that object was originally decoded from.
type Resource struct {
	// Unstructured is the decoded resource object.
	unstructured.Unstructured

	// filename where the resource was originally decoded from. This value is
	// only set if the resource was decoded from a file (opposed to being
	// decoded from an io.Reader).
	filename string
}

// GetFilename returns the filename from which this resource was originally
// decoded. Exists to align with the other unstructured.Unstructured helper
// functions.
func (i Resource) GetFilename() string {
	return i.filename
}

// ResourceFunc is a callback function that is passed each resource encountered
// while decoding.
type ResourceFunc func(Resource)

// Decode decodes Kubernetes resources from the given generic source. The
// ResourceFunc callback is executed with each decoded resource.
//
// Behavior notes:
// - If the string "" or "-" is given, resources are read from os.Stdin.
func Decode(source any, handler ResourceFunc) error {
	switch s := source.(type) {
	case io.Reader:
		return Reader(s, handler)

	case string:
		switch s {
		case "", "-":
			return Reader(os.Stdin, handler)

		default:
			if fi, err := os.Stat(s); err != nil {
				return err
			} else if fi.IsDir() {
				return Directory(s, handler)
			}

			return File(s, handler)
		}

	default:
		return fmt.Errorf("unsupported source type: %T", source)
	}
}

// Directory decodes Kubernetes resources from files discovered while walking
// the given directory. The ResourceFunc callback is executed with each decoded
// resource.
//
// Behavior notes:
// - Any directories named ".git" or "node_modules" are skipped.
// - Any files not ending with ".yaml" are skipped.
// - Any decoding errors are ignored.
func Directory(directory string, handler ResourceFunc) error {
	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		switch {
		case err != nil:
			// There was an actual error.
			return err

		case info.IsDir():
			// Ignore anything that isn't a file.
			switch info.Name() {
			case ".git":
				// Completely stop recursing into the .git directory.
				return filepath.SkipDir
			case "node_modules":
				// Completely stop recursing into the node_modules directory.
				return filepath.SkipDir
			default:
				return nil
			}

		case !strings.HasSuffix(info.Name(), ".yaml"):
			// Ignore anything that isn't a .yaml file.
			return nil

		default:
			_ = File(path, handler)

			// Intentionally do not propagate errors encountered while decoding
			// a discovered file. Walking through an arbitrary directory can
			// commonly result in the attempted decoding of invalid yaml files
			// (Helm charts for example). This should not interrupt the
			// continued walking of the directory tree.
			return nil
		}
	})
}

// File decodes Kubernetes resources from the given filename. The ResourceFunc
// callback is executed with each decoded resource.
func File(filename string, handler ResourceFunc) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close() //nolint:errcheck

	return Reader(file, func(item Resource) {
		item.filename = filename
		handler(item)
	})
}

// Reader decodes Kubernetes resources from the given io.Reader. The
// ResourceFunc callback is executed with each decoded resource.
func Reader(reader io.Reader, handler ResourceFunc) error {
	return decodeReader(reader, func(uu unstructured.Unstructured) {
		handler(Resource{Unstructured: uu})
	})
}

func decodeReader(reader io.Reader, handler func(unstructured.Unstructured)) error {
	decoder := yaml.NewYAMLOrJSONDecoder(reader, 100) //nolint:mnd

	for {
		// Attempt to decode a single object from the stream.
		var uu unstructured.Unstructured
		if err := decoder.Decode(&uu.Object); err != nil {
			if errors.Is(err, io.EOF) {
				// No more objects in the stream.
				return nil
			}

			return err
		}

		// Verify that this object has the minimum set of properties to even be
		// considered a Kubernetes resource.
		// We don't resourceCheck the name value as that is omitted from v1.List
		// resources.
		// We don't resourceCheck the namespace value as that is commonly omitted from
		// local manifests.
		if uu.GetAPIVersion() == "" || uu.GetKind() == "" {
			continue
		}

		// This object does not contain a list of other objects. Handle object
		// directly.
		if !uu.IsList() {
			// Verify that this (non-v1.List) object has a name value.
			if uu.GetName() == "" {
				continue
			}

			handler(uu)

			continue
		}

		// This object does contain a list of other objects. Handle each
		// contained object individually.
		ul, err := uu.ToList()
		if err != nil {
			return err
		}

		for _, uli := range ul.Items {
			// Verify that this (v1.List item) object has a name value.
			if uli.GetName() == "" {
				continue
			}

			handler(uli)
		}
	}
}
