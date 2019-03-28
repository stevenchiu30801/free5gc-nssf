/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package flog

import (
	"log"
    "os"
)

// TODO: Support multiple logging backends and different logging levels.
//       Or using go-logging (https://github.com/op/go-logging).

var(
    debugLogger *log.Logger
    infoLogger *log.Logger
    warnLogger *log.Logger
    errorLogger *log.Logger
)

func LogInit() {
    debugLogger = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime)
    infoLogger = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
    warnLogger = log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime)
    errorLogger = log.New(os.Stdout, "[ERROR] ", log.Ldate|log.Ltime)
}

func Debug(format string, v ...interface{}) {
    debugLogger.Printf(format, v...)
}

func Info(format string, v ...interface{}) {
    infoLogger.Printf(format, v...)
}

func Warn(format string, v ...interface{}) {
    warnLogger.Printf(format, v...)
}

func Error(format string, v ...interface{}) {
    errorLogger.Printf(format, v...)
}

func Fatal(v ...interface{}) {
    errorLogger.Fatal(v...)
}
