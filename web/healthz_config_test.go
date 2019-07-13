package web

import (
	"testing"
	"time"

	"go-sdk/assert"
	"go-sdk/yaml"
)

func TestHealthzConfigYAML(t *testing.T) {
	assert := assert.New(t)

	yml := `
bindAddr: ":4444"
gracePeriod: "10s"
recoverPanics: false
`
	var cfg HealthzConfig
	assert.Nil(yaml.Unmarshal([]byte(yml), &cfg))

	assert.Equal(":4444", cfg.GetBindAddr())
	assert.Equal(10*time.Second, cfg.GetGracePeriod())
	assert.Equal(false, cfg.GetRecoverPanics())
}
