package modules

import "github.com/spf13/cobra"

const dnfModuleName = "dnf"

func init() { Register(dnfModuleName, func() Module { return &dnfModule{} }) }

// dnfModule is used to manage the dnf scope in the credospell configuration.
type dnfModule struct{}

// Apply implements Module.
func (d *dnfModule) Apply(any) error {
	panic("unimplemented")
}

// BulkApply implements Module.
func (d *dnfModule) BulkApply(config *Config) error {
	panic("unimplemented")
}

// BulkSave implements Module.
func (d *dnfModule) BulkSave(config *Config) error {
	panic("unimplemented")
}

// CliConfig implements Module.
func (d *dnfModule) CliConfig(config *Config) *cobra.Command {
	panic("unimplemented")
}

// Commit implements Module.
func (d *dnfModule) Commit(config *Config, result any) error {
	panic("unimplemented")
}

// Save implements Module.
func (d *dnfModule) Save(any) error {
	panic("unimplemented")
}
