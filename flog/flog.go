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

// TODO: Use a base log function and log structs for each level to remove duplicate code

// TODO: Support multiple logging backends and different logging levels.
//       Or using go-logging (https://github.com/op/go-logging).

// Service logger
var (
    serviceName string
    debugLogger *log.Logger
    infoLogger *log.Logger
    warnLogger *log.Logger
    errorLogger *log.Logger
    mute bool = false
)

// Log style and color
const (
    logStyle = NoEffect
    debugColor = HiMagenta
    infoColor = HiWhite
    warnColor = HiYellow
    errorColor = HiRed
)

// Default log style if loggers are not initizalized
const defaultServiceName string = "DefaultFlog"
const defaultLogWithColor bool = false

func getEscapeCode(style int, color int) string {
    return fmt.Sprintf("%s[%d;%dm", Escape, style, color)
}

func reset() string {
    return fmt.Sprintf("%s[%dm", Escape, NoEffect)
}

func InitLog(service string, colorInd bool) {
    serviceName = service
    mute = false

    if colorInd == true {
        debugLogger = log.New(os.Stdout, getEscapeCode(logStyle, debugColor), log.Ldate|log.Ltime)
        infoLogger = log.New(os.Stdout, getEscapeCode(logStyle, infoColor), log.Ldate|log.Ltime)
        warnLogger = log.New(os.Stdout, getEscapeCode(logStyle, warnColor), log.Ldate|log.Ltime)
        errorLogger = log.New(os.Stdout, getEscapeCode(logStyle, errorColor), log.Ldate|log.Ltime)
    } else {
        debugLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
        infoLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
        warnLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
        errorLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
    }
}

func MuteLog() {
    mute = true
}

func Debugf(format string, v ...interface{}) {
    if mute == true {
        return
    }
    if debugLogger == nil {
        InitLog(defaultServiceName, defaultLogWithColor)
    }
    format = "- " + serviceName + " - DEBUG - " + format + reset()
    debugLogger.Printf(format, v...)
}

func Infof(format string, v ...interface{}) {
    if mute == true {
        return
    }
    if infoLogger == nil {
        InitLog(defaultServiceName, defaultLogWithColor)
    }
    format = "- " + serviceName + " - INFO - " + format + reset()
    infoLogger.Printf(format, v...)
}

func Warnf(format string, v ...interface{}) {
    if mute == true {
        return
    }
    if warnLogger == nil {
        InitLog(defaultServiceName, defaultLogWithColor)
    }
    format = "- " + serviceName + " - WARN - " + format + reset()
    warnLogger.Printf(format, v...)
}

func Errorf(format string, v ...interface{}) {
    if mute == true {
        return
    }
    if errorLogger == nil {
        InitLog(defaultServiceName, defaultLogWithColor)
    }
    format = "- " + serviceName + " - ERROR - " + format + reset()
    errorLogger.Printf(format, v...)
}

func Fatal(v ...interface{}) {
    if errorLogger == nil {
        InitLog(defaultServiceName, false)
    }
    errorLogger.Fatal(v...)
}
