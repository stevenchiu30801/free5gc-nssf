/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package flog

import (
	"log"
    "strconv"
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

func getEscapeCode(style int, color int) string {
    return "\033[" + strconv.Itoa(style) + ";" + strconv.Itoa(color) + "m"
}

func InitLog(service string) {
    serviceName = service
    debugLogger = log.New(os.Stdout, getEscapeCode(logStyle, debugColor), log.Ldate|log.Ltime)
    infoLogger = log.New(os.Stdout, getEscapeCode(logStyle, infoColor), log.Ldate|log.Ltime)
    warnLogger = log.New(os.Stdout, getEscapeCode(logStyle, warnColor), log.Ldate|log.Ltime)
    errorLogger = log.New(os.Stdout, getEscapeCode(logStyle, errorColor), log.Ldate|log.Ltime)
}

func Debug(format string, v ...interface{}) {
    format = "- " + serviceName + " - DEBUG - " + format
    debugLogger.Printf(format, v...)
}

func Info(format string, v ...interface{}) {
    format = "- " + serviceName + " - INFO - " + format
    infoLogger.Printf(format, v...)
}

func Warn(format string, v ...interface{}) {
    format = "- " + serviceName + " - WARN - " + format
    warnLogger.Printf(format, v...)
}

func Error(format string, v ...interface{}) {
    format = "- " + serviceName + " - ERROR - " + format
    errorLogger.Printf(format, v...)
}

func Fatal(v ...interface{}) {
    errorLogger.Fatal(v...)
}
