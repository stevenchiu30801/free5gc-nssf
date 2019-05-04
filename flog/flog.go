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
type ServiceLogger struct {
    serviceName string
    debugLogger *log.Logger
    infoLogger *log.Logger
    warnLogger *log.Logger
    errorLogger *log.Logger
    mute bool
}

// Service logger type
var (
    Default ServiceLogger
    System ServiceLogger
    Nsselection ServiceLogger
    Nssaiavailability ServiceLogger
)

// Service name
const (
    SYSTEM_SERVICE string = "Nnssf-System"
    NSSELECTION_SERVICE string = "Nnssf-Nsselection"
    NSSAIAVAILABILITY_SERVICE string = "Nnssf-Nssaiavailability"
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
const defaultServiceName string = "Nnssf-Default"
const defaultLogWithColor bool = false

func getEscapeCode(style int, color int) string {
    return fmt.Sprintf("%s[%d;%dm", Escape, style, color)
}

func reset() string {
    return fmt.Sprintf("%s[%dm", Escape, NoEffect)
}

func init() {
    Default.InitLogger(defaultServiceName, defaultLogWithColor)
    System.InitLogger(SYSTEM_SERVICE, true)
    Nsselection.InitLogger(NSSELECTION_SERVICE, true)
    Nssaiavailability.InitLogger(NSSAIAVAILABILITY_SERVICE, true)
}

func (s *ServiceLogger) InitLogger(service string, colorInd bool) {
    if service == "" {
        s.serviceName = defaultServiceName
    } else {
        s.serviceName = service
    }
    s.mute = false

    if colorInd == true {
        s.debugLogger = log.New(os.Stdout, getEscapeCode(logStyle, debugColor), log.Ldate|log.Ltime)
        s.infoLogger = log.New(os.Stdout, getEscapeCode(logStyle, infoColor), log.Ldate|log.Ltime)
        s.warnLogger = log.New(os.Stdout, getEscapeCode(logStyle, warnColor), log.Ldate|log.Ltime)
        s.errorLogger = log.New(os.Stdout, getEscapeCode(logStyle, errorColor), log.Ldate|log.Ltime)
    } else {
        s.debugLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
        s.infoLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
        s.warnLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
        s.errorLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
    }
}

func MuteLog() {
    Default.MuteLog()
}

func (s *ServiceLogger) MuteLog() {
    s.mute = true
}

func Debugf(format string, v ...interface{}) {
    Default.Debugf(format, v...)
}

func (s *ServiceLogger) Debugf(format string, v ...interface{}) {
    if s.mute == true {
        return
    }
    format = "- " + s.serviceName + " - DEBUG - " + format + reset()
    s.debugLogger.Printf(format, v...)
}

func Infof(format string, v ...interface{}) {
    Default.Infof(format, v...)
}

func (s *ServiceLogger) Infof(format string, v ...interface{}) {
    if s.mute == true {
        return
    }
    format = "- " + s.serviceName + " - INFO - " + format + reset()
    s.infoLogger.Printf(format, v...)
}

func Warnf(format string, v ...interface{}) {
    Default.Warnf(format, v...)
}

func (s *ServiceLogger) Warnf(format string, v ...interface{}) {
    if s.mute == true {
        return
    }
    format = "- " + s.serviceName + " - WARN - " + format + reset()
    s.warnLogger.Printf(format, v...)
}

func Errorf(format string, v ...interface{}) {
    Default.Errorf(format, v...)
}

func (s *ServiceLogger) Errorf(format string, v ...interface{}) {
    if s.mute == true {
        return
    }
    format = "- " + s.serviceName + " - ERROR - " + format + reset()
    s.errorLogger.Printf(format, v...)
}

func Fatal(v ...interface{}) {
    Default.Fatal(v...)
}

func (s *ServiceLogger) Fatal(v ...interface{}) {
    s.errorLogger.Fatal(v...)
}
