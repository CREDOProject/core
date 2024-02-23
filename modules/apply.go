package modules

import (
	"credo/logger"
	"log"
)

type ApplyModule struct {
	logger *log.Logger
}

func (m *ApplyModule) Marshaler() interface{} {
	return nil
}

func (m *ApplyModule) Commit(config *Config, result any) error {
	return nil
}

func (m *ApplyModule) BareRun(c *Config, p *Parameters) any {
	err := m.bareRun(c, p)
	if err != nil {
		logger.Get().Fatal(err)
	}
	return nil
}

// This bare run is a real run.
func (m *ApplyModule) bareRun(c *Config, p *Parameters) error {

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
	return nil
}

func (m *ApplyModule) BulkRun(config *Config) error {
	return nil
}

func init() { Register("apply", func() Module { return &ApplyModule{} }) }
