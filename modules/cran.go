package modules

import (
	"strings"

	"github.com/spf13/cobra"
)

const cranModuleName = "cran"

const cranModuleShort = "Retrieves a CRAN package and its dependencies."

const cranModuleExample = `
`

// Registers the carnModule.
func init() { Register(cranModuleName, func() Module { return &cranModule{} }) }

// cranModule is used to manage the CARN scope in the credospell configuration.
type cranModule struct{}

type cranSpell struct {
	packageName      string
	packageDirectory string
	repository       string
}

// equals checks if two cranSpell objects are equal.
func (c cranSpell) equals(t equatable) bool {
	// TODO: implement equality check.
	s, ok := t.(cranSpell)
	if !ok {
		return false
	}
	return strings.Compare(s.packageName, c.packageName) == 0 &&
		strings.Compare(s.repository, c.repository) == 0
}

// BulkRun implements Module.
func (c *cranModule) BulkRun(config *Config) error {
	for _, cs := range config.Cran {
		if err := c.Run(cs); err != nil {
			return err
		}
	}
	return nil
}

// cobraArgs is used to validate the arguments passed to the cran command.
//
// This function is intended to be used by cobra.
func (c *cranModule) cobraArgs() func(*cobra.Command, []string) error {
	return func(c *cobra.Command, s []string) error {
		return nil
	}
}

// cobraRun is used to run the module from the command line.
// It serves as an entry point to the cranModule.
//
// This function is inteded to be used by cobra.
func (c *cranModule) cobraRun(_ *Config) func(*cobra.Command, []string) {
	return func(c *cobra.Command, s []string) {
		// TODO: Implement cobraRun
	}
}

// CliConfig implements Module.
func (c *cranModule) CliConfig(config *Config) *cobra.Command {
	return &cobra.Command{
		Args:    c.cobraArgs(),
		Example: pipModuleExample,
		Run:     c.cobraRun(config),
		Short:   cranModuleShort,
		Use:     cranModuleName,
	}
}

// Commit implements Module.
func (c *cranModule) Commit(config *Config, result any) error {
	// TODO: implement Commit
	return nil
}

// Run implements Module.
func (c *cranModule) Run(any) error {
	// TODO: implement Run
	return nil
}
