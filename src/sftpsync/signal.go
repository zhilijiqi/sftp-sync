package sftpsync

import (
	"os"
	"os/signal"
	"syscall"
)

func signalListen() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP)
	signal.Stop(c)

	for {
		s := <-c
		Log.Info(s)
	}
}
