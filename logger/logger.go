package logger

type LogLevel string

const (
	Info  LogLevel = "Info"
	Warn  LogLevel = "Warn"
	Debug LogLevel = "Debug"
	Error LogLevel = "Error"
	Fatal LogLevel = "Fatal"
)

type Logger interface {
	Infoln(fn, tid string, msg string)
	Infof(fn, tid string, format string, args ...interface{})
	Warnln(fn, tid string, msg string)
	Errorln(fn, tid string, msg string)
	Errorf(fn, tid string, format string, args ...interface{})
	Print(level LogLevel, fn, tid string, msg string)
}

func GetDefaultStructLogger(verbose bool) Logger {
	logger := NewZeroLevelLogger(verbose)
	logger.Infoln("GetDefaultStructLogger", "0000000000", "Running in verbose mode")
	return logger
}
