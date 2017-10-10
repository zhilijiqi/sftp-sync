package main

import (
	"flag"
	"fmt"
	"os"
	ss "sftpsync"
	"time"
)

func main() {

	var confFile, logLevel string
	//var cmdConf ss.Config
	var printVer bool
	var interval int

	flag.BoolVar(&printVer, "version", false, "print version")
	flag.StringVar(&confFile, "c", "config.json", "specify config file")
	//	flag.StringVar(&cmdServer, "h", "", "server address and port(host:port)")
	//	flag.StringVar(&cmdLocalPath, "l", "", "local dir")
	//	flag.StringVar(&cmdServerPath, "s", "", "sync sftp server dir")
	//	flag.StringVar(&cmdConfig.Password, "k", "", "password")
	//	flag.IntVar(&cmdConfig.ServerPort, "p", 0, "server port")

	flag.StringVar(&logLevel, "log", "info", "log Level(debug,info,warn,error)")
	flag.IntVar(&interval, "i", 2, "interval in seconds")

	//  flag.BoolVar((*bool)(&debug), "d", false, "print debug message")

	flag.Parse()

	if printVer {
		ss.PrintVersion()
		os.Exit(0)
	}

	ss.LogLevelByName(logLevel)

	conf, err := ss.NewConfig(confFile)
	if err != nil {
		ss.Log.Fatal(err)
	}
	ss.Log.Info("conf info:", conf)

	ss.GetPid()
	run(interval, conf)
}
func run(interval int, conf *ss.Config) {
	ss.Log.Info("running")

	fs, err := ss.NewLocalFs(conf.LocalPath)
	if err != nil {
		ss.Log.Error(err)
	}
	ss.Log.Info("add cache dir:", conf.LocalPath)

	req := make(chan struct{}, 1)
	done := make(chan struct{}, 1)
	go ss.Syncd(conf, req, done)

	i := time.Duration(interval) * time.Second
	timer := time.NewTimer(i)

	var counter int64 = 1
	for {
		<-timer.C
		upFs(fs, counter, i)
		req <- struct{}{}
		//go ss.Syncd(conf, req, done)
		<-done
		counter++
		timer.Reset(i)
	}
}

func upFs(fs *ss.LocalFs, counter int64, i time.Duration) {
	if counter%5 == 0 || int64(i)*counter >= time.Hour.Nanoseconds() {
		if err := fs.Reload(); err != nil {
			ss.Log.Error(err)
		}
	}
}
func delay() {
	ticker := time.NewTicker(time.Minute * 1)
	go func() {
		for _ = range ticker.C {
			fmt.Printf("ticked at %v\n", time.Now())
		}
	}()
	for {
		time.Sleep(time.Hour)
	}
}
