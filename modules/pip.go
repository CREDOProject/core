package modules

import (
	"credo/logger"
	"errors"
	"log"
	"os"

	gopip "github.com/CREDOProject/go-pip"
	"github.com/CREDOProject/go-pip/utils"
)

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

func (m *PipModule) BareRun(c *Config, p *Parameters) any {
	spell, err := m.bareRun(p)
	if err != nil {
		logger.Get().Fatal(err)
	}
	return spell
}

func (m *PipModule) bareRun(p *Parameters) (PipSpell, error) {
	if len(p.Env) < 1 {
		return PipSpell{},
			errors.New("Pip module requires at least one parameter.")
	}
	// Setup a spell entry.
	spell := PipSpell{
		Name: p.Env["name"],
	}

	pipBinary, err := utils.DetectPipBinary()
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

func init() { Register("pip", func() Module { return &PipModule{} }) }
