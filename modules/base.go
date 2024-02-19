package modules

import (
	"credo/logger"
)

type Parameters struct {
	Env map[string]string
}

type Module interface {
	// Function used to bare run the Module.
	BareRun(*Parameters) any
	Marshaler() interface{}
	Commit(config *Config, result any) error
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
