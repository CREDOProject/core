package modules

import (
	"github.com/spf13/cobra"
)

const saveModuleName = "save"

// Registers the applyModule.
func init() { Register(saveModuleName, func() Module { return &saveModule{} }) }

// saveModule is used to apply the credospell configuration in the current
// working directory.
type saveModule struct{}

// Apply implements Module.
func (m *saveModule) Apply(any) error {
	panic("unimplemented")
}

// BulkApply implements Module.
func (m *saveModule) BulkApply(config *Config) error {
	panic("unimplemented")
}

func (m *saveModule) bulkRun(c *Config) error {
	// Iterates ovet the modules to call the bulkRun function of each Module.
	for k := range Modules {
		module := Modules[k]() // Returns the module.
		err := module.BulkSave(c)
		if err != nil {
			return err
		}
	}
	return nil
}

// CliConfig implements Module.
func (m *saveModule) CliConfig(conifig *Config) *cobra.Command {
	return &cobra.Command{
		Use:   saveModuleName,
		Short: "Runs the credospell.yaml configuration in the current directory and saves every dependency.",
		Run: func(cmd *cobra.Command, args []string) {
			m.bulkRun(conifig)
		},
		Args: cobra.NoArgs,
	}
}

// This is a stub method. It should always return nil.
func (m *saveModule) Commit(config *Config, result any) error { return nil }

// This is a stub method. It should always return nil.
func (m *saveModule) Save(anySpell any) error { return nil }

// This is a stub method. It should always return nil.
func (m *saveModule) BulkSave(config *Config) error { return nil }
