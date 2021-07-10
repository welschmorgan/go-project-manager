// The log package simplify configuration of the logrus package.
package log

import (
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/fs"
)

// The number of frames to climb to find the real log caller
const CALLER_FRAME = 3

// The padding size of the function column
const CALLER_FUNC_PAD_SIZE = 30

// The padding size of the file column
const CALLER_FILE_PAD_SIZE = 20

// The default logger name
const LOGGER_MAIN = "main"

// The map of created loggers
var loggers = map[string]*log.Logger{}

// PrettifyCallers retrieves the correct caller infos, without project directory.
func PrettifyCaller(callerFrameOffset int) func(f *runtime.Frame) (function, file string) {
	return func(f *runtime.Frame) (function, file string) {
		function, file = fileInfo(CALLER_FRAME + callerFrameOffset)
		function = fmt.Sprintf(fmt.Sprintf("%%-%ds", CALLER_FUNC_PAD_SIZE+3+2), function)
		if len(function) > CALLER_FUNC_PAD_SIZE {
			removeCount := len(function) - (CALLER_FUNC_PAD_SIZE + 3 + 2)
			halfPos := len(function) / 2
			function = fmt.Sprintf("%s...%s", string(function[0:halfPos-removeCount/2]), string(function[halfPos+removeCount/2:]))
		}
		// function = fmt.Sprintf(fmt.Sprintf("%%%d.%ds", CALLER_FUNC_PAD_SIZE, CALLER_FUNC_PAD_SIZE), function)
		file = fmt.Sprintf(fmt.Sprintf("%%%-d.%ds", CALLER_FILE_PAD_SIZE, CALLER_FILE_PAD_SIZE), file)
		return
	}
}

// Setup prepares the whole logging mecanism.
// It creates file and console loggers with file rolling.
// It should be called after config, but before everything else
func Setup() error {
	log.SetLevel(config.Get().Verbose.LogLevel())

	logDir := config.Get().Workspace.LogFolder()
	fs.Mkdir(logDir.Expand())

	// instanciate main logger
	Logger(LOGGER_MAIN)

	return nil
}

// SetOutput defines all loggers output
func SetOutput(w io.Writer) {
	log.SetOutput(w)
	for _, logger := range loggers {
		logger.SetOutput(w)
	}
}

// SetLevel defines all loggers level
func SetLevel(level log.Level) {
	log.SetLevel(level)
	for _, logger := range loggers {
		logger.SetLevel(level)
	}
}

// Logger retrieves or creates a logger by name
func Logger(n string) *log.Logger {
	if _, ok := loggers[n]; !ok {
		loggers[n] = log.New()
		loggers[n].SetReportCaller(true)
		loggers[n].SetFormatter(&log.TextFormatter{
			DisableLevelTruncation: true,
			PadLevelText:           true,
			CallerPrettyfier:       PrettifyCaller(7),
		})
		loggers[n].SetLevel(log.GetLevel())

		// file rotation
		// loggers[n].SetOutput(io.MultiWriter(os.Stdout, rotator))

		logDir := config.Get().Workspace.LogFolder()
		if rotatorHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
			Filename:   filepath.Join(logDir.Expand(), n+".log"),
			MaxSize:    5, // megabytes
			MaxBackups: 3,
			MaxAge:     28, //days
			Level:      log.GetLevel(),
			Formatter: &log.TextFormatter{
				TimestampFormat:        time.RFC822,
				DisableColors:          true,
				PadLevelText:           true,
				DisableLevelTruncation: true,
				CallerPrettyfier:       PrettifyCaller(9),
			},
		}); err != nil {
			panic(err)
		} else {
			// log.SetOutput(io.MultiWriter(os.Stdout, rotator))
			loggers[n].AddHook(rotatorHook)
		}
	}
	return loggers[n]
}

// Main retrieves the main logger instance
func Main() *log.Logger {
	return Logger(LOGGER_MAIN)
}

// MainEntry creates an entry in the main logger
func MainEntry() *log.Entry {
	// fn, file := fileInfo(CALLER_FRAME)
	return Main().WithFields(log.Fields{})
	// .WithField("file", file).WithField("function", fn)
}

// Loggers retrieves all loggers instances
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
