package modules

import (
	"credo/logger"
	"credo/project"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/CREDOProject/go-apt-client"
	"github.com/spf13/cobra"
)

const aptModuleName = "apt"

const aptModuleShort = "Retrieves an apt package and its depenencies."

const aptModuleExample = `
Install a apt package:
	credo apt python3
`

func init() { Register(aptModuleName, func() Module { return &aptModule{} }) }

type aptModule struct{}

type aptSpell struct {
	Name         string     `yaml:"name"`
	Depencencies []aptSpell `yaml:"dependencies,omitempty"`
}

// Function to check equality of two aptSpells
func (a aptSpell) equals(t equatable) bool {
	if o, ok := t.(aptSpell); ok {
		equality := len(o.Depencencies) == len(a.Depencencies)
		if !equality {
			return false
		}
		for i := range o.Depencencies {

			equality = equality &&
				strings.Compare(
					o.Depencencies[i].Name, a.Depencencies[i].Name) == 0
		}
		return strings.Compare(a.Name, o.Name) == 0
	}
	return false
}

// BulkRun implements Module.
func (m *aptModule) BulkRun(config *Config) error {
	for _, as := range config.Apt {
		for _, dep := range as.Depencencies {
			err := m.Run(dep)
			if err != nil {
				return err
			}
		}
		err := m.Run(as)
		if err != nil {
			return err
		}
	}
	return nil
}

// CliConfig implements Module.
func (m *aptModule) CliConfig(config *Config) *cobra.Command {
	return &cobra.Command{
		Args:    m.cobraArgs(),
		Example: aptModuleExample,
		Run:     m.cobraRun(config),
		Short:   aptModuleShort,
		Use:     aptModuleName,
	}
}

// Function used to validate the arguments passed to the apt command.
// If no arguments are passed, it returns an error.
// Otherwise it returns nil.
//
// Intended to be used by cobra.
func (m *aptModule) cobraArgs() func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("%s module requires at least one argument.",
				aptModuleName)
		}
		return nil
	}
}

// Function used to run the module from the command line.
// It serves as an entry point to the bare run of the aptModule.
//
// Intended to be used by cobra.
func (m *aptModule) cobraRun(config *Config) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		name := args[0]
		spell, err := m.bareRun(aptSpell{
			Name: name,
		})
		if err != nil {
			logger.Get().Fatal(err)
		}
		err = m.Commit(config, spell)
		if err != nil {
			logger.Get().Fatal(err)
		}
	}
}

func (*aptModule) bareRun(spell aptSpell) (aptSpell, error) {
	aptPack := &apt.Package{
		Name: spell.Name,
	}
	output, err := apt.InstallDry(aptPack)
	logger.Get().Print(string(output))
	if err != nil {
		return aptSpell{}, err
	}
	depList, err := apt.GetDependencies(aptPack)
	if err != nil {
		return aptSpell{}, err
	}
	for _, dependency := range depList {
		spell.Depencencies = append(spell.Depencencies, aptSpell{
			Name: dependency,
		})
	}
	return spell, nil
}

// Commit implements Module.
func (*aptModule) Commit(config *Config, result any) error {
	newEntry := result.(aptSpell)
	if Contains(config.Apt, newEntry) {
		return ErrAlreadyPresent
	}
	config.Apt = append(config.Apt, newEntry)
	return nil
}

// Run implements Module.
func (*aptModule) Run(anySpell any) error {
	spell, ok := anySpell.(aptSpell)
	if !ok {
		return fmt.Errorf("Error converting to aptSpell")
	}
	project, err := project.ProjectPath()
	if err != nil {
		return err
	}
	downloadPath := path.Join(*project, aptModuleName)
	os.MkdirAll(downloadPath, 0755)
	aptPack := &apt.Package{
		Name: spell.Name,
	}
	out, err := apt.Download(aptPack, downloadPath)
	logger.Get().Print(string(out))
	return err
}
