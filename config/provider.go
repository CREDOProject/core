package config

import "credo/modules"

// Provider is an interface that abstracts a configuration provider.
type Provider interface {
	// Retrieves the configuration.
	Get() (*modules.Config, error)
	// Writes the configuration.
	Write(*modules.Config) error
}
