package modules

import "github.com/spf13/cobra"

const condaModuleName = "conda"

func init() { Register(condaModuleName, func() Module { return &condaModule{} }) }

type condaModule struct{}

type condaSpell struct{}

// BulkRun implements Module.
func (c *condaModule) BulkRun(config *Config) error {
	return nil
}

// CliConfig implements Module.
func (c *condaModule) CliConfig(config *Config) *cobra.Command {
	return nil
}

// Commit implements Module.
func (c *condaModule) Commit(config *Config, result any) error {
	return nil
}

// Run implements Module.
func (c *condaModule) Run(any) error {
	return nil
}
