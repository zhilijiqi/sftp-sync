package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	ss "sftpsync"
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

	flag.StringVar(&logLevel, "log", "debug", "log Level(debug,info,warn,error)")
	flag.IntVar(&interval, "i", 60, "interval in seconds")

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
	fmt.Println("running")

	req := make(chan struct{}, 1)
	done := make(chan struct{}, 1)
	i := time.Duration(interval) * time.Second
	timer := time.NewTimer(i)
	for {
		<-timer.C
		req <- struct{}{}
		go ss.Syncd(conf, req, done)
		<-done
		timer.Reset(i)
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
