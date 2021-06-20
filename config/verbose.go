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
)

var VerboseLevels = []VerboseLevel{
	NoVerbose,
	LowVerbose,
	Verbose,
	HighVerbose,
}

func (v VerboseLevel) TextualRepresentations() []string {
	switch v {
	case NoVerbose:
		return []string{"none", "no", "0"}
	case LowVerbose:
		return []string{"low"}
	case Verbose:
		return []string{"normal", ""}
	case HighVerbose:
		return []string{"high", "max"}
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
		return logrus.TraceLevel
	default:
		panic(fmt.Sprintf("unknown verbose level: %d", v))
	}
}

func (v VerboseLevel) String() string {
	return v.Name()
}
