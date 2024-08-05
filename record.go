package logging

import (
	"bytes"
	"fmt"
	"time"
)

// Redactor is an interface for types that may contain sensitive information
// (like passwords), which shouldn't be printed to the log. The idea was found
// in relog as part of the vitness project.
type Redactor interface {
	Redacted() interface{}
}

// Record represents a log record and contains the timestamp when the record
// was created, an increasing id, filename and line and finally the actual
// formatted log line.
type Record struct {
	ID     uint64
	Time   time.Time
	Module string
	Level  Level
	Args   []interface{}

	// message is kept as a pointer to have shallow copies update this once
	// needed.
	message   *string
	fmt       *string
	formatter Formatter
	formatted string
}

// Formatted returns the formatted log record string.
func (r *Record) Formatted(calldepth int) string {
	if r.formatted == "" {
		var buf bytes.Buffer
		r.formatter.Format(calldepth+1, r, &buf)
		r.formatted = buf.String()
	}
	return r.formatted
}

// Message returns the log record message.
func (r *Record) Message() string {
	if r.message == nil {
		// Redact the arguments that implements the Redactor interface
		for i, arg := range r.Args {
			if redactor, ok := arg.(Redactor); ok {
				r.Args[i] = redactor.Redacted()
			}
		}
		var buf bytes.Buffer
		if r.fmt != nil {
			fmt.Fprintf(&buf, *r.fmt, r.Args...)
		} else {
			// use Fprintln to make sure we always get space between arguments
			fmt.Fprintln(&buf, r.Args...)
			buf.Truncate(buf.Len() - 1) // strip newline
		}
		msg := buf.String()
		r.message = &msg
	}
	return *r.message
}
