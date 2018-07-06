package log

import (
	log "github.com/sirupsen/logrus"
)

// Initialize log level
func Initialize(level string) {

	switch level {
	case "debug":
		log.SetLevel(log.DebugLevel)
		break
	case "info":
		log.SetLevel(log.InfoLevel)
		break
	case "error":
		log.SetLevel(log.ErrorLevel)
		break
	case "warn":
		log.SetLevel(log.WarnLevel)
		break
	case "fatal":
		log.SetLevel(log.FatalLevel)
		break
	case "panic":
		log.SetLevel(log.PanicLevel)
		break
	}

}

// Info level
func Info(text string, arg ...interface{}) {
	log.Infof(text, arg)
}

// Error level
func Error(text string, arg ...interface{}) {
	log.Errorf(text, arg)
}

// Debug level
func Debug(text string, arg ...interface{}) {
	log.Debugf(text, arg)
}
