package logging

import (
	"errors"
	"strings"
	"sync"

	"github.com/Valdenirmezadri/core-go/safe"
)

// ErrInvalidLogLevel is used when an invalid log level has been used.
var ErrInvalidLogLevel = errors.New("logger: invalid log level")

// Level defines all available log levels for log messages.
type Level int

// Log levels.
const (
	CRITICAL Level = iota
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
)

var levelNames = []string{
	"CRITICAL",
	"ERROR",
	"WARNING",
	"NOTICE",
	"INFO",
	"DEBUG",
}

func (p Level) New(name string) Level {
	name = strings.ToUpper(strings.TrimSpace(name))

	if name == CRITICAL.String() {
		return CRITICAL
	}

	if name == ERROR.String() {
		return ERROR
	}

	if name == WARNING.String() {
		return WARNING
	}

	if name == NOTICE.String() {
		return NOTICE
	}

	if name == INFO.String() {
		return INFO
	}

	return DEBUG

}

// String returns the string representation of a logging level.
func (p Level) String() string {
	return levelNames[p]
}

// LogLevel returns the log level from a string representation.
func LogLevel(level string) (Level, error) {
	for i, name := range levelNames {
		if strings.EqualFold(name, level) {
			return Level(i), nil
		}
	}
	return ERROR, ErrInvalidLogLevel
}

// Leveled interface is the interface required to be able to add leveled
// logging.
type Leveled interface {
	GetLevel() Level
	SetLevel(string)
}

// LeveledBackend is a log backend with additional knobs for setting levels on
// individual modules to different levels.
type LeveledBackend interface {
	Backend
	Leveled
}

type moduleLeveled struct {
	level     safe.Item[Level]
	backend   Backend
	formatter Formatter
	once      sync.Once
}

// GetLevel returns the log level
func (l *moduleLeveled) GetLevel() Level {
	return l.level.Get()

}

// SetLevel sets the log level
func (l *moduleLeveled) SetLevel(lv string) {
	var level Level
	l.level.Set(level.New(lv))
}

func (l *moduleLeveled) Log(level Level, calldepth int, rec *Record) (err error) {
	rec.formatter = l.getFormatterAndCacheCurrent()
	return l.backend.Log(level, calldepth+1, rec)
}

func (l *moduleLeveled) getFormatterAndCacheCurrent() Formatter {
	l.once.Do(func() {
		if l.formatter == nil {
			l.formatter = getFormatter()
		}
	})
	return l.formatter
}
