package golog

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

type Formatter logrus.Formatter
type Fields logrus.Fields

var _ Formatter = &JSONFormatter{}

type JSONFormatter struct {
	*logrus.JSONFormatter
}

func NewLogrusLogger(level Level) *Logrus {
	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)

	var llevel logrus.Level

	switch level {
	case DebugLevel:
		llevel = logrus.DebugLevel
	case TraceLevel:
		llevel = logrus.TraceLevel
	case InfoLevel:
		llevel = logrus.InfoLevel
	case WarningLevel:
		llevel = logrus.WarnLevel
	case ErrorLevel:
		llevel = logrus.ErrorLevel
	case DisabledLevel:
		llevel = logrus.InfoLevel
	}

	return &Logrus{
		logger:      l,
		level:       level,
		logrusLevel: llevel,
	}
}

type Logrus struct {
	logger      *logrus.Logger
	level       Level
	logrusLevel logrus.Level
}

func (l *Logrus) Printf(format string, v ...interface{}) {}
func (l *Logrus) Print(v ...interface{}) {
	if l.isEnabled() {
		l.logger.Log(l.logrusLevel, v...)
	}
	return
}
func (l *Logrus) Println(v ...interface{}) {
	if l.isEnabled() {
		l.logger.Log(l.logrusLevel, v...)
	}
}
func (l *Logrus) Fatal(v ...interface{}) {
	if l.isEnabled() {
		l.logger.Log(l.logrusLevel, v...)
		os.Exit(1)
	}
	return
}
func (l *Logrus) Fatalf(format string, v ...interface{}) {
	if l.isEnabled() {
		l.logger.Logf(l.logrusLevel, format, v...)
	}
}
func (l *Logrus) SetOutput(w io.Writer) {
	l.logger.SetOutput(w)
}
func (l *Logrus) SetFormatter(formatter Formatter) {
	l.logger.SetFormatter(formatter)
}
func (l *Logrus) WithFields(fields Fields) Logger {
	return l
}

func (l *Logrus) isEnabled() bool {
	gstate := getState()
	if l.level < gstate.currentLevel {
		return false
	}
	return true
}
