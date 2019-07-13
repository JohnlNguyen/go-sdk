package jobkit

import (
	"context"

	"go-sdk/cron"
	"go-sdk/logger"
)

// Debugf prints an info message if the logger is set.
func Debugf(ctx context.Context, log logger.Log, format string, args ...interface{}) {
	if log == nil {
		return
	}
	ji := cron.GetJobInvocation(ctx)
	log.SubContext(ji.ID).Debugf(format, args...)
}

// Infof prints an info message if the logger is set.
func Infof(ctx context.Context, log logger.Log, format string, args ...interface{}) {
	if log == nil {
		return
	}
	ji := cron.GetJobInvocation(ctx)
	log.SubContext(ji.ID).Infof(format, args...)
}

// Warningf prints a warning message if the logger is set.
func Warningf(ctx context.Context, log logger.Log, format string, args ...interface{}) {
	if log == nil {
		return
	}
	ji := cron.GetJobInvocation(ctx)
	log.SubContext(ji.ID).Warningf(format, args...)
}

// Warning prints an warning if the logger is set.
func Warning(ctx context.Context, log logger.Log, err error) {
	if log == nil {
		return
	}
	ji := cron.GetJobInvocation(ctx)
	log.SubContext(ji.ID).Warning(err)
}

// Errorf prints an error message if the logger is set.
func Errorf(ctx context.Context, log logger.Log, format string, args ...interface{}) {
	if log == nil {
		return
	}
	ji := cron.GetJobInvocation(ctx)
	log.SubContext(ji.ID).Errorf(format, args...)
}

// Error prints an error if the logger is set.
func Error(ctx context.Context, log logger.Log, err error) {
	if log == nil {
		return
	}
	ji := cron.GetJobInvocation(ctx)
	log.SubContext(ji.ID).Error(err)
}

// Fatalf prints a fatal error message if the logger is set.
func Fatalf(ctx context.Context, log logger.Log, format string, args ...interface{}) {
	if log == nil {
		return
	}
	ji := cron.GetJobInvocation(ctx)
	log.SubContext(ji.ID).Fatalf(format, args...)
}

// Fatal prints a fatal error if the logger is set.
func Fatal(ctx context.Context, log logger.Log, err error) {
	if log == nil {
		return
	}
	ji := cron.GetJobInvocation(ctx)
	log.SubContext(ji.ID).Fatal(err)
}
