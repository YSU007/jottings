package log

import (
	"fmt"
	"os"
)

type FmtLog struct{}

func (FmtLog) Info(format string, a ...interface{}) {
	fmt.Println("InfoLog:", fmt.Sprintf(format, a...))
}

func (FmtLog) Debug(format string, a ...interface{}) {
	fmt.Println("DebugLog:", fmt.Sprintf(format, a...))
}

func (FmtLog) Warn(format string, a ...interface{}) {
	fmt.Println("WarnLog:", fmt.Sprintf(format, a...))
}

func (FmtLog) Error(format string, a ...interface{}) {
	fmt.Println("ErrorLog:", fmt.Sprintf(format, a...))
}

func (FmtLog) Fatal(format string, a ...interface{}) {
	defer os.Exit(0)
	fmt.Println("FatalLog:", fmt.Sprintf(format, a...))
}
