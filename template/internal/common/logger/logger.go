package logger

import (
	"github.com/sirupsen/logrus"
)

var _LoggerEntry *logrus.Entry

func SetGlobalLogger(l *logrus.Logger, sname string) {
	_LoggerEntry = l.WithField("service", sname)
}

func GetGlobalLogger() *logrus.Entry {
	return _LoggerEntry
}
