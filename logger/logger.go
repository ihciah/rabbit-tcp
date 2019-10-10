package logger

import (
	"log"
	"os"
)

const (
	LogLevelOff = iota
	LogLevelFatal
	LogLevelError
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

var LEVEL int = LogLevelOff

type Logger struct {
	logger *log.Logger
	level  int
}

func NewLogger(prefix string) *Logger {
	return &Logger{
		logger: log.New(os.Stdout, prefix, log.LstdFlags),
		level:  LEVEL,
	}
}

func (l *Logger) Debugln(v string) {
	if l.level >= LogLevelDebug {
		l.logger.Println("[Debug] " + v)
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.level >= LogLevelDebug {
		l.logger.Printf("[Debug] "+format, v...)
	}
}

func (l *Logger) Infoln(v string) {
	if l.level >= LogLevelInfo {
		l.logger.Println("[Info] " + v)
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if l.level >= LogLevelInfo {
		l.logger.Printf("[Info] "+format, v...)
	}
}

func (l *Logger) Warnln(v string) {
	if l.level >= LogLevelWarn {
		l.logger.Println("[Warn] " + v)
	}
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.level >= LogLevelWarn {
		l.logger.Printf("[Warn] "+format, v...)
	}
}

func (l *Logger) Errorln(v string) {
	if l.level >= LogLevelError {
		l.logger.Println("[Error] " + v)
	}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.level >= LogLevelError {
		l.logger.Printf("[Error] "+format, v...)
	}
}

func (l *Logger) Fatalln(v string) {
	if l.level >= LogLevelFatal {
		l.logger.Println("[Fatal] " + v)
	}
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	if l.level >= LogLevelFatal {
		l.logger.Printf("[Fatal] "+format, v...)
	}
}
