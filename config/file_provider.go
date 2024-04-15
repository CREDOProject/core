package config

import (
	"credo/modules"
	"credo/storage"

	"gopkg.in/yaml.v3"
)

// Implements Provider.
type FileProvider struct{}

// fromFile parses the configuration file and outputs a modules.Config.
func fromFile(configFile []byte) (modules.Config, error) {
	var config modules.Config

	error := yaml.Unmarshal(configFile, &config)
	if error != nil {
		return modules.Config{}, error
	}
	return config, nil
}

var _internalFileStorage *storage.FileStorage

// getFileStorage returns an instance of storage.FileStorage.
func getFileStorage() *storage.FileStorage {
	if _internalFileStorage == nil {
		_internalFileStorage = &storage.FileStorage{
			Filename: "credospell.yaml",
		}
	}
	return _internalFileStorage
}

// Get retrieves the configuration from file
func (*FileProvider) Get() (*modules.Config, error) {
	store := getFileStorage()
	prevFile := store.Read()
	fullConfig, err := fromFile(prevFile)
	return &fullConfig, err
}

// Write commits the configuration to file
func (*FileProvider) Write(config *modules.Config) error {
	store := getFileStorage()
	marshal, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	store.Write(marshal)
	return nil
}
