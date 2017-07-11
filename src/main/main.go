package main

import (
	"fmt"
	"os"
	"time"

	ss "sftpsync"
)

func main() {
	//	var confFile, cmdServer, cmdLocal string
	//	var cmdConf ss.Config
	//	var printVer bool

	//	flag.BoolVar(&printVer, "version", false, "print version")
	//	flag.StringVar(&confFile, "c", "config.json", "specify config file")
	//	flag.StringVar(&cmdServer, "s", "", "server address")
	//	flag.StringVar(&cmdLocal, "b", "", "local address, listen only to this address if specified")
	//	flag.StringVar(&cmdConfig.Password, "k", "", "password")
	//	flag.IntVar(&cmdConfig.ServerPort, "p", 0, "server port")
	//	flag.IntVar(&cmdConfig.Timeout, "t", 300, "timeout in seconds")
	//	flag.IntVar(&cmdConfig.LocalPort, "l", 0, "local socks5 proxy port")
	//	flag.StringVar(&cmdConfig.Method, "m", "", "encryption method, default: aes-256-cfb")
	//	flag.BoolVar((*bool)(&debug), "d", false, "print debug message")

	//	flag.Parse()

	//	if printVer {
	//		ss.PrintVersion()
	//		os.Exit(0)
	//	}

	args := os.Args
	conf, err := ss.NewConfig(args)
	if err != nil {
		ss.Log.Fatal(err)
	}
	ss.Log.Info("conf info:", conf)
	ss.GetPid()
	//start(conf)
	delay()
}
func start(conf *ss.Config) {
	req := make(chan struct{}, 1)
	done := make(chan struct{}, 1)
	timer := time.NewTimer(time.Second)
	for {
		<-timer.C
		req <- struct{}{}
		go ss.Syncd(conf, req, done)
		<-done
		timer.Reset(time.Second)
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
