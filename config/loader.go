package config

import (
	"credo/modules"

	"gopkg.in/yaml.v3"
)

// Parses configuration file and outputs a Config structure.
func FromFile(configFile []byte) (modules.Config, error) {
	var config modules.Config

	error := yaml.Unmarshal(configFile, &config)
	if error != nil {
		return modules.Config{}, error
	}
	return config, nil
}
