package runutil

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

// IgnoreErr prints function error to stderr, if it happens.
func IgnoreErr(f func() error) {
	err := f()
	if err == nil {
		return
	}

	_, _ = fmt.Fprintf(os.Stderr, "Ignore error: %v\n", err)
}

// MustNoErr function must return no error.
func MustNoErr(f func() error) {
	err := f()
	if err == nil {
		return
	}

	_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	panic(err)
}

// LogOnError put error into zap.log, if it happens.
func LogOnError(f func() error, log *zap.Logger, msg string, fields ...zap.Field) {
	err := f()
	if err == nil {
		return
	}

	log.With(zap.Error(err)).Error(msg, fields...)
}
