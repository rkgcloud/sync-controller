package internal

import "carvel.dev/imgpkg/pkg/imgpkg/plainimage"

// imgpkg logger

var _ plainimage.Logger = &NoopLogger{}

// NewNoopLogger creates a new noop logger
func NewNoopLogger() *NoopLogger {
	return &NoopLogger{}
}

// NoopLogger this logger will not print
type NoopLogger struct{}

// Logf does nothing
func (n NoopLogger) Logf(string, ...interface{}) {}
