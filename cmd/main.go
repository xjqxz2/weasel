package main

import (
	"flag"
	"log"
	"peon.top/weasel/internal/service"
)

var (
	listenIPAddr       string
	listenPort         int
	networkEventNotify string
)

func main() {
	flag.StringVar(&listenIPAddr, "Host", "0.0.0.0", "The WebService Host,Default is 0.0.0.0")
	flag.IntVar(&listenPort, "Port", 8080, "The WebService Post,Default is 8080")
	flag.StringVar(&networkEventNotify, "EventNotifyDomain", "", "If device connect, this can notify")
	flag.Parse()

	srv := service.New(listenIPAddr, listenPort, networkEventNotify)

	if err := srv.Listen(); err != nil {
		log.Fatalln(err.Error())
	}
}
