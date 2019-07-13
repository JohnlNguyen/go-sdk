package grpcutil

import (
	"net"

	"go-sdk/async"
	"go-sdk/logger"
	"google.golang.org/grpc"
)

// NewGraceful returns a new graceful host for a grpc server.
func NewGraceful(listener net.Listener, server *grpc.Server) *Graceful {
	return &Graceful{
		Listener: listener,
		Server:   server,
	}
}

// Graceful is a shim for graceful hosting grpc servers.
type Graceful struct {
	Log      logger.FullReceiver
	Latch    async.Latch
	Listener net.Listener
	Server   *grpc.Server
}

// WithLogger sets the logger.
func (gz *Graceful) WithLogger(log logger.FullReceiver) *Graceful {
	gz.Log = log
	return gz
}

// Start starts the server.
func (gz *Graceful) Start() error {
	gz.Latch.Starting()
	gz.Latch.Started()
	logger.MaybeSyncInfof(gz.Log, "grpc server starting, listening on %v %s", gz.Listener.Addr().Network(), gz.Listener.Addr().String())
	return gz.Server.Serve(gz.Listener)
}

// Stop shuts the server down.
func (gz *Graceful) Stop() error {
	gz.Latch.Stopping()
	logger.MaybeSyncInfof(gz.Log, "grpc server shutting down")
	gz.Server.GracefulStop()
	gz.Latch.Stopped()
	return nil
}

// IsRunning returns if the server is running.
func (gz *Graceful) IsRunning() bool {
	return gz.Latch.IsRunning()
}

// NotifyStarted returns the notify started signal.
func (gz *Graceful) NotifyStarted() <-chan struct{} {
	return gz.Latch.NotifyStarted()
}

// NotifyStopped returns the notify stopped signal.
func (gz *Graceful) NotifyStopped() <-chan struct{} {
	return gz.Latch.NotifyStopped()
}
