package main

import "peon.top/weasel/service"

func main() {
	service.New("", 8080).Listen()
}
