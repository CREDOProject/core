package modules

import (
	"credo/logger"
	"credo/project"
	"fmt"
	"log"
	"os"
	"path"

	gopip "github.com/CREDOProject/go-pip"
	"github.com/CREDOProject/go-pip/utils"
	pythonvenv "github.com/CREDOProject/go-pythonvenv"
	"github.com/spf13/cobra"
)

const pipModuleName = "pip"

const pipModuleExample = `
Install a pip package:
	credo pip numpy

Install a pip package pinning it to a version:
	credo pip numpy==1.26.0
`

type PipModule struct {
	logger *log.Logger
}

type PipSpell struct {
	Name string `yaml:"name"`
}

// Function to check equality of two PipSpells
func (s *PipSpell) equals(t *PipSpell) bool {
	return s.Name == t.Name
}

func (m *PipModule) Marshaler() interface{} {
	return PipSpell{}
}

func (m *PipModule) Commit(config *Config, result any) error {
	newEntry := result.(PipSpell)

	for _, spell := range config.Pip {
		if spell.equals(&newEntry) {
			return nil // Break from the for loop.
		}
	}

	config.Pip = append(config.Pip, newEntry)
	return nil
}

func setupPythonVenv(path string) (string, error) {
	venv, err := pythonvenv.Create(path)
	if err != nil {
		return "", err
	}
	return venv.Path, nil
}

func (m *PipModule) bareRun(p PipSpell) (PipSpell, error) {
	// Setup a spell entry.
	spell := PipSpell{
		Name: p.Name,
	}

	// Obtain the project path
	projectPath, err := project.ProjectPath()
	if err != nil {
		return PipSpell{}, err
	}

	venvPath, err := setupPythonVenv(path.Join(*projectPath, "venv"))
	if err != nil {
		return PipSpell{}, err
	}

	pipBinary, err := utils.PipBinaryFrom(path.Join(venvPath, "bin"))
	if err != nil {
		return PipSpell{}, err
	}

	cmd, err := gopip.New(pipBinary).Install(spell.Name).DryRun().Seal()
	if err != nil {
		return PipSpell{}, err
	}

	err = cmd.Run(&gopip.RunOptions{
		Output: os.Stdout,
	})
	if err != nil {
		return PipSpell{}, err
	}

	return spell, nil
}

func (m *PipModule) Run(anySpell any) error {
	return nil
}

func (m *PipModule) BulkRun(config *Config) error {
	for _, ps := range config.Pip {
		err := m.Run(ps)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *PipModule) CliConfig(conifig *Config) *cobra.Command {
	return &cobra.Command{
		Use:     pipModuleName,
		Short:   "Retrieves a python package.",
		Example: pipModuleExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("%s module requires at least one argument.",
					pipModuleName)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			spell, err := m.bareRun(PipSpell{
				Name: args[0],
			})
			if err != nil {
				logger.Get().Fatal(err)
			}
			err = m.Commit(conifig, spell)
			if err != nil {
				logger.Get().Fatal(err)
			}
		},
	}
}

func init() { Register(pipModuleName, func() Module { return &PipModule{} }) }
