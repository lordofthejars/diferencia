package log

import (
	"log"
)

func Info(text string, arg ...interface{}) {
	text = colorInfo("[INFO] ") + text
	log.Printf(text, arg...)
}

func Error(text string, arg ...interface{}) {
	text = colorErr("[ERROR] ") + text
	log.Printf(text, arg...)
}

func Debug(text string, arg ...interface{}) {
	text = colorDebug("[DEBUG] ") + text
	log.Printf(text, arg...)
}
