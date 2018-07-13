package log

import (
	"github.com/sirupsen/logrus"
)

// Initialize log level
func Initialize(level string) {

	switch level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
		break
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
		break
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
		break
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
		break
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
		break
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
		break
	}

}
