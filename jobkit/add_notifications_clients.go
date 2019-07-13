package jobkit

import (
	"go-sdk/airbrake"
	"go-sdk/aws"
	"go-sdk/aws/ses"
	"go-sdk/datadog"
	"go-sdk/diagnostics"
	"go-sdk/email"
	"go-sdk/slack"
	"go-sdk/stats"
)

// AddNotificationClients adds notification clients to a given job.
func AddNotificationClients(job *Job, cfg *Config) error {
	var err error
	// set up myriad of notification targets
	var emailClient email.Sender
	if !cfg.AWS.IsZero() {
		emailClient = ses.New(aws.MustNewSession(&cfg.AWS))
	}
	var slackClient slack.Sender
	if !cfg.Slack.IsZero() {
		slackClient = slack.New(&cfg.Slack)
	}
	var statsClient stats.Collector
	if !cfg.Datadog.IsZero() {
		statsClient, err = datadog.NewCollector(&cfg.Datadog)
		if err != nil {
			return err
		}
	}
	var errorClient diagnostics.Notifier
	if !cfg.Airbrake.IsZero() {
		errorClient = airbrake.MustNew(&cfg.Airbrake)
	}

	job.WithEmailClient(emailClient).
		WithStatsClient(statsClient).
		WithSlackClient(slackClient).
		WithErrorClient(errorClient)

	return nil
}
