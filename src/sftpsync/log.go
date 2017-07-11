package sftpsync

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	Debug = iota
	Info
	Warn
	Error
)

var Log = NewLog(Debug)

type logger struct {
	level  int
	logger *log.Logger
}

func NewLog(level int) *logger {
	var f, err = os.OpenFile(filepath.Join(baseDir, "sftp_sync.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		f = os.Stderr
	}
	return &logger{level: level, logger: log.New(f, "", log.LstdFlags)}
}
func (l *logger) Debug(v ...interface{}) {
	if l.level <= Debug {
		l.logger.Output(3, fmt.Sprintln(v...))
	}
}

func (l *logger) Info(v ...interface{}) {
	if l.level <= Info {
		l.logger.Output(3, fmt.Sprintln(v...))
	}
}

func (l *logger) Error(v ...interface{}) {
	if l.level <= Error {
		l.logger.Output(3, fmt.Sprintln(v...))
	}
}

func (l *logger) Fatal(v ...interface{}) {
	l.logger.Output(3, fmt.Sprintln(v...))
	os.Exit(1)
}
