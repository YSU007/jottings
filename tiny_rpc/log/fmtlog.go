package log

import "fmt"

type FmtLog struct{}

func (FmtLog) Info(format string, a ...interface{}) {
	fmt.Println("InfoLog: ", fmt.Sprintf(format, a...))
}

func (FmtLog) Debug(format string, a ...interface{}) {
	fmt.Println("DebugLog: ", fmt.Sprintf(format, a...))
}

func (FmtLog) Warn(format string, a ...interface{}) {
	fmt.Println("WarnLog: ", fmt.Sprintf(format, a...))
}

func (FmtLog) Error(format string, a ...interface{}) {
	fmt.Println("ErrorLog: ", fmt.Sprintf(format, a...))
}
