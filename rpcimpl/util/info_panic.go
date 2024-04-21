package util

import (
	"fmt"

	"tiny_rpc/log"
)

func InfoPanic(format string, v ...interface{}) {
	if err := recover(); err != nil {
		log.Error("panic %v info %s", err, fmt.Sprintf(format, v...))
	}
}
