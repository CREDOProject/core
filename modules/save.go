package modules

import (
	"credo/logger"

	"github.com/spf13/cobra"
)

const saveModuleName = "save"

// Registers the applyModule.
func init() { Register(saveModuleName, func() Module { return &saveModule{} }) }

// CliConfig implements Module.
func (m *saveModule) CliConfig(config *Config) *cobra.Command {
	return &cobra.Command{
		Use:   saveModuleName,
		Short: "Runs the credospell.yaml configuration in the current directory and saves every dependency.",
		Run: func(cmd *cobra.Command, args []string) {
			for k := range Modules {
				module := Modules[k]()
				err := module.BulkSave(config)
				if err != nil {
					logger.Get().Fatal(err)
				}
			}
		},
		Args: cobra.NoArgs,
	}
}

// saveModule is used to apply the credospell configuration in the current
// working directory.
type saveModule struct{}

// This is a stub method. It should always return nil.
func (m *saveModule) Apply(any) error { return nil }

// This is a stub method. It should always return nil.
func (m *saveModule) BulkApply(config *Config) error { return nil }

// This is a stub method. It should always return nil.
func (m *saveModule) Commit(config *Config, result any) error { return nil }

// This is a stub method. It should always return nil.
func (m *saveModule) Save(anySpell any) error { return nil }

// This is a stub method. It should always return nil.
func (m *saveModule) BulkSave(config *Config) error { return nil }
