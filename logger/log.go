package logger

import (
	log "github.com/sirupsen/logrus"
)

var (
	LogInfo  = log.Infof
	LogWarn  = log.Warnf
	LogPanic = log.Panicf
)
