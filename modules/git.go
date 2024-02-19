package modules

import (
	"errors"
	"log"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

type GitModule struct {
	logger *log.Logger
}

type GitSpell struct {
	URL     string `yaml:"url"`
	Version string `yaml:"version"`
}

func (m *GitModule) Marshaler() interface{} {
	return GitSpell{}
}

func (m *GitModule) Commit(config *Config, result any) error {
	v := result.(GitSpell)
	config.Git = append(config.Git, v)
	return nil
}

func (m *GitModule) BareRun(p *Parameters) any {
	spell, err := m.bareRun(p)
	if err != nil {
		m.logger.Fatal(err)
	}
	return spell
}

func (m *GitModule) bareRun(p *Parameters) (GitSpell, error) {
	if len(p.Env) < 1 {
		return GitSpell{}, errors.New("Git module requires at least one parameter.")
	}
	// Logic to get the latest version or the specified version.
	version := p.Env["version"]
	if len(version) == 0 {
		version = "HEAD"
	}

	// Setup a spell entry.
	spell := GitSpell{
		URL:     p.Env["url"],
		Version: version,
	}

	// Try Cloning
	_, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:               p.Env["url"],
		Depth:             1,
		SingleBranch:      true,
		RecurseSubmodules: 1,
		Tags:              0,
	})

	if err != nil {
		return spell, err
	}

	return spell, nil
}

func init() {
	Register("git", func() Module {
		return &GitModule{}
	})
}
