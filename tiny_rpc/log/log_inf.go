package log

// Interface ----------------------------------------------------------------------------------------------------
type Interface interface {
	Info(format string, a ...interface{})
	Debug(format string, a ...interface{})
	Warn(format string, a ...interface{})
	Error(format string, a ...interface{})
}

var DefLog Interface

func SetDefLog(log Interface) {
	DefLog = log
}

func Info(format string, a ...interface{}) {
	DefLog.Info(format, a...)
}

func Debug(format string, a ...interface{}) {
	DefLog.Debug(format, a...)
}

func Warn(format string, a ...interface{}) {
	DefLog.Warn(format, a...)
}

func Error(format string, a ...interface{}) {
	DefLog.Error(format, a...)
}
