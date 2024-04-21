package log

import (
	"path"
	"runtime"

	log "github.com/sirupsen/logrus"
)

func retLogrus() Interface {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File)
			return frame.Function, fileName
		},
	})
	log.SetLevel(log.TraceLevel)

	return log.StandardLogger()
}
