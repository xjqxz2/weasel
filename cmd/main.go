package main

import (
	"flag"
	"log"
	"peon.top/weasel/service"
)

var (
	listenIPAddr string
	listenPort   int
)

func main() {
	flag.StringVar(&listenIPAddr, "Host", "0.0.0.0", "The WebService Host,Default is 0.0.0.0")
	flag.IntVar(&listenPort, "Port", 8080, "The WebService Post,Default is 8080")
	flag.Parse()

	srv := service.New(listenIPAddr, listenPort)

	if err := srv.Listen(); err != nil {
		log.Fatalln(err.Error())
	}
}
