package modules

import (
	"errors"
	"log"
)

type PipModule struct {
	logger *log.Logger
}

type PipSpell struct {
	Name string `yaml:"name"`
}

func (m *PipModule) Marshaler() interface{} {
	return PipSpell{}
}

func (m *PipModule) Commit(config *Config, result any) error {
	v := result.(PipSpell)
	config.Pip = append(config.Pip, v)
	return nil
}

func (m *PipModule) BareRun(p *Parameters) any {
	spell, err := m.bareRun(p)
	if err != nil {
		m.logger.Fatal(err)
	}
	return spell
}

func (m *PipModule) bareRun(p *Parameters) (PipSpell, error) {
	if len(p.Env) < 1 {
		return PipSpell{}, errors.New("Pip module requires at least one parameter.")
	}
	// Setup a spell entry.
	spell := PipSpell{
		Name: p.Env["name"],
	}

	return spell, nil
}

func init() {
	Register("pip", func() Module {
		return &PipModule{}
	})
}
