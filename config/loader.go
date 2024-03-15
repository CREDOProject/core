package config

import (
	"credo/modules"

	"gopkg.in/yaml.v3"
)

// FromFile parses the configuration file and outputs a modules.Config.
func FromFile(configFile []byte) (modules.Config, error) {
	var config modules.Config

	error := yaml.Unmarshal(configFile, &config)
	if error != nil {
		return modules.Config{}, error
	}
	return config, nil
}
