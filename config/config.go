// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package config provides functionality for initializing and parsing a yaml
// configuration file.
package config

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"

	"github.com/joshdk/krf/resolver"
)

// Configuration represents the contents of a krf configuration file.
type Configuration struct {
	Resources []resolver.Resource `yaml:"resources"`
}

//go:embed files/configuration.yaml
var configurationData []byte

// InitAndLoad creates the named configuration file if it does not exist and
// then subsequently loads it.
func InitAndLoad(filename string) (*Configuration, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(filename), 0o700); err != nil {
			return nil, err
		}

		if err := os.WriteFile(filename, configurationData, 0o600); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return Load(filename)
}

// Load parses the named configuration file.
func Load(filename string) (*Configuration, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg struct {
		Configuration `yaml:",inline"`

		APIVersion string `yaml:"apiVersion"`
		Kind       string `yaml:"kind"`
	}

	if err := yaml.UnmarshalStrict(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.APIVersion != "krf.joshdk.github.com/v1beta1" {
		return nil, fmt.Errorf("unsupported apiversion: %s", cfg.APIVersion)
	}

	if cfg.Kind != "Configuration" {
		return nil, fmt.Errorf("unsupported kind %s", cfg.Kind)
	}

	return &cfg.Configuration, nil
}
