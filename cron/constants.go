package cron

import (
	"time"

	"go-sdk/logger"
)

const (
	// EnvVarHeartbeatInterval is an environment variable name.
	EnvVarHeartbeatInterval = "CRON_HEARTBEAT_INTERVAL"
)

// Retention defaults
const (
	DefaultMaxCount = 10
	DefaultMaxAge   = 6 * time.Hour
)

const (
	// DefaultHeartbeatInterval is the interval between schedule next run checks.
	DefaultHeartbeatInterval = 50 * time.Millisecond
)

const (
	// DefaultEnabled is a default.
	DefaultEnabled = true
	// DefaultSerial is a default.
	DefaultSerial = false
	// DefaultShouldWriteOutput is a default.
	DefaultShouldWriteOutput = true
	// DefaultShouldTriggerListeners is a default.
	DefaultShouldTriggerListeners = true
)

const (
	// FlagStarted is an event flag.
	FlagStarted logger.Flag = "cron.started"
	// FlagFailed is an event flag.
	FlagFailed logger.Flag = "cron.failed"
	// FlagCancelled is an event flag.
	FlagCancelled logger.Flag = "cron.cancelled"
	// FlagComplete is an event flag.
	FlagComplete logger.Flag = "cron.complete"
	// FlagBroken is an event flag.
	FlagBroken logger.Flag = "cron.broken"
	// FlagFixed is an event flag.
	FlagFixed logger.Flag = "cron.fixed"
	// FlagEnabled is an event flag.
	FlagEnabled logger.Flag = "cron.enabled"
	// FlagDisabled is an event flag.
	FlagDisabled logger.Flag = "cron.disabled"
)

// State is a job state.
type State string

const (
	//StateRunning is the running state.
	StateRunning State = "running"
	// StateEnabled is the enabled state.
	StateEnabled State = "enabled"
	// StateDisabled is the disabled state.
	StateDisabled State = "disabled"
)

// JobStatus is a job status.
type JobStatus string

// Status values.
const (
	JobStatusRunning   JobStatus = "running"
	JobStatusCancelled JobStatus = "cancelled"
	JobStatusFailed    JobStatus = "failed"
	JobStatusComplete  JobStatus = "complete"
)
