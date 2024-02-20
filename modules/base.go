package modules

import (
	"credo/logger"
)

type Parameters struct {
	Env map[string]string
}

type Module interface {
	// Function used to do a bare run of a Module.
	BareRun(*Config, *Parameters) any
	// Function used to return the Marshaler of a Module
	Marshaler() interface{}
	// Function used to commit a Module into the configuration.
	Commit(config *Config, result any) error
	// Function used to run a Module.
	Run(any) error
	// Function used to run the config entry of a module.
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
