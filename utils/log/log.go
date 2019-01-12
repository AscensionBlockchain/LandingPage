package log

import (
	"fmt"
	"os"
	"strings"

	raven "github.com/getsentry/raven-go"
	"github.com/tchap/zapext/zapsentry"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log *zap.SugaredLogger
	// Debug, Info, Error, Fatal func(string, ...interface{})
)

func init() {
	EnableLogging()
}

func EnableLogging() {
	initZapLogger(
		os.Getenv("SENTRY_SINK") != "",
	)
}

func DisableLogging() {
	noLogs := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return false
	})
	config := zap.NewDevelopmentEncoderConfig()
	jsonEncoder := zapcore.NewJSONEncoder(config)
	consoleOutput := zapcore.Lock(os.Stderr)
	core := zapcore.NewCore(jsonEncoder, consoleOutput, noLogs)
	Log = zap.New(core).Sugar()
}

func initZapLogger(withSentry bool) {
	//	logger, err := zap.NewDevelopment()
	//	if err != nil {
	//		panic(err)
	//	}

	// Log filter
	allLogs := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		// return lvl < zapcore.ErrorLevel
		return true
	})

	// Log output encoding
	// textEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	config := zap.NewDevelopmentEncoderConfig()
	jsonEncoder := zapcore.NewJSONEncoder(config)

	// Log sinks

	//   Basic STDERR output sink
	consoleOutput := zapcore.Lock(os.Stderr)

	var core zapcore.Core

	if withSentry {
		dsn := os.Getenv("SENTRY_SINK")
		if dsn == "" {
			panic("requested Sentry logging, but SENTRY_SINK is not set in env")
		}

		dsn = strings.Trim(dsn, "'")

		fmt.Printf("Sending Zap logs to Sentry at: %s\n", dsn)

		client, err := raven.New(dsn)
		if err != nil {
			panic(err)
		}

		core = zapcore.NewTee(
			zapcore.NewCore(jsonEncoder, consoleOutput, allLogs),
			zapsentry.NewCore(allLogs, client),
		)
	} else {
		core = zapcore.NewCore(jsonEncoder, consoleOutput, allLogs)
	}

	Log = zap.New(core).Sugar()
}
