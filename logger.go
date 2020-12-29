package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

const dateTimeFormat = "2006-01-02 15:04:05"

type Logger struct {
	File       *os.File
	FileLogger *logrus.Logger
	StdLogger  *logrus.Logger
}

func formatMessage(msg string) string {
	t := time.Now()
	return t.Format(dateTimeFormat+" : ") + msg + "\n"
}

func NewLogger(logFileName string) *Logger {

	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	newLogger := new(Logger)
	// File logger
	newLogger.File = file
	newLogger.FileLogger = logrus.New()
	newLogger.FileLogger.SetOutput(file)
	newLogger.FileLogger.SetFormatter(&logrus.TextFormatter{})
	newLogger.FileLogger.SetLevel(logrus.TraceLevel)

	// STDOUT logger
	newLogger.StdLogger = logrus.New()
	newLogger.StdLogger.SetOutput(os.Stdout)
	newLogger.StdLogger.SetFormatter(&logrus.TextFormatter{})
	if configuration.Debug {
		newLogger.StdLogger.SetLevel(logrus.TraceLevel)
	}

	return newLogger
}

func (l *Logger) Destroy() {
	l.File.Sync()
	l.File.Close()
}

func (l *Logger) Info(msg string) {
	l.StdLogger.Info(msg)
	l.FileLogger.Info(msg)
}

func (l *Logger) Trace(msg string) {
	l.StdLogger.Trace(msg)
	l.FileLogger.Trace(msg)
}

func (l *Logger) Error(msg string) {
	l.StdLogger.Error(msg)
	l.FileLogger.Error(msg)
}

func (l *Logger) Panic(msg string) {
	l.StdLogger.Error(msg)
	l.FileLogger.Error(msg)
	panic(msg)
}
