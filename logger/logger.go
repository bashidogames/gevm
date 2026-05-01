package logger

type Logger interface {
	Errorf(format string, a ...any)
	Warningf(format string, a ...any)
	Infof(format string, a ...any)
	Debugf(format string, a ...any)
	Tracef(format string, a ...any)

	Error(msg string)
	Warning(msg string)
	Info(msg string)
	Debug(msg string)
	Trace(msg string)
}
