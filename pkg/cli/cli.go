// Package cli provides standardized exit codes for command-line tools.
// Use with pkg/errors.ExitCode(err) to map errors to exit codes.
package cli

// Exit codes standardized (compatible with pkg/errors.ExitCode).
const (
	ExitOK           = 0
	ExitErr          = 1
	ExitInvalidInput = 2
	ExitNotFound     = 3
	ExitUnauthorized = 4
	ExitConflict     = 5
	ExitUnavailable  = 6
)
