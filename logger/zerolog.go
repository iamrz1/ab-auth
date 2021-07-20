package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type zeroLevelLogger struct {
	lgr     zerolog.Logger
	verbose bool
}

func NewZeroLevelLogger(verbose bool) Logger {
	//if verbose{
	//	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	//} else {
	//	zerolog.SetGlobalLevel(zerolog.Disabled)
	//}

	sublogger := log.With().CallerWithSkipFrameCount(3).Stack().
		Str("service", "auth").
		Logger()

	return &zeroLevelLogger{
		lgr:     sublogger,
		verbose: verbose,
	}
}
func (l zeroLevelLogger) Infoln(fn, tid string, msg string) {
	level := zerolog.Disabled
	if l.verbose {
		level = zerolog.InfoLevel
	}
	l.lgr.WithLevel(level).Str("function", fn).Str("tid", tid).Msg(msg)
}

func (l zeroLevelLogger) Infof(fn, tid string, format string, args ...interface{}) {
	level := zerolog.Disabled
	if l.verbose {
		level = zerolog.InfoLevel
	}
	l.lgr.WithLevel(level).Str("function", fn).Str("tid", tid).Msgf(format, args...)
}

func (l zeroLevelLogger) Warnln(fn, tid string, msg string) {
	level := zerolog.Disabled
	if l.verbose {
		level = zerolog.WarnLevel
	}
	l.lgr.WithLevel(level).Str("function", fn).Str("tid", tid).Msg(msg)
}

func (l zeroLevelLogger) Errorln(fn, tid string, msg string) {
	level := zerolog.Disabled
	if l.verbose {
		level = zerolog.ErrorLevel
	}
	l.lgr.WithLevel(level).Str("function", fn).Str("tid", tid).Msg(msg)
}
func (l zeroLevelLogger) Errorf(fn, tid string, format string, args ...interface{}) {
	level := zerolog.Disabled
	if l.verbose {
		level = zerolog.ErrorLevel
	}
	l.lgr.WithLevel(level).Str("function", fn).Str("tid", tid).Msgf(format, args...)
}

func (l zeroLevelLogger) Print(level LogLevel, fn, tid string, msg string) {
	l.lgr.Log().Str("level", getZeroLevel(level).String()).Str("function", fn).Str("tid", tid).Msg(msg)
}

func getZeroLevel(level LogLevel) zerolog.Level {
	switch level {
	case Info:
		return zerolog.InfoLevel
	case Warn:
		return zerolog.WarnLevel
	case Debug:
		return zerolog.DebugLevel
	case Error:
		return zerolog.ErrorLevel
	case Fatal:
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}
