package modules

import (
	"credo/logger"
	"errors"

	"github.com/spf13/cobra"
)

var (
	// ErrAlreadyPresent SHOULD be used by a Module to attest that it's
	// configuration in the context of a `spell` is already present.
	ErrAlreadyPresent = errors.New("Entry already present.")

	// ErrConverting SHOULD be used by a Module to communicate an error in
	// converting a Spell.
	ErrConverting = errors.New("Error converting spell.")
)

// equatable is an interface that provides a method to check equality between
// two objects.
type equatable interface {
	// equals returns true if the target object passed is equal to the
	// calling object.
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

// Interface used to define the functionality of a module.
// A Module should implement this interface to be used in CREDO.
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

// Registers modules to a subcommand.
func RegisterModulesCli(cmd *cobra.Command, config *Config) {
	for _, module := range Modules {
		if cfg := module().CliConfig(config); cfg != nil && cmd != cfg {
			cmd.AddCommand(cfg)
		}
	}
}
