package log

import (
	"github.com/fatih/color"
)

var colorInfo = color.New(color.FgGreen).SprintFunc()
var colorWarn = color.New(color.FgYellow).SprintFunc()
var colorErr = color.New(color.FgRed).SprintFunc()
var colorDebug = color.New(color.FgBlue).SprintFunc()
