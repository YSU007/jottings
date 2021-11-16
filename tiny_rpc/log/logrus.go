package log

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetLevel(log.TraceLevel)
}

type Logrus struct{}

func (Logrus) Info(format string, a ...interface{}) {
	log.Info(fmt.Sprintf(format, a...))
}

func (Logrus) Debug(format string, a ...interface{}) {
	log.Debug(fmt.Sprintf(format, a...))
}

func (Logrus) Warn(format string, a ...interface{}) {
	log.Warn(fmt.Sprintf(format, a...))
}

func (Logrus) Error(format string, a ...interface{}) {
	log.Error(fmt.Sprintf(format, a...))
}

func (Logrus) Fatal(format string, a ...interface{}) {
	log.Fatal(fmt.Sprintf(format, a...))
}
