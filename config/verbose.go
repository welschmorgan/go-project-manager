package config

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type VerboseLevel uint8

const (
	NoVerbose VerboseLevel = iota
	LowVerbose
	Verbose
	HighVerbose
	MaxVerbose
)

var VerboseLevels = []VerboseLevel{
	NoVerbose,
	LowVerbose,
	Verbose,
	HighVerbose,
	MaxVerbose,
}

func (v VerboseLevel) TextualRepresentations() []string {
	switch v {
	case NoVerbose:
		return []string{"none", "no", "0", "error"}
	case LowVerbose:
		return []string{"low", "warn"}
	case Verbose:
		return []string{"normal", "info"}
	case HighVerbose:
		return []string{"high", "debug"}
	case MaxVerbose:
		return []string{"max", "trace"}
	default:
		panic(fmt.Sprintf("unknown verbose level: %d", v))
	}
}

func (v VerboseLevel) Name() string {
	switch v {
	case NoVerbose:
		return "none"
	case LowVerbose:
		return "low"
	case Verbose:
		return "verbose"
	case HighVerbose:
		return "high"
	case MaxVerbose:
		return "max"
	default:
		panic(fmt.Sprintf("unknown verbose level: %d", v))
	}
}

func (v VerboseLevel) LogLevel() logrus.Level {
	switch v {
	case NoVerbose:
		return logrus.ErrorLevel
	case LowVerbose:
		return logrus.WarnLevel
	case Verbose:
		return logrus.InfoLevel
	case HighVerbose:
		return logrus.DebugLevel
	case MaxVerbose:
		return logrus.TraceLevel
	default:
		panic(fmt.Sprintf("unknown verbose level: %d", v))
	}
}

func (v VerboseLevel) String() string {
	return v.Name()
}
