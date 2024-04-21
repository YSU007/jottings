package log

import (
	"fmt"
	"os"
)

type FmtLog struct{}

func (l FmtLog) Debug(v ...interface{}) {
	fmt.Println(v...)
}

func (l FmtLog) Debugf(format string, v ...interface{}) {
	fmt.Println("DebugLog:", fmt.Sprintf(format, v...))
}

func (l FmtLog) Info(v ...interface{}) {
	fmt.Println(v...)
}

func (l FmtLog) Infof(format string, v ...interface{}) {
	fmt.Println("InfoLog:", fmt.Sprintf(format, v...))
}

func (l FmtLog) Warn(v ...interface{}) {
	fmt.Println(v...)
}

func (l FmtLog) Warnf(format string, v ...interface{}) {
	fmt.Println("WarnLog:", fmt.Sprintf(format, v...))
}

func (l FmtLog) Error(v ...interface{}) {
	fmt.Println(v...)
}

func (l FmtLog) Errorf(format string, v ...interface{}) {
	fmt.Println("ErrorLog:", fmt.Sprintf(format, v...))
}

func (l FmtLog) Fatal(v ...interface{}) {
	fmt.Println(v...)
}

func (l FmtLog) Fatalf(format string, v ...interface{}) {
	defer os.Exit(1)
	fmt.Println("FatalLog:", fmt.Sprintf(format, v...))
}

func (l FmtLog) Panic(v ...interface{}) {
	panic(fmt.Sprint(v...))
}

func (l FmtLog) Panicf(format string, v ...interface{}) {
	panic(fmt.Sprintf(format, v...))
}
