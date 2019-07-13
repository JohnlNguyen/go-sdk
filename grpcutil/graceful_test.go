package grpcutil

import "go-sdk/graceful"

// Validate the interface is satisfied.
var (
	_ (graceful.Graceful) = (*Graceful)(nil)
)
