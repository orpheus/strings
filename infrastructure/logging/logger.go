package logging

import "log"

type Logger interface {
	Log(v ...interface{})
	Logf(format string, v ...interface{})
}

type TmpLogger struct{}

func (l *TmpLogger) Log(v ...interface{}) {
	log.Println(v)
}

func (l *TmpLogger) Logf(format string, v ...interface{}) {
	log.Printf(format, v)
}
