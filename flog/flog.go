/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package flog

import (
    "fmt"
	"log"
    "os"
)

// TODO: Support multiple logging backends and different logging levels.
//       Or using go-logging (https://github.com/op/go-logging).

// Service logger
var (
    serviceName string
    debugLogger *log.Logger
    infoLogger *log.Logger
    warnLogger *log.Logger
    errorLogger *log.Logger
)

// Log style and color
const (
    logStyle = NoEffect
    debugColor = HiMagenta
    infoColor = HiWhite
    warnColor = HiYellow
    errorColor = HiRed
)

// Default service name shown in log if loggers are not initizalized
const defaultServiceName string = "DefaultFlog"

func getEscapeCode(style int, color int) string {
    return fmt.Sprintf("%s[%d;%dm", Escape, style, color)
}

func reset() string {
    return fmt.Sprintf("%s[%dm", Escape, NoEffect)
}

func InitLog(service string) {
    serviceName = service

    debugLogger = log.New(os.Stdout, getEscapeCode(logStyle, debugColor), log.Ldate|log.Ltime)
    infoLogger = log.New(os.Stdout, getEscapeCode(logStyle, infoColor), log.Ldate|log.Ltime)
    warnLogger = log.New(os.Stdout, getEscapeCode(logStyle, warnColor), log.Ldate|log.Ltime)
    errorLogger = log.New(os.Stdout, getEscapeCode(logStyle, errorColor), log.Ldate|log.Ltime)
}

func Debug(format string, v ...interface{}) {
    if debugLogger == nil {
        InitLog(defaultServiceName)
    }
    format = "- " + serviceName + " - DEBUG - " + format + reset()
    debugLogger.Printf(format, v...)
}

func Info(format string, v ...interface{}) {
    if infoLogger == nil {
        InitLog(defaultServiceName)
    }
    format = "- " + serviceName + " - INFO - " + format + reset()
    infoLogger.Printf(format, v...)
}

func Warn(format string, v ...interface{}) {
    if warnLogger == nil {
        InitLog(defaultServiceName)
    }
    format = "- " + serviceName + " - WARN - " + format + reset()
    warnLogger.Printf(format, v...)
}

func Error(format string, v ...interface{}) {
    if errorLogger == nil {
        InitLog(defaultServiceName)
    }
    format = "- " + serviceName + " - ERROR - " + format + reset()
    errorLogger.Printf(format, v...)
}

func Fatal(v ...interface{}) {
    errorLogger.Fatal(v...)
}
