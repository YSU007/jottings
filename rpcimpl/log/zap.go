package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const logSkip = 1

var (
	outDir    = "./"
	outPath   = []string{"stdout", outDir + "log.log"}
	errorPath = []string{"stderr", outDir + "err.log"}
)

func retZap() Interface {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.Encoding = "json"
	config.OutputPaths = outPath
	config.ErrorOutputPaths = errorPath
	ZapLog, err := config.Build(zap.AddCallerSkip(logSkip))
	if err != nil {
		panic(fmt.Sprintf("zap log: %v", err))
	}

	return ZapLog.Sugar()
}
