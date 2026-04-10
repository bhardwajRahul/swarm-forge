package swarmlog

import (
	"fmt"
	"io"
)

// Logger writes timestamped log entries to one or more writers.
type Logger struct {
	writers []io.Writer
}

// New creates a Logger that writes to all supplied writers.
func New(writers ...io.Writer) *Logger {
	return &Logger{writers: writers}
}

// Write formats and writes a log entry with the given role and message.
func (l *Logger) Write(role, message string) error {
	entry := fmt.Sprintf("[%s] %s\n", role, message)
	for _, w := range l.writers {
		if _, err := fmt.Fprint(w, entry); err != nil {
			return err
		}
	}
	return nil
}
