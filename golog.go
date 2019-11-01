package golog

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	// mu protects read-write of the state
	mu sync.RWMutex
	// state is the global state of the package log level
	state = &globalState{
		currentLevel: InfoLevel,
	}
)
var (
	// DebugLogger is a standard logger use for debugging.
	DebugLogger = &stdLogger{level: DebugLevel, l: log.New(os.Stdout, DebugLevel.String(), log.LstdFlags|log.Lshortfile)}
	// TraceLogger is a standard logger use for tracing.
	TraceLogger = &stdLogger{level: TraceLevel, l: log.New(os.Stdout, TraceLevel.String(), log.LstdFlags|log.Lshortfile)}
	// InfoLogger is a standard logger info log.
	InfoLogger = &stdLogger{level: InfoLevel, l: log.New(os.Stdout, InfoLevel.String(), log.LstdFlags)}
	// WarningLogger is a standard logger warning log.
	WarningLogger = &stdLogger{level: WarningLevel, l: log.New(os.Stdout, WarningLevel.String(), log.LstdFlags)}
	// ErrorLogger is a standard logger use for printing errors.
	ErrorLogger = &stdLogger{level: ErrorLevel, l: log.New(os.Stdout, ErrorLevel.String(), log.LstdFlags|log.Lshortfile)}
	// DisabledLogger is a standard logger use to disable all logs.
	DisabledLogger = &stdLogger{level: DisabledLevel, l: log.New(ioutil.Discard, DisabledLevel.String(), log.LstdFlags|log.Lshortfile)}
)
var _ Logger = (*stdLogger)(nil)

// Level represents the log level of severity
// of the package.
type Level int

// When the global state level is lower than the
// logger level that currently in use log will print,
// otherwise log will not print.
//
// See stdLogger.isPrint
// debug < trace < info < warning < error < disabled
const (
	DebugLevel    Level = iota // DebugLevel use level for debug state. Useful for debugging that contains detailed information logging.
	TraceLevel                 // TraceLevel use level for trace state. Useful for tracing the code flow.
	InfoLevel                  // InfoLevel use level for info state
	WarningLevel               // WarningLevel use level for warning state
	ErrorLevel                 // ErrorLevel use level for error state. Higher than info because API needs the error message
	DisabledLevel              // DisabledLevel use level for disabled state
)

// Logger represents a general logger interface.
type Logger interface {
	Printf(format string, v ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	SetOutput(w io.Writer)
}

// String is to implement Stringer interface
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG: "
	case InfoLevel:
		return "INFO: "
	case ErrorLevel:
		return "ERROR: "
	case DisabledLevel:
		return "DISABLED: "
	case TraceLevel:
		return "TRACE: "
	case WarningLevel:
		return "WARNING: "
	}
	return ""
}

type globalState struct {
	currentLevel Level
}

func getState() *globalState {
	mu.Lock()
	defer mu.Unlock()
	return state
}
func setGlobalStateLevel(lvl Level) {
	mu.Lock()
	defer mu.Unlock()
	state.currentLevel = lvl
}

// SetLevel accepts log level to be set on the
// global state.
func SetLevel(lvl Level) {
	setGlobalStateLevel(lvl)
}

// NewStdLogger accepts level and return a
// standard logger that is bind to the level.
func NewStdLogger(level Level) Logger {
	return loggerFactory(level)
}
func loggerFactory(level Level) Logger {
	var l Logger
	switch level {
	case DebugLevel:
		l = DebugLogger
	case TraceLevel:
		l = TraceLogger
	case InfoLevel:
		l = InfoLogger
	case WarningLevel:
		l = WarningLogger
	case ErrorLevel:
		l = ErrorLogger
	case DisabledLevel:
		l = DisabledLogger
	}
	return l
}

const (
	stdCallDepth = 3
)

type stdLogger struct {
	level Level
	l     *log.Logger
}

func (l *stdLogger) Print(v ...interface{}) {
	if !l.isPrint() {
		return
	}
	l.Output(stdCallDepth, fmt.Sprint(v...))
}
func (l *stdLogger) Printf(format string, v ...interface{}) {
	if !l.isPrint() {
		return
	}
	l.Output(stdCallDepth, fmt.Sprintf(format, v...))
}
func (l *stdLogger) Println(v ...interface{}) {
	if !l.isPrint() {
		return
	}
	l.Output(stdCallDepth, fmt.Sprintln(v...))
}
func (l *stdLogger) Fatal(v ...interface{}) {
	if !l.isPrint() {
		return
	}
	l.Output(stdCallDepth, fmt.Sprint(v...))
	os.Exit(1)
}
func (l *stdLogger) Fatalf(format string, v ...interface{}) {
	if !l.isPrint() {
		return
	}
	l.Output(stdCallDepth, fmt.Sprintf(format, v...))
	os.Exit(1)
}
func (l *stdLogger) isPrint() bool {
	gstate := getState()
	if l.level < gstate.currentLevel {
		return false
	}
	return true
}
func (l *stdLogger) SetOutput(w io.Writer) {
	l.l.SetOutput(w)
}
func (l *stdLogger) Output(calldepth int, s string) {
	l.l.Output(calldepth, s)
}

// Debug is a convenient function that will be use for debugging.
func Debug(v ...interface{}) {
	if !DebugLogger.isPrint() {
		return
	}
	DebugLogger.Output(stdCallDepth, fmt.Sprintln(v...))
}

// Debugf is a convenient function that accepts format string
// and arguments that will be use for debugging.
func Debugf(format string, v ...interface{}) {
	if !DebugLogger.isPrint() {
		return
	}
	DebugLogger.Output(stdCallDepth, fmt.Sprintf(format, v...))
}

// Error is a convenient function that accepts arguments v
// and will be use to log error
func Error(v ...interface{}) {
	if !ErrorLogger.isPrint() {
		return
	}
	ErrorLogger.Output(stdCallDepth, fmt.Sprintln(v...))
}

// Errorf is a convenient function that accepts format string
// and arguments that will be use for to log error.
func Errorf(format string, v ...interface{}) {
	if !ErrorLogger.isPrint() {
		return
	}
	ErrorLogger.Output(stdCallDepth, fmt.Sprintf(format, v...))
}

// Info is a convenient function that accepts arguments v
// and will be use for info log.
func Info(v ...interface{}) {
	if !InfoLogger.isPrint() {
		return
	}
	InfoLogger.Output(stdCallDepth, fmt.Sprintln(v...))
}

// Infof is a convenient function that accepts format string
// and arguments that will be use for info log.
func Infof(format string, v ...interface{}) {
	if !InfoLogger.isPrint() {
		return
	}
	InfoLogger.Output(stdCallDepth, fmt.Sprintf(format, v...))
}

// Trace is a convenient function that accepts argument v
// and will be use for tracing.
func Trace(v ...interface{}) {
	if !TraceLogger.isPrint() {
		return
	}
	TraceLogger.Output(stdCallDepth, fmt.Sprintln(v...))
}

// Tracef is a convenient function that accepts format string
// and arguments v. It will be useful for tracing operation.
func Tracef(format string, v ...interface{}) {
	if !TraceLogger.isPrint() {
		return
	}
	TraceLogger.Output(stdCallDepth, fmt.Sprintf(format, v...))
}

// Warning is a convenient function that accepts argument v
// and logs the v in a warning state.
func Warning(v ...interface{}) {
	if !WarningLogger.isPrint() {
		return
	}
	WarningLogger.Output(stdCallDepth, fmt.Sprintln(v...))
}

// Warningf is a convenient function that accepts argument v and
// a format string and logs in a warning state.
func Warningf(format string, v ...interface{}) {
	if !WarningLogger.isPrint() {
		return
	}
	WarningLogger.Output(stdCallDepth, fmt.Sprintf(format, v...))
}

// Fatal is a convenient function that accepts argument v
// and will be use to abort the program with an error log.
func Fatal(v ...interface{}) {
	if !ErrorLogger.isPrint() {
		return
	}
	ErrorLogger.Output(stdCallDepth, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is a convenient function that accepts format string
// and arguments that will be use to abort the program with
// an error log.
func Fatalf(format string, v ...interface{}) {
	if !ErrorLogger.isPrint() {
		return
	}
	ErrorLogger.Output(stdCallDepth, fmt.Sprintf(format, v...))
	os.Exit(1)
}
