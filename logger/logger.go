package logger

type LogLevel string

const (
	Info  LogLevel = "Info"
	Warn  LogLevel = "Warn"
	Debug LogLevel = "Debug"
	Error LogLevel = "Error"
	Fatal LogLevel = "Fatal"
)

// Logger is an interface to print logs
type Logger interface {
	Print(v ...interface{})
	Println(v ...interface{})
	Printf(format string, v ...interface{})
	Errorf(format string, args ...interface{})
}

type LaveledLogger interface {
	Println(v ...interface{})
	Warnln(v ...interface{})
	Errorln(v ...interface{})
	Printf(format string, args ...interface{})
	Print(level LogLevel, v ...interface{})
	Errorf(level LogLevel, format string, args ...interface{})
}

type StructLogger interface {
	Infoln(fn, tid string, msg string)
	Infof(fn, tid string, format string, args ...interface{})
	Warnln(fn, tid string, msg string)
	Errorln(fn, tid string, msg string)
	Errorf(fn, tid string, format string, args ...interface{})
	Print(level LogLevel, fn, tid string, msg string)
}

var (
	DefaultOutStructLogger StructLogger
)

func init() {
	DefaultOutStructLogger = NewZeroLevelLogger()
}
