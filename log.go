package KolonseWeb

import (
	"github.com/kolonse/logs"
	"strings"
)

// Log levels to control the logging output.
const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

// SetLogLevel sets the global log level used by the simple
// logger.
func SetLevel(l int) {
	BeeLogger.SetLevel(l)
}

func SetLogFuncCall(b bool) {
	BeeLogger.EnableFuncCallDepth(b)
	BeeLogger.SetLogFuncCallDepth(3)
}

func DefaultLogs() *logs.BeeLogger {
	ret := logs.NewLogger(10000)
	err := ret.SetLogger("console", "")
	if err != nil {
		panic(err)
	}
	ret.EnableFuncCallDepth(true)
	ret.SetLogFuncCallDepth(3)
	return ret
}

// logger references the used application logger.
var BeeLogger *logs.BeeLogger

// SetLogger sets a new logger.
func SetLogger(adaptername string, config string) error {
	err := BeeLogger.SetLogger(adaptername, config)
	if err != nil {
		return err
	}
	return nil
}

func Emergency(v ...interface{}) {
	BeeLogger.Emergency(generateFmtStr(len(v)), v...)
}

func Alert(v ...interface{}) {
	BeeLogger.Alert(generateFmtStr(len(v)), v...)
}

// Critical logs a message at critical level.
func Critical(v ...interface{}) {
	BeeLogger.Critical(generateFmtStr(len(v)), v...)
}

// Error logs a message at error level.
func Error(v ...interface{}) {
	BeeLogger.Error(generateFmtStr(len(v)), v...)
}

// Warning logs a message at warning level.
func Warning(v ...interface{}) {
	BeeLogger.Warning(generateFmtStr(len(v)), v...)
}

// compatibility alias for Warning()
func Warn(v ...interface{}) {
	BeeLogger.Warn(generateFmtStr(len(v)), v...)
}

func Notice(v ...interface{}) {
	BeeLogger.Notice(generateFmtStr(len(v)), v...)
}

// Info logs a message at info level.
func Informational(v ...interface{}) {
	BeeLogger.Informational(generateFmtStr(len(v)), v...)
}

// compatibility alias for Warning()
func Info(v ...interface{}) {
	BeeLogger.Info(generateFmtStr(len(v)), v...)
}

// Debug logs a message at debug level.
func Debug(v ...interface{}) {
	BeeLogger.Debug(generateFmtStr(len(v)), v...)
}

// Trace logs a message at trace level.
// compatibility alias for Warning()
func Trace(v ...interface{}) {
	BeeLogger.Trace(generateFmtStr(len(v)), v...)
}

func generateFmtStr(n int) string {
	return strings.Repeat("%v ", n)
}

func init() {
	BeeLogger = DefaultLogs()
}
