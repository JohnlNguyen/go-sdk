package jobkit

import (
	"go-sdk/airbrake"
	"go-sdk/aws"
	"go-sdk/configutil"
	"go-sdk/cron"
	"go-sdk/datadog"
	"go-sdk/email"
	"go-sdk/logger"
	"go-sdk/slack"
	"go-sdk/web"
)

// Config is the jobkit config.
type Config struct {
	cron.Config `json:",inline" yaml:",inline"`

	MaxLogBytes int `json:"maxLogBytes" yaml:"maxLogBytes"`

	Logger logger.Config `json:"logger" yaml:"logger"`
	Web    web.Config    `json:"web" yaml:"web"`

	Airbrake airbrake.Config `json:"airbrake" yaml:"airbrake"`
	AWS      aws.Config      `json:"aws" yaml:"aws"`
	Email    email.Message   `json:"email" yaml:"email"`
	Datadog  datadog.Config  `json:"datadog" yaml:"datadog"`
	Slack    slack.Config    `json:"slack" yaml:"slack"`
}

// MaxLogBytesOrDefault is a the maximum amount of log data to buffer.
func (c Config) MaxLogBytesOrDefault() int {
	return configutil.CoalesceInt(c.MaxLogBytes, DefaultMaxLogBytes)
}
