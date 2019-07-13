package oauth

import (
	"encoding/base64"

	"go-sdk/configutil"
	"go-sdk/env"
)

// NewConfigFromEnv creates a new config from the environment.
func NewConfigFromEnv() (*Config, error) {
	var cfg Config
	err := env.Env().ReadInto(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// MustNewConfigFromEnv returns a new config from the environment
// and panics if there is an error.
func MustNewConfigFromEnv() *Config {
	cfg, err := NewConfigFromEnv()
	if err != nil {
		panic(err)
	}
	return cfg
}

// Config is the config options.
type Config struct {
	// Secret is an encryption key used to verify oauth state.
	Secret string `json:"secret,omitempty" yaml:"secret,omitempty" env:"OAUTH_SECRET"`
	// RedirectURI is the oauth return url.
	RedirectURI string `json:"redirectURI,omitempty" yaml:"redirectURI,omitempty" env:"OAUTH_REDIRECT_URI"`
	// HostedDomain is a specific domain we want to filter identities to.
	HostedDomain string `json:"hostedDomain,omitempty" yaml:"hostedDomain,omitempty" env:"OAUTH_HOSTED_DOMAIN"`

	Scopes []string `json:"scopes,omitempty" yaml:"scopes,omitempty"`

	// ClientID is part of the oauth credential pair.
	ClientID string `json:"clientID,omitempty" yaml:"clientID,omitempty" env:"OAUTH_CLIENT_ID"`
	// ClientSecret is part of the oauth credential pair.
	ClientSecret string `json:"clientSecret,omitempty" yaml:"clientSecret,omitempty" env:"OAUTH_CLIENT_SECRET"`
}

// IsZero returns if the config is set or not.
func (c Config) IsZero() bool {
	return len(c.ClientID) == 0 || len(c.ClientSecret) == 0
}

// GetSecret gets the secret if set or a default.
func (c Config) GetSecret(defaults ...[]byte) ([]byte, error) {
	if len(c.Secret) > 0 {
		decoded, err := base64.StdEncoding.DecodeString(c.Secret)
		if err != nil {
			return nil, err
		}
		return decoded, nil
	}
	if len(defaults) > 0 {
		return defaults[0], nil
	}
	return nil, nil
}

// GetRedirectURI returns a property or a default.
func (c Config) GetRedirectURI(inherited ...string) string {
	return configutil.CoalesceString(c.RedirectURI, "", inherited...)
}

// GetHostedDomain returns a property or a default.
func (c Config) GetHostedDomain(inherited ...string) string {
	return configutil.CoalesceString(c.HostedDomain, "", inherited...)
}

// GetScopes gets oauth scopes to authenticate with.
func (c Config) GetScopes(inherited ...[]string) []string {
	return configutil.CoalesceStrings(c.Scopes, DefaultScopes, inherited...)
}

// GetClientID returns a property or a default.
func (c Config) GetClientID(inherited ...string) string {
	return configutil.CoalesceString(c.ClientID, "", inherited...)
}

// GetClientSecret returns a property or a default.
func (c Config) GetClientSecret(inherited ...string) string {
	return configutil.CoalesceString(c.ClientSecret, "", inherited...)
}
