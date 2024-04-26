package modules

import "github.com/spf13/cobra"

const applyModuleName = "apply"

func init() { Register(applyModuleName, func() Module { return applyModule{} }) }

type applyModule struct{}

// Apply implements Module.
func (a applyModule) Apply(any) error {
	panic("unimplemented")
}

// BulkApply implements Module.
func (a applyModule) BulkApply(config *Config) error {
	panic("unimplemented")
}

// BulkSave implements Module.
func (a applyModule) BulkSave(config *Config) error {
	panic("unimplemented")
}

// CliConfig implements Module.
func (a applyModule) CliConfig(config *Config) *cobra.Command {
	panic("unimplemented")
}

// Commit implements Module.
func (a applyModule) Commit(config *Config, result any) error {
	panic("unimplemented")
}

// Save implements Module.
func (a applyModule) Save(any) error {
	panic("unimplemented")
}
