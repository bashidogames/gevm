package logging

import (
	"fmt"

	"github.com/bashmills/gevm/internal/utils"
	"github.com/bashmills/gevm/logger"
)

type Level int64

const (
	NOTHING Level = iota
	ERROR
	WARNING
	INFO
	DEBUG
	TRACE
)

type Logging struct {
	Level Level
}

func (l Logging) Errorf(format string, a ...any) {
	if !l.ShouldLog(ERROR) {
		return
	}

	l.Error(fmt.Sprintf(format, a...))
}

func (l Logging) Error(msg string) {
	if !l.ShouldLog(ERROR) {
		return
	}

	utils.Printlnf("ERROR: %s", msg)
}

func (l Logging) Warningf(format string, a ...any) {
	if !l.ShouldLog(WARNING) {
		return
	}

	l.Warning(fmt.Sprintf(format, a...))
}

func (l Logging) Warning(msg string) {
	if !l.ShouldLog(WARNING) {
		return
	}

	utils.Printlnf("WARNING: %s", msg)
}

func (l Logging) Infof(format string, a ...any) {
	if !l.ShouldLog(INFO) {
		return
	}

	l.Info(fmt.Sprintf(format, a...))
}

func (l Logging) Info(msg string) {
	if !l.ShouldLog(INFO) {
		return
	}

	utils.Println(msg)
}

func (l Logging) Debugf(format string, a ...any) {
	if !l.ShouldLog(DEBUG) {
		return
	}

	l.Debug(fmt.Sprintf(format, a...))
}

func (l Logging) Debug(msg string) {
	if !l.ShouldLog(DEBUG) {
		return
	}

	utils.Println(msg)
}

func (l Logging) Tracef(format string, a ...any) {
	if !l.ShouldLog(TRACE) {
		return
	}

	l.Trace(fmt.Sprintf(format, a...))
}

func (l Logging) Trace(msg string) {
	if !l.ShouldLog(TRACE) {
		return
	}

	utils.Println(msg)
}

func (l Logging) ShouldLog(level Level) bool {
	return l.Level >= level
}

func New(level Level) (logger.Logger, error) {
	return &Logging{
		Level: level,
	}, nil
}
