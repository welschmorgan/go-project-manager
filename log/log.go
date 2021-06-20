package log

import (
	"fmt"
	"path"
	"runtime"

	log "github.com/sirupsen/logrus"
)

var loggers = map[string]*log.Entry{
	"main": nil,
}

func init() {
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function, file string) {
			filename := path.Base(f.File)
			function = fmt.Sprintf("%s()", path.Base(f.Function))
			file = fmt.Sprintf("%s:%d", filename, f.Line)
			return
		},
	})
	loggers["main"] = log.WithFields(log.Fields{
		"module": "app",
	})
}

func SetLevel(level log.Level) {
	log.SetLevel(level)
}

func Main() *log.Entry {
	return loggers["main"]
}

func Loggers() map[string]*log.Entry {
	return loggers
}

func Log(level log.Level, args ...interface{}) {
	Main().Log(level, args...)
}

func Trace(args ...interface{}) {
	Main().Trace(args...)
}

func Debug(args ...interface{}) {
	Main().Debug(args...)
}

func Print(args ...interface{}) {
	Main().Print(args...)
}

func Info(args ...interface{}) {
	Main().Info(args...)
}

func Warn(args ...interface{}) {
	Main().Warn(args...)
}

func Warning(args ...interface{}) {
	Main().Warning(args...)
}

func Error(args ...interface{}) {
	Main().Error(args...)
}

func Fatal(args ...interface{}) {
	Main().Fatal(args...)
}

func Panic(args ...interface{}) {
	Main().Panic(args...)
}

// Entry Printf family functions

func Logf(level log.Level, format string, args ...interface{}) {
	Main().Logf(level, format, args...)
}

func Tracef(format string, args ...interface{}) {
	Main().Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	Main().Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	Main().Infof(format, args...)
}

func Printf(format string, args ...interface{}) {
	Main().Printf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	Main().Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	Main().Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	Main().Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	Main().Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	Main().Panicf(format, args...)
}

// Entry Println family functions

func Logln(level log.Level, args ...interface{}) {
	Main().Logln(level, args...)
}

func Traceln(args ...interface{}) {
	Main().Traceln(args...)
}

func Debugln(args ...interface{}) {
	Main().Debugln(args...)
}

func Infoln(args ...interface{}) {
	Main().Infoln(args...)
}

func Println(args ...interface{}) {
	Main().Println(args...)
}

func Warnln(args ...interface{}) {
	Main().Warnln(args...)
}

func Warningln(args ...interface{}) {
	Main().Warningln(args...)
}

func Errorln(args ...interface{}) {
	Main().Errorln(args...)
}

func Fatalln(args ...interface{}) {
	Main().Fatalln(args...)
}

func Panicln(args ...interface{}) {
	Main().Panicln(args...)
}
