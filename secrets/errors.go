package secrets

import "go-sdk/exception"

// Common error codes.
const (
	ErrNotFound     exception.Class = "secrets; not found"
	ErrUnauthorized exception.Class = "secrets; not authorized"
	ErrServerError  exception.Class = "secrets; remote error"
)
