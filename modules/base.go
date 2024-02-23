package modules

import (
	"credo/logger"
)

type Parameters struct {
	Env map[string]string
}

type Module interface {
	// BareRun is used to run a module without making any change to the
	// file system other than adding an entry to the credospell file.
	// It returns a spell entry of a module.
	BareRun(*Config, *Parameters) any

	// Marshaler returns the interface used by a module to specify its
	// parameters.
	Marshaler() interface{}

	// Commit adds a configuration entry for a said module.
	Commit(config *Config, result any) error

	// Run is used to execute a Module making changes to the filesystem.
	Run(any) error

	// BulkRun is used to run the config entry ofa each sub-entry of a module.
	BulkRun(config *Config) error
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
