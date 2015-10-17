package command

import (
	"io"
	"log"
	"os"
)

// TurnOffLogging turns off logging by redirecting it to a NOOP
func TurnOffLogging() io.Writer {
	w := &noOpWriter{}
	log.SetOutput(w)
	return w
}

// TurnOnLogging turns on logging
func TurnOnLogging() io.Writer {
	w := os.Stdout
	log.SetOutput(w)
	return w
}

// NoOpWriter is a basic type that logs to nowhere
type noOpWriter struct{}

func (w *noOpWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}
