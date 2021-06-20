package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

const CALLER_FRAME = 3
const LOGGER_MAIN = "main"

var loggers = map[string]*log.Logger{}

func init() {
	CreateLogger(LOGGER_MAIN)
}

func CreateLogger(n string) *log.Logger {
	if _, ok := loggers[n]; !ok {
		loggers[n] = log.New()
		loggers[n].SetReportCaller(true)
		loggers[n].SetFormatter(&log.TextFormatter{
			DisableLevelTruncation: true,
			PadLevelText:           true,
			CallerPrettyfier: func(f *runtime.Frame) (function, file string) {
				function, file = fileInfo(CALLER_FRAME + 7)
				return
			},
		})
	}
	return loggers[n]
}

func SetLevel(level log.Level) {
	log.SetLevel(level)
}

func Logger(n string) *log.Logger {
	return loggers[n]
}

func Main() *log.Logger {
	return Logger(LOGGER_MAIN)
}

func MainEntry() *log.Entry {
	// fn, file := fileInfo(CALLER_FRAME)
	return Main().WithFields(log.Fields{})
	// .WithField("file", file).WithField("function", fn)
}

func Loggers() map[string]*log.Logger {
	return loggers
}

func Log(level log.Level, args ...interface{}) {
	MainEntry().Log(level, args...)
}

func Trace(args ...interface{}) {
	MainEntry().Trace(args...)
}

func Debug(args ...interface{}) {
	MainEntry().Debug(args...)
}

func Print(args ...interface{}) {
	MainEntry().Print(args...)
}

func Info(args ...interface{}) {
	MainEntry().Info(args...)
}

func Warn(args ...interface{}) {
	MainEntry().Warn(args...)
}

func Warning(args ...interface{}) {
	MainEntry().Warning(args...)
}

func Error(args ...interface{}) {
	MainEntry().Error(args...)
}

func Fatal(args ...interface{}) {
	MainEntry().Fatal(args...)
}

func Panic(args ...interface{}) {
	MainEntry().Panic(args...)
}

// Logger Printf family functions

func Logf(level log.Level, format string, args ...interface{}) {
	MainEntry().Logf(level, format, args...)
}

func Tracef(format string, args ...interface{}) {
	MainEntry().Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	MainEntry().Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	MainEntry().Infof(format, args...)
}

func Printf(format string, args ...interface{}) {
	MainEntry().Printf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	MainEntry().Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	MainEntry().Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	MainEntry().Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	MainEntry().Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	MainEntry().Panicf(format, args...)
}

// Logger Println family functions

func Logln(level log.Level, args ...interface{}) {
	MainEntry().Logln(level, args...)
}

func Traceln(args ...interface{}) {
	MainEntry().Traceln(args...)
}

func Debugln(args ...interface{}) {
	MainEntry().Debugln(args...)
}

func Infoln(args ...interface{}) {
	MainEntry().Infoln(args...)
}

func Println(args ...interface{}) {
	MainEntry().Println(args...)
}

func Warnln(args ...interface{}) {
	MainEntry().Warnln(args...)
}

func Warningln(args ...interface{}) {
	MainEntry().Warningln(args...)
}

func Errorln(args ...interface{}) {
	MainEntry().Errorln(args...)
}

func Fatalln(args ...interface{}) {
	MainEntry().Fatalln(args...)
}

func Panicln(args ...interface{}) {
	MainEntry().Panicln(args...)
}

func fileInfo(skip int) (function, file string) {
	var line int
	var ok bool
	var pc uintptr
	pc, file, line, ok = runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		if fnFrame := runtime.FuncForPC(pc); fnFrame != nil {
			function = fnFrame.Name()
		}

		filePrefix := filepath.Dir(filepath.Dir(file))

		file = strings.TrimPrefix(file, filePrefix)
		file = strings.TrimPrefix(file, "/")

		function = filepath.Base(function)

	}
	return fmt.Sprintf("%s()", function), fmt.Sprintf("%s:%d", file, line)
}
