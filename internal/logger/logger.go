package logger

import (
	"log"
	"sync"
)

type Logger interface {
	Infof(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}

var once *sync.Once = &sync.Once{}
var logger Logger

func InitLogger(l Logger) {
	once.Do(func() {
		logger = l
	})
}

func handlePanicf(format string, v ...interface{}) {
	if a := recover(); a != nil {
		// ignore any panic errors in logger package and log to default log.
		log.Printf(format, v...)
	}
}

func Infof(format string, v ...interface{}) {
	defer handlePanicf(format, v...)
	if logger != nil {
		logger.Infof(format, v...)
	}
}

func Errorf(format string, v ...interface{}) {
	defer handlePanicf(format, v...)
	if logger != nil {
		logger.Errorf(format, v...)
	}
}

func Debugf(format string, v ...interface{}) {
	defer handlePanicf(format, v...)
	if logger != nil {
		logger.Debugf(format, v...)
	}
}

func Fatalf(format string, v ...interface{}) {
	defer handlePanicf(format, v...)
	if logger != nil {
		logger.Fatalf(format, v...)
	}
}
