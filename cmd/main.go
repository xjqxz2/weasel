package main

import (
	"flag"
	"log"
	"peon.top/weasel/internal/service"
	"peon.top/weasel/internal/weasel"
)

var (
	listenIPAddr       string
	listenPort         int
	networkEventNotify string
	enableKeeper       bool
)

func main() {
	flag.StringVar(&listenIPAddr, "Host", "0.0.0.0", "The WebService Host,Default is 0.0.0.0")
	flag.IntVar(&listenPort, "Port", 8080, "The WebService Post,Default is 8080")
	flag.StringVar(&networkEventNotify, "EventNotifyDomain", "", "If device connect, this can notify")
	flag.BoolVar(&enableKeeper, "EnableKeeper", false, "Is enable message keeper")
	flag.Parse()

	//	加载必要的组件
	event, keeper := loadExtendComponent()

	//	创建服务
	srv := service.New(listenIPAddr, listenPort, event, keeper)

	//	监听服务
	if err := srv.Listen(); err != nil {
		log.Fatalln(err.Error())
	}
}

func loadExtendComponent() (keeper weasel.Keeper, event weasel.Event) {
	keeper, event = weasel.NewNoMessageKeeper(), weasel.NewNoNotifyEvent()

	if networkEventNotify != "" {
		event = weasel.NewRemoteEvent(networkEventNotify)
	}

	if enableKeeper {
		keeper = weasel.NewLRUKeeper()
	}

	return keeper, event
}
