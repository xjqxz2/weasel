package main

import (
	"log"
	"peon.top/weasel/service"
)

func main() {
	if err := service.New("", 8080).Listen(); err != nil {
		log.Fatalf(err.Error())
	}
}
