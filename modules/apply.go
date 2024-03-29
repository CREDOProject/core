package modules

import (
	"github.com/spf13/cobra"
)

const applyModuleName = "apply"

// Registers the applyModule.
func init() { Register(applyModuleName, func() Module { return &applyModule{} }) }

// applyModule is used to apply the credospell configuration in the current
// working directory.
type applyModule struct{}

func (m *applyModule) bulkRun(c *Config) error {
	// Iterates ovet the modules to call the bulkRun function of each Module.
	for k := range Modules {
		module := Modules[k]() // Returns the module.
		err := module.BulkRun(c)
		if err != nil {
			return err
		}
	}
	return nil
}

// CliConfig implements Module.
func (m *applyModule) CliConfig(conifig *Config) *cobra.Command {
	return &cobra.Command{
		Use:   applyModuleName,
		Short: "Applies the credospell.yaml configuration in the current directory.",
		Run: func(cmd *cobra.Command, args []string) {
			m.bulkRun(conifig)
		},
		Args: cobra.NoArgs,
	}
}

// This is a stub method. It should always return nil.
func (m *applyModule) Commit(config *Config, result any) error { return nil }

// This is a stub method. It should always return nil.
func (m *applyModule) Run(anySpell any) error { return nil }

// This is a stub method. It should always return nil.
func (m *applyModule) BulkRun(config *Config) error { return nil }
