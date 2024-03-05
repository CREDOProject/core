package modules

// Application configuration.
type Config struct {
	Git []gitSpell `yaml:"git,omitempty"`
	Pip []pipSpell `yaml:"pip,omitempty"`
	Apt []aptSpell `yaml:"apt,omitempty"`
}
