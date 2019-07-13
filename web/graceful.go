package web

import "go-sdk/graceful"

// GracefulShutdown shuts an app down gracefull.
// It is an alias to graceful.Shutdown
var GracefulShutdown = graceful.Shutdown

// StartWithGracefulShutdown shuts an app down gracefull.
// It is an alias to graceful.Shutdown
var StartWithGracefulShutdown = graceful.Shutdown
