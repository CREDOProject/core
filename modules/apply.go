package modules

import (
	"credo/project"

	"github.com/spf13/cobra"
)

const applyModuleName = "apply"

func init() { Register(applyModuleName, func() Module { return &ApplyModule{} }) }

type ApplyModule struct{}

func (m *ApplyModule) bulkRun(c *Config) error {
	_, err := project.ProjectPath()
	if err != nil {
		return err
	}
	for k := range Modules {
		module := Modules[k]()
		err := module.BulkRun(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *ApplyModule) CliConfig(conifig *Config) *cobra.Command {
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
func (m *ApplyModule) Commit(config *Config, result any) error {
	return nil
}

// This is a stub method. It should always return nil.
func (m *ApplyModule) Run(anySpell any) error {
	return nil
}

// This is a stub method. It should always return nil.
func (m *ApplyModule) BulkRun(config *Config) error {
	return nil
}
