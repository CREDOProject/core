package modules

import (
	"credo/logger"

	"github.com/spf13/cobra"
)

type Module interface {
	// Marshaler returns the interface used by a module to specify its
	// parameters.
	Marshaler() interface{}

	// Commit adds a configuration entry for a said module.
	Commit(config *Config, result any) error

	// Run is used to execute a Module making changes to the filesystem.
	Run(any) error

	// BulkRun is used to run the config entry ofa each sub-entry of a module.
	BulkRun(config *Config) error

	// Returns a cobra.Command to use in the command line.
	CliConfig(config *Config) *cobra.Command
}

// Factory to provide a closure to get the Module.
type ModuleFactory = func() Module

// Module registry.
var Modules = map[string]ModuleFactory{}

// Register SHOULD BE called by the init() function of a provider.
func Register(moduleName string, module ModuleFactory) {
	if _, present := Modules[moduleName]; present {
		logger.Get().Fatalf("Module %s already defined.", moduleName)
	}
	Modules[moduleName] = module
}
