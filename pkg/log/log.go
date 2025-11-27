package log

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
}

type Record struct {
	message string
}

func Warrningf(format string, args ...any) {
	s := fmt.Sprintf(format, args...)
	record := Record{
		s,
	}
	logrus.WithFields(logrus.Fields{}).Warning(record.message)
}

func Warrning(s string) {
	record := Record{
		s,
	}
	logrus.WithFields(logrus.Fields{}).Warning(record.message)
}

func Infof(format string, args ...any) {
	s := fmt.Sprintf(format, args...)
	record := Record{
		s,
	}
	logrus.WithFields(logrus.Fields{}).Info(record.message)
}

func Info(s string) {
	record := Record{
		s,
	}
	logrus.WithFields(logrus.Fields{}).Info(record.message)
}

func Debugging(s string) {
	record := Record{
		s,
	}
	logrus.WithFields(logrus.Fields{}).Debug(record.message)
}

func Debuggingf(format string, args ...any) {
	s := fmt.Sprintf(format, args...)
	record := Record{
		s,
	}
	logrus.WithFields(logrus.Fields{}).Debug(record.message)
}
