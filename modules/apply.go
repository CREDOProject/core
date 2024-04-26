package modules

import "github.com/spf13/cobra"

const applyModuleName = "apply"

func init() { Register(applyModuleName, func() Module { return applyModule{} }) }

type applyModule struct{}

// This is a stub method. It should always return nil.
func (a applyModule) Apply(any) error {
	return nil
}

// This is a stub method. It should always return nil.
func (a applyModule) BulkApply(config *Config) error {
	return nil
}

// This is a stub method. It should always return nil.
func (a applyModule) BulkSave(config *Config) error {
	return nil
}

// CliConfig implements Module.
func (a applyModule) CliConfig(config *Config) *cobra.Command {
	return &cobra.Command{
		Use:   applyModuleName,
		Short: "Runs the credospell.yaml configuration in the current directory and installs all the dependencies.",
		Run: func(cmd *cobra.Command, args []string) {
		},
		Args: cobra.NoArgs,
	}
}

// This is a stub method. It should always return nil.
func (a applyModule) Commit(config *Config, result any) error {
	return nil
}

// This is a stub method. It should always return nil.
func (a applyModule) Save(any) error {
	return nil
}
