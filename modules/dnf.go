package modules

import (
	"credo/logger"
	"fmt"

	"github.com/spf13/cobra"
)

const dnfModuleName = "dnf"

const dnfModuleShort = "Retrieves a dnf package and its dependencies"

const dnfModuleExample = `
Install a dnf package
	credo dnf vim
`

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
	return &cobra.Command{
		Args:    d.cobraArgs(),
		Example: dnfModuleExample,
		Run:     d.cobraRun(config),
		Short:   dnfModuleShort,
		Use:     dnfModuleName,
	}
}

// Function used to run the module from the command line.
// It serves as an entry point to the bare run of the dnfModule.
//
// Intended to be used by cobra.
func (d *dnfModule) cobraRun(config *Config) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		// TODO: implement run
		err := d.Commit(config, nil)
		if err != nil && err != ErrAlreadyPresent {
			logger.Get().Fatal(err)
		}
	}
}

// Function used to validate the arguments passed to the dnf command.
// If no arguments are passed, it returns an error.
// Otherwise it returns nil.
//
// Intended to be used by cobra.
func (*dnfModule) cobraArgs() func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("%s module requires at least one argument.",
				aptModuleName)
		}
		return nil
	}
}

// Commit implements Module.
func (d *dnfModule) Commit(config *Config, result any) error {
	panic("unimplemented")
}

// Save implements Module.
func (d *dnfModule) Save(any) error {
	panic("unimplemented")
}
