package logging

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/Valdenirmezadri/core-go/safe"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// Logger is the actual logger which creates log records based on the functions
// called and passes them to the underlying logging backend.
type logger struct {
	backends safe.Item[[]LeveledBackend]
	// Sequence number is incremented and utilized for all log records created.
	sequenceNo uint64
	// timeNow is a customizable for testing purposes.
	timeNow func() time.Time
	// ExtraCallDepth can be used to add additional call depth when getting the
	// calling function. This is normally used when wrapping a logger.
	ExtraCalldepth int
}

// TODO call NewLogger and remove MustGetLogger?

type Logger interface {
	// Fatal is equivalent to l.Critical(fmt.Sprint()) followed by a call to os.Exit(1).
	Fatal(args ...interface{})
	// Fatalf is equivalent to l.Critical followed by a call to os.Exit(1).
	Fatalf(format string, args ...interface{})
	// Panic is equivalent to l.Critical(fmt.Sprint()) followed by a call to panic().
	Panic(args ...interface{})
	// Panicf is equivalent to l.Critical followed by a call to panic().
	Panicf(format string, args ...interface{})

	// Critical logs a message using CRITICAL as log level.
	Critical(args ...interface{})
	// Criticalf logs a message using CRITICAL as log level.
	Criticalf(format string, args ...interface{})
	// Error logs a message using ERROR as log level.
	Error(args ...interface{})

	// Errorf logs a message using ERROR as log level.
	Errorf(format string, args ...interface{})
	// Warning logs a message using WARNING as log level.
	Warning(args ...interface{})

	// Warningf logs a message using WARNING as log level.
	Warningf(format string, args ...interface{})

	Warnf(format string, args ...interface{})

	// Notice logs a message using NOTICE as log level.
	Notice(args ...interface{})

	// Noticef logs a message using NOTICE as log level.
	Noticef(format string, args ...interface{})

	// Info logs a message using INFO as log level.
	Info(args ...interface{})

	// Infof logs a message using INFO as log level.
	Infof(format string, args ...interface{})

	// Printf logs a message using INFO as log level.
	Printf(format string, args ...interface{})

	// Debug logs a message using DEBUG as log level.
	Debug(args ...interface{})

	// Debugf logs a message using DEBUG as log level.
	Debugf(format string, args ...interface{})

	Sub(module string) waLog.Logger

	SetLevel(string)
}

func New(lv string, backends ...Backend) (Logger, error) {
	if len(backends) == 0 {
		return nil, fmt.Errorf("need at least one backend")
	}

	logger := &logger{
		backends: safe.NewItem[[]LeveledBackend](),
		timeNow:  time.Now,
	}

	logger.prepareBackends(lv, backends)

	return logger, nil
}

func (l *logger) prepareBackends(lv string, backends []Backend) {
	var list []LeveledBackend
	for _, backend := range backends {
		list = append(list, l.prepareBackend(lv, backend))
	}

	l.backends.Set(list)
}

func (l logger) prepareBackend(lv string, backend Backend) (leveled LeveledBackend) {
	var level Level
	leveled, ok := backend.(LeveledBackend)
	if !ok {
		leveled = &moduleLeveled{
			level:   safe.NewItemWithData(level.New(lv)),
			backend: backend,
		}
	}

	return leveled
}

func (l *logger) SetLevel(lv string) {
	l.backends.Update(func(lb []LeveledBackend) []LeveledBackend {
		list := make([]LeveledBackend, len(lb))
		for k := range lb {
			item := lb[k]
			item.SetLevel(lv)
			list[k] = item
		}

		return list
	})
}

func (l *logger) log(lvl Level, format *string, args ...interface{}) {
	record := &Record{
		ID:    atomic.AddUint64(&l.sequenceNo, 1),
		Time:  l.timeNow(),
		Level: lvl,
		fmt:   format,
		Args:  args,
	}

	for _, backend := range l.backends.Get() {
		blv := backend.GetLevel()
		if blv >= lvl {
			backend.Log(lvl, 2+l.ExtraCalldepth, record)
		}
	}
}

// Fatal is equivalent to l.Critical(fmt.Sprint()) followed by a call to os.Exit(1).
func (l *logger) Fatal(args ...interface{}) {
	l.log(CRITICAL, nil, args...)
	os.Exit(1)
}

// Fatalf is equivalent to l.Critical followed by a call to os.Exit(1).
func (l *logger) Fatalf(format string, args ...interface{}) {
	l.log(CRITICAL, &format, args...)
	os.Exit(1)
}

// Panic is equivalent to l.Critical(fmt.Sprint()) followed by a call to panic().
func (l *logger) Panic(args ...interface{}) {
	l.log(CRITICAL, nil, args...)
	panic(fmt.Sprint(args...))
}

// Panicf is equivalent to l.Critical followed by a call to panic().
func (l *logger) Panicf(format string, args ...interface{}) {
	l.log(CRITICAL, &format, args...)
	panic(fmt.Sprintf(format, args...))
}

// Critical logs a message using CRITICAL as log level.
func (l *logger) Critical(args ...interface{}) {
	l.log(CRITICAL, nil, args...)
}

// Criticalf logs a message using CRITICAL as log level.
func (l *logger) Criticalf(format string, args ...interface{}) {
	l.log(CRITICAL, &format, args...)
}

// Error logs a message using ERROR as log level.
func (l *logger) Error(args ...interface{}) {
	l.log(ERROR, nil, args...)
}

// Errorf logs a message using ERROR as log level.
func (l *logger) Errorf(format string, args ...interface{}) {
	l.log(ERROR, &format, args...)
}

// Warning logs a message using WARNING as log level.
func (l *logger) Warning(args ...interface{}) {
	l.log(WARNING, nil, args...)
}

// Warningf logs a message using WARNING as log level.
func (l *logger) Warningf(format string, args ...interface{}) {
	l.log(WARNING, &format, args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.Warningf(format, args...)
}

// Notice logs a message using NOTICE as log level.
func (l *logger) Notice(args ...interface{}) {
	l.log(NOTICE, nil, args...)
}

// Noticef logs a message using NOTICE as log level.
func (l *logger) Noticef(format string, args ...interface{}) {
	l.log(NOTICE, &format, args...)
}

// Info logs a message using INFO as log level.
func (l *logger) Info(args ...interface{}) {
	l.log(INFO, nil, args...)
}

// Infof logs a message using INFO as log level.
func (l *logger) Infof(format string, args ...interface{}) {
	l.log(INFO, &format, args...)
}

// Printf logs a message using INFO as log level.
func (l *logger) Printf(format string, args ...interface{}) {
	l.log(INFO, &format, args...)
}

// Debug logs a message using DEBUG as log level.
func (l *logger) Debug(args ...interface{}) {
	l.log(DEBUG, nil, args...)
}

// Debugf logs a message using DEBUG as log level.
func (l *logger) Debugf(format string, args ...interface{}) {
	l.log(DEBUG, &format, args...)
}

func (l *logger) Sub(module string) waLog.Logger {
	return l
}
