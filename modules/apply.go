package modules

import (
	"credo/project"

	"github.com/spf13/cobra"
)

type ApplyModule struct{}

func (m *ApplyModule) Commit(config *Config, result any) error {
	panic("unimplemented")
}

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

func (m *ApplyModule) Run(anySpell any) error {
	panic("unimplemented")
}

func (m *ApplyModule) BulkRun(config *Config) error {
	panic("unimplemented")
}

func (m *ApplyModule) CliConfig(conifig *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "apply",
		Short: "Applies the credospell.yaml configuration in the current directory.",
		Run: func(cmd *cobra.Command, args []string) {
			m.bulkRun(conifig)
		},
		Args: cobra.NoArgs,
	}
}

func init() { Register("apply", func() Module { return &ApplyModule{} }) }
