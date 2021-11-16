package log

// Interface ----------------------------------------------------------------------------------------------------
type Interface interface {
	Info(format string, a ...interface{})
	Debug(format string, a ...interface{})
	Warn(format string, a ...interface{})
	Error(format string, a ...interface{})
	Fatal(format string, a ...interface{})
}

var defLog Interface

func SetDefLog(log Interface) {
	defLog = log
}

func Info(format string, a ...interface{}) {
	defLog.Info(format, a...)
}

func Debug(format string, a ...interface{}) {
	defLog.Debug(format, a...)
}

func Warn(format string, a ...interface{}) {
	defLog.Warn(format, a...)
}

func Error(format string, a ...interface{}) {
	defLog.Error(format, a...)
}
