package golog

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"strings"
	"testing"
)

// Requirements:
// - When set level is info level, error logs should not print.
// - When set level is error level, error and info log will print.
func TestLogger_Print(t *testing.T) {
	t.Run("Contains cases", func(t *testing.T) {
		containsCases := []struct {
			name   string
			lvl    Level
			setLvl Level
			want   string
			input  string
			out    *bytes.Buffer
		}{
			// Info Level
			{name: "when printing message should be in the output", lvl: InfoLevel, setLvl: InfoLevel, input: "Hello Logging World!", want: "Hello Logging World!", out: &bytes.Buffer{}},
			{name: "when the global level is higher than the stdLogger level", lvl: InfoLevel, setLvl: ErrorLevel, input: "Hello World", want: "", out: &bytes.Buffer{}},
			{name: "when the global level is lower than the stdlogger level", lvl: InfoLevel, setLvl: DebugLevel, input: "Hello World!", want: "Hello World!", out: &bytes.Buffer{}},
			// Debug Level
			{name: "when printing message should be in the output", lvl: DebugLevel, setLvl: DebugLevel, input: "Hello Logging World!", want: "Hello Logging World!", out: &bytes.Buffer{}},
			{name: "when the global level is higher than the stdLogger level", lvl: DebugLevel, setLvl: InfoLevel, input: "Hello World!", want: "", out: nil},
			{name: "when the global level is lower than or equal to the stdlogger level", lvl: DebugLevel, setLvl: DebugLevel, input: "Hello World!", want: "Hello World!", out: &bytes.Buffer{}},
			{name: "when printing the line captured should be the client call level", lvl: DebugLevel, setLvl: DebugLevel, input: "Hello World!", want: "log_test.go", out: &bytes.Buffer{}},
			// Error Level
			{name: "when printing message should be in the output", lvl: ErrorLevel, setLvl: DebugLevel, input: "Hello Logging World!", want: "Hello Logging World!", out: &bytes.Buffer{}},
			{name: "when the global level is higher than the stdLogger level", lvl: ErrorLevel, setLvl: DisabledLevel, input: "Hello World!", want: "", out: &bytes.Buffer{}},
			{name: "when printing the line captured should be the client call level", lvl: ErrorLevel, setLvl: ErrorLevel, input: "Hello World!", want: "log_test.go", out: &bytes.Buffer{}},
		}
		for _, c := range containsCases {
			t.Run(fmt.Sprintf("%s:%s", c.lvl, c.name), func(t *testing.T) {
				l := &stdLogger{
					level: c.lvl,
					l:     log.New(c.out, c.lvl.String(), log.Lshortfile),
				}
				SetLevel(c.setLvl)
				l.Print(c.input, c.input)
				got := c.out.String()
				assert.Contains(t, got, c.want)
			})
		}
	})
	t.Run("Empty value", func(t *testing.T) {
		t.Run("when set level is info level, error logs should not print", func(t *testing.T) {
			out := &bytes.Buffer{}
			l := &stdLogger{
				level: ErrorLevel,
				l:     log.New(out, ErrorLevel.String(), log.LstdFlags|log.Lshortfile),
			}
			SetLevel(InfoLevel)
			l.Print("Hello world")
			want := "Hello world"
			got := out.String()
			assert.Contains(t, got, want)
		})
	})
}
func TestTraceLoggerPrintln(t *testing.T) {
	t.Run("when global state is higher than current log level", func(t *testing.T) {
		out := &bytes.Buffer{}
		TraceLogger.SetOutput(out)
		TraceLogger.Println("Hello World")
		assert.Empty(t, out.String())
	})
	t.Run("when global state is lower then current log level", func(t *testing.T) {
		SetLevel(DebugLevel)
		out := &bytes.Buffer{}
		TraceLogger.SetOutput(out)
		TraceLogger.Println("Hello World")
		assert.Contains(t, out.String(), "Hello World")
	})
	t.Run("when global state is equal to the current log level", func(t *testing.T) {
		SetLevel(TraceLevel)
		out := &bytes.Buffer{}
		TraceLogger.SetOutput(out)
		TraceLogger.Println("Hello World")
		assert.Contains(t, out.String(), "Hello World")
	})
}
func TestWarningLoggerPrintln(t *testing.T) {
	t.Run("when global state is higher than current log level", func(t *testing.T) {
		out := &bytes.Buffer{}
		SetLevel(ErrorLevel)
		WarningLogger.SetOutput(out)
		WarningLogger.Println("Hello World")
		assert.Empty(t, out.String())
	})
	t.Run("when global state is lower then current log level", func(t *testing.T) {
		SetLevel(DebugLevel)
		out := &bytes.Buffer{}
		WarningLogger.SetOutput(out)
		WarningLogger.Println("Hello World")
		assert.Contains(t, out.String(), "Hello World")
	})
	t.Run("when global state is equal to the current log level", func(t *testing.T) {
		SetLevel(WarningLevel)
		out := &bytes.Buffer{}
		WarningLogger.SetOutput(out)
		WarningLogger.Println("Hello World")
		assert.Contains(t, out.String(), "Hello World")
	})
	t.Run("contains warning prefix", func(t *testing.T) {
		SetLevel(WarningLevel)
		out := &bytes.Buffer{}
		WarningLogger.SetOutput(out)
		WarningLogger.Println("Hello World")
		assert.Contains(t, out.String(), "Hello World")
		assert.True(t, strings.HasPrefix(out.String(), WarningLevel.String()))
	})
}
func TestLogger_Printf(t *testing.T) {
	cases := []struct {
		name   string
		lvl    Level
		setLvl Level
		want   string
		format string
		input  string
		out    *bytes.Buffer
	}{
		// Info Level
		{name: "when printing message should be in the output", lvl: InfoLevel, setLvl: InfoLevel, format: "Hello Logging %s!", input: "World", want: "Hello Logging World!", out: &bytes.Buffer{}},
		{name: "when the global level is higher than the stdLogger level", lvl: InfoLevel, setLvl: ErrorLevel, format: "Hello %s!", input: "World", want: "", out: &bytes.Buffer{}},
		{name: "when the global level is lower than the stdlogger level", lvl: InfoLevel, setLvl: DebugLevel, format: "Hello %s!", input: "World", want: "Hello World!", out: &bytes.Buffer{}},
		// Debug Level
		{name: "when printing message should be in the output", lvl: DebugLevel, setLvl: DebugLevel, format: "Hello Logging %s!", input: "World", want: "Hello Logging World!", out: &bytes.Buffer{}},
		{name: "when the global level is higher than the stdLogger level", lvl: DebugLevel, setLvl: InfoLevel, format: "Hello %s!", input: "World", want: "", out: &bytes.Buffer{}},
		{name: "when the global level is lower than or equal to the stdlogger level", lvl: DebugLevel, setLvl: DebugLevel, format: "Hello %s!", input: "World", want: "Hello World!", out: &bytes.Buffer{}},
		{name: "when printing the line captured should be the client call level", lvl: DebugLevel, setLvl: DebugLevel, format: "Hello %s!", input: "World", want: "log_test.go", out: &bytes.Buffer{}},
		// Error Level
		{name: "when printing message should be in the output", lvl: ErrorLevel, setLvl: DebugLevel, format: "Hello Logging %s!", input: "World", want: "Hello Logging World!", out: &bytes.Buffer{}},
		{name: "when the global level is higher than the stdLogger level", lvl: ErrorLevel, setLvl: DisabledLevel, format: "Hello %s!", input: "World", want: "", out: &bytes.Buffer{}},
		{name: "when printing the line captured should be the client call level", lvl: ErrorLevel, setLvl: ErrorLevel, format: "Hello %s!", input: "World", want: "log_test.go", out: &bytes.Buffer{}},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%s:%s", c.lvl, c.name), func(t *testing.T) {
			l := &stdLogger{
				level: c.lvl,
				l:     log.New(c.out, c.lvl.String(), log.Lshortfile),
			}
			SetLevel(c.setLvl)
			l.Printf(c.format, c.input)
			got := c.out.String()
			assert.Contains(t, got, c.want)
		})
	}
}
func TestLogger_Dryrun_Println(t *testing.T) {
	t.SkipNow()
	l := NewStdLogger(InfoLevel)
	l.Println("Hello Dry Run!")
	debugLog := NewStdLogger(DebugLevel)
	SetLevel(DebugLevel)
	debugLog.Println("Hello Dry Run Debug!")
}
func TestLogger_Dryrun_Printf(t *testing.T) {
	t.SkipNow()
	l := NewStdLogger(InfoLevel)
	l.Println("Hello Dry Run!")
	debugLog := NewStdLogger(DebugLevel)
	SetLevel(DebugLevel)
	debugLog.Printf("Hello Dry Run Debug %s!\n", "World")
	l.Println("Log when debug level")
	setGlobalStateLevel(ErrorLevel)
	l.Println("Log when Error level")
	setGlobalStateLevel(InfoLevel)
	l.Println("Log when Error level")
}
func TestDebug(t *testing.T) {
	t.SkipNow()
	SetLevel(DebugLevel)
	Debug("Hello Debug Log!")
	SetLevel(InfoLevel)
	Debug("Will not print")
	SetLevel(ErrorLevel)
	Debug("Will not print")
}
func TestInfo(t *testing.T) {
	t.SkipNow()
	Info("Hello Info Log!")
	SetLevel(DebugLevel)
	Info("Hello Info Log in Debug Level!")
	SetLevel(DisabledLevel)
	Info("Will not print")
}
func TestError(t *testing.T) {
	t.SkipNow()
	Error("Hello Error Log!")
	SetLevel(DebugLevel)
	Error("Hello Error Log in Debug Level!")
	SetLevel(ErrorLevel)
	Error("Hello Error Log!")
	SetLevel(InfoLevel)
	Error("Hello Error Log in Info Level")
	SetLevel(DisabledLevel)
	Error("Will not print")
}
func TestDisable(t *testing.T) {
	t.SkipNow()
	SetLevel(DisabledLevel)
	Debug("Hello Debug Log!")
	Info("Hello Info Log!")
	Error("Hello Info Log in Debug Level!")
}
