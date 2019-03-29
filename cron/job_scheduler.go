package cron

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
)

// NewJobScheduler returns a job scheduler for a given job.
func NewJobScheduler(job Job, options ...JobSchedulerOption) *JobScheduler {
	js := &JobScheduler{
		Latch: async.NewLatch(),
		Name:  job.Name(),
		Job:   job,
	}

	if typed, ok := job.(DescriptionProvider); ok {
		js.Description = typed.Description()
	}

	if typed, ok := job.(ScheduleProvider); ok {
		js.Schedule = typed.Schedule()
	}

	if typed, ok := job.(TimeoutProvider); ok {
		js.TimeoutProvider = typed.Timeout
	} else {
		js.TimeoutProvider = func() time.Duration { return 0 }
	}

	if typed, ok := job.(EnabledProvider); ok {
		js.EnabledProvider = typed.Enabled
	} else {
		js.EnabledProvider = func() bool { return DefaultEnabled }
	}

	if typed, ok := job.(SerialProvider); ok {
		js.SerialProvider = typed.Serial
	} else {
		js.SerialProvider = func() bool { return DefaultSerial }
	}

	if typed, ok := job.(ShouldTriggerListenersProvider); ok {
		js.ShouldTriggerListenersProvider = typed.ShouldTriggerListeners
	} else {
		js.ShouldTriggerListenersProvider = func() bool { return DefaultShouldTriggerListeners }
	}

	if typed, ok := job.(ShouldWriteOutputProvider); ok {
		js.ShouldWriteOutputProvider = typed.ShouldWriteOutput
	} else {
		js.ShouldWriteOutputProvider = func() bool { return DefaultShouldWriteOutput }
	}

	for _, option := range options {
		option(js)
	}

	return js
}

// JobScheduler is a job instance.
type JobScheduler struct {
	sync.Mutex   `json:"-"`
	*async.Latch `json:"-"`

	Name        string `json:"name"`
	Description string `json:"description"`
	Job         Job    `json:"-"`

	Tracer        Tracer        `json:"-"`
	Log           logger.Log    `json:"-"`
	HistoryConfig HistoryConfig `json:"-"`

	// Meta Fields
	Disabled    bool            `json:"disabled"`
	NextRuntime time.Time       `json:"nextRuntime"`
	Current     *JobInvocation  `json:"current"`
	Last        *JobInvocation  `json:"last"`
	History     []JobInvocation `json:"history"`

	Schedule                       Schedule             `json:"-"`
	EnabledProvider                func() bool          `json:"-"`
	SerialProvider                 func() bool          `json:"-"`
	TimeoutProvider                func() time.Duration `json:"-"`
	ShouldTriggerListenersProvider func() bool          `json:"-"`
	ShouldWriteOutputProvider      func() bool          `json:"-"`
}

// Start starts the scheduler.
// This call blocks.
func (js *JobScheduler) Start() error {
	if !js.Latch.CanStart() {
		return fmt.Errorf("already started")
	}
	js.Latch.Starting()
	js.RunLoop()
	return nil
}

// StartAsync starts the job scheduler in the background.
func (js *JobScheduler) StartAsync() error {
	if !js.Latch.CanStart() {
		return fmt.Errorf("already started")
	}
	js.Latch.Starting()
	go js.RunLoop()
	<-js.Latch.NotifyStarted()
	return nil
}

// Stop stops the scheduler.
func (js *JobScheduler) Stop() error {
	if !js.Latch.CanStop() {
		return fmt.Errorf("already stopped")
	}
	js.Latch.Stopping()
	<-js.Latch.NotifyStopped()
	return nil
}

// NotifyStarted notifies the job scheduler has started.
func (js *JobScheduler) NotifyStarted() <-chan struct{} {
	return js.Latch.NotifyStarted()
}

// NotifyStopped notifies the job scheduler has stopped.
func (js *JobScheduler) NotifyStopped() <-chan struct{} {
	return js.Latch.NotifyStopped()
}

// Enable sets the job as enabled.
func (js *JobScheduler) Enable() {
	js.Lock()
	defer js.Unlock()

	js.Disabled = false
	if js.Log != nil && js.ShouldTriggerListenersProvider() {
		event := NewEvent(FlagEnabled, js.Name, OptEventWritable(js.ShouldWriteOutputProvider()))
		js.Log.Trigger(context.Background(), event)
	}
	if typed, ok := js.Job.(OnEnabledReceiver); ok {
		typed.OnEnabled(context.Background())
	}
}

// Disable sets the job as disabled.
func (js *JobScheduler) Disable() {
	js.Lock()
	defer js.Unlock()

	js.Disabled = true
	if js.Log != nil && js.ShouldTriggerListenersProvider() {
		event := NewEvent(FlagDisabled, js.Name, OptEventWritable(js.ShouldWriteOutputProvider()))
		js.Log.Trigger(context.Background(), event)
	}
	if typed, ok := js.Job.(OnDisabledReceiver); ok {
		typed.OnDisabled(context.Background())
	}
}

// Cancel stops an execution in process.
func (js *JobScheduler) Cancel() {
	if js.Current != nil {
		js.Current.Cancel()
	}
}

// RunLoop is the main scheduler loop.
// it alarms on the next runtime and forks a new routine to run the job.
// It can be aborted with the scheduler's async.Latch.
func (js *JobScheduler) RunLoop() {
	js.Latch.Started()

	if js.Schedule != nil {
		// sniff the schedule, see if a next runtime is called for (or if the job is on demand).
		js.NextRuntime = js.Schedule.Next(js.NextRuntime)
	}
	if js.NextRuntime.IsZero() {
		js.Latch.Stopped()
		return
	}

	for {
		if js.NextRuntime.IsZero() {
			js.Latch.Stopped()
			return
		}
		runAt := time.After(js.NextRuntime.UTC().Sub(Now()))
		select {
		case <-runAt:
			// start the job
			go js.Run()
			// set up the next runtime.
			js.NextRuntime = js.Schedule.Next(js.NextRuntime)
		case <-js.Latch.NotifyStopping():
			js.Latch.Stopped()
			return
		}
	}
}

// Run forces the job to run.
// It checks if the job should be allowed to execute.
// It blocks on the job execution to enforce or clear timeouts.
func (js *JobScheduler) Run() {
	// check if the job can run
	if !js.canRun() {
		return
	}

	// mark the start time
	start := Now()

	timeout := js.TimeoutProvider()

	// create the root context.
	ctx, cancel := js.createContextWithTimeout(timeout)

	// create a job invocation, or a record of each
	// individual execution of a job.
	ji := JobInvocation{
		ID:      NewJobInvocationID(),
		JobName: js.Name,
		Status:  JobStatusRunning,
		Started: start,
		Context: ctx,
		Cancel:  cancel,
	}
	if timeout > 0 {
		ji.Timeout = start.Add(timeout)
	}
	js.setCurrent(&ji)

	var err error
	var tf TraceFinisher
	// load the job invocation into the context
	ctx = WithJobInvocation(ctx, &ji)

	// this defer runs all cleanup actions
	// it recovers panics
	// it cancels the timeout (if relevant)
	// it rotates the current and last references
	// it fires lifecycle events
	defer func() {
		if r := recover(); r != nil {
			err = exception.New(err)
		}
		cancel()
		if tf != nil {
			tf.Finish(ctx)
		}

		ji.Finished = Now()
		ji.Elapsed = ji.Finished.Sub(ji.Started)
		ji.Err = err

		if err != nil && IsJobCancelled(err) {
			ji.Cancelled = ji.Finished
			js.onCancelled(ctx, &ji)
		} else if ji.Err != nil {
			js.onFailure(ctx, &ji)
		} else {
			js.onComplete(ctx, &ji)
		}

		js.addHistory(ji)
		js.setCurrent(nil)
		js.setLast(&ji)
	}()

	// if the tracer is set, create a trace context
	if js.Tracer != nil {
		ctx, tf = js.Tracer.Start(ctx)
	}
	// fire the on start event
	js.onStart(ctx, &ji)

	// check if the job has been canceled
	// or if it's finished.
	select {
	case <-ctx.Done():
		err = ErrJobCancelled
	case err = <-js.safeAsyncExec(ctx):
	}
}

//
// exported utility methods
//

// GetInvocationByID returns an invocation by id.
func (js *JobScheduler) GetInvocationByID(id string) *JobInvocation {
	for _, ji := range js.History {
		if ji.ID == id {
			return &ji
		}
	}
	return nil
}

//
// utility functions
//

func (js *JobScheduler) setCurrent(ji *JobInvocation) {
	js.Lock()
	js.Current = ji
	js.Unlock()
}

func (js *JobScheduler) setLast(ji *JobInvocation) {
	js.Lock()
	js.Last = ji
	js.Unlock()
}

// safeAsyncExec runs a given job's body and recovers panics.
func (js *JobScheduler) safeAsyncExec(ctx context.Context) chan error {
	errors := make(chan error)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errors <- exception.New(r)
			}
		}()
		errors <- js.Job.Execute(ctx)
	}()
	return errors
}

func (js *JobScheduler) createContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout > 0 {
		return context.WithTimeout(context.Background(), timeout)
	}
	return context.WithCancel(context.Background())
}

// canRun returns if a job can execute.
func (js *JobScheduler) canRun() bool {
	js.Lock()
	defer js.Unlock()

	if js.Disabled {
		return false
	}

	if js.EnabledProvider != nil {
		if !js.EnabledProvider() {
			return false
		}
	}

	if js.SerialProvider != nil && js.SerialProvider() {
		if js.Current != nil {
			return false
		}
	}
	return true
}

func (js *JobScheduler) onStart(ctx context.Context, ji *JobInvocation) {
	if js.Log != nil && js.ShouldTriggerListenersProvider() {
		event := NewEvent(FlagStarted, ji.JobName, OptEventJobInvocation(ji.ID), OptEventWritable(js.ShouldWriteOutputProvider()))
		js.Log.Trigger(ctx, event)
	}
	if typed, ok := js.Job.(OnStartReceiver); ok {
		typed.OnStart(ctx)
	}
}

func (js *JobScheduler) onCancelled(ctx context.Context, ji *JobInvocation) {
	ji.Status = JobStatusCancelled

	if js.Log != nil && js.ShouldTriggerListenersProvider() {
		event := NewEvent(FlagCancelled, ji.JobName, OptEventJobInvocation(ji.ID), OptEventElapsed(ji.Elapsed), OptEventWritable(js.ShouldWriteOutputProvider()))
		js.Log.Trigger(ctx, event)
	}
	if typed, ok := js.Job.(OnCancellationReceiver); ok {
		typed.OnCancellation(ctx)
	}
}

func (js *JobScheduler) onComplete(ctx context.Context, ji *JobInvocation) {
	ji.Status = JobStatusComplete

	if js.Log != nil && js.ShouldTriggerListenersProvider() {
		event := NewEvent(FlagComplete, ji.JobName, OptEventJobInvocation(ji.ID), OptEventElapsed(ji.Elapsed), OptEventWritable(js.ShouldWriteOutputProvider()))
		js.Log.Trigger(ctx, event)
	}
	if typed, ok := js.Job.(OnCompleteReceiver); ok {
		typed.OnComplete(ctx)
	}

	if js.Last != nil && js.Last.Err != nil {
		if js.Log != nil {
			event := NewEvent(FlagFixed, ji.JobName, OptEventElapsed(ji.Elapsed), OptEventWritable(js.ShouldWriteOutputProvider()))
			js.Log.Trigger(ctx, event)
		}

		if typed, ok := js.Job.(OnFixedReceiver); ok {
			typed.OnFixed(ctx)
		}
	}
}

func (js *JobScheduler) onFailure(ctx context.Context, ji *JobInvocation) {
	ji.Status = JobStatusFailed

	if js.Log != nil && js.ShouldTriggerListenersProvider() {
		event := NewEvent(FlagFailed, ji.JobName, OptEventErr(ji.Err), OptEventJobInvocation(ji.ID), OptEventElapsed(ji.Elapsed), OptEventWritable(js.ShouldWriteOutputProvider()))

		js.Log.Trigger(ctx, event)
	}
	if ji.Err != nil {
		logger.MaybeError(js.Log, ji.Err)
	}
	if typed, ok := js.Job.(OnFailureReceiver); ok {
		typed.OnFailure(ctx)
	}
	if js.Last != nil && js.Last.Err == nil {
		if js.Log != nil {
			event := NewEvent(FlagBroken, ji.JobName, OptEventJobInvocation(ji.ID), OptEventElapsed(ji.Elapsed), OptEventWritable(js.ShouldWriteOutputProvider()))
			js.Log.Trigger(ctx, event)
		}

		if typed, ok := js.Job.(OnBrokenReceiver); ok {
			typed.OnBroken(ctx)
		}
	}
}

func (js *JobScheduler) addHistory(ji JobInvocation) {
	js.Lock()
	defer js.Unlock()
	js.History = append(js.cullHistory(), ji)
}

func (js *JobScheduler) cullHistory() []JobInvocation {
	count := len(js.History)
	maxCount := js.HistoryConfig.MaxCountOrDefault()
	maxAge := js.HistoryConfig.MaxAgeOrDefault()
	now := time.Now().UTC()
	var filtered []JobInvocation
	for index, h := range js.History {
		if maxCount > 0 {
			if index < (count - maxCount) {
				continue
			}
		}
		if maxAge > 0 {
			if now.Sub(h.Started) > maxAge {
				continue
			}
		}
		filtered = append(filtered, h)
	}
	return filtered
}
