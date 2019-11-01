package golog

import (
	"github.com/sirupsen/logrus"
	"io"
)

type Formatter logrus.Formatter
type Fields logrus.Fields

var _ Formatter = &JSONFormatter{}

type JSONFormatter struct {
	*logrus.JSONFormatter
}

type Logrus struct {
	logger logrus.Logger
	level  Level
}

func (l *Logrus) Printf(format string, v ...interface{}) {}
func (l *Logrus) Print(v ...interface{})                 {}
func (l *Logrus) Println(v ...interface{})               {}
func (l *Logrus) Fatal(v ...interface{})                 {}
func (l *Logrus) Fatalf(format string, v ...interface{}) {}
func (l *Logrus) SetOutput(w io.Writer)                  {}
func (l *Logrus) SetFormatter(formatter Formatter)       {}
func (l *Logrus) WithFields(fields Fields) Logger {
	return l
}
