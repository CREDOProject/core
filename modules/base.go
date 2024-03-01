package modules

import (
	"credo/logger"
	"errors"

	"github.com/spf13/cobra"
)

var (
	ErrAlreadyPresent = errors.New("Entry already present.")
)

type equatable interface {
	equals(T equatable) bool
}

// Returns true if s contains e, returns false otherwise.
func Contains[T equatable](s []T, e T) bool {
	for _, v := range s {
		if v.equals(e) {
			return true
		}
	}
	return false
}

// Factory to provide a closure to get the Module.
type ModuleFactory = func() Module

// Module registry.
var Modules = map[string]ModuleFactory{}

type Module interface {
	// Commit adds a configuration entry for a said module.
	Commit(config *Config, result any) error

	// Run is used to execute a Module making changes to the filesystem.
	Run(any) error

	// BulkRun is used to run the config entry ofa each sub-entry of a module.
	BulkRun(config *Config) error

	// Returns a cobra.Command to use in the command line.
	CliConfig(config *Config) *cobra.Command
}

// Register SHOULD BE called by the init() function of a provider.
func Register(moduleName string, module ModuleFactory) {
	if _, present := Modules[moduleName]; present {
		logger.Get().Fatalf("Module %s already defined.", moduleName)
	}
	Modules[moduleName] = module
}
