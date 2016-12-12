package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	wemo "github.com/danward79/go.wemo"
)

func main() {
	pin := flag.String("pin", "87654321", "Pin number for wemo accessories, default 87654321")
	netInterface := flag.String("i", "en0", "Network Interface, default en0")
	listenerAddress := flag.String("l", getIPAddress()+":6767", "Listener address")
	flag.Parse()

	if *listenerAddress == "" {
		log.Fatal("No IP address specified or found")
	}

	discover(*netInterface, *pin)

	cs := make(chan wemo.SubscriptionEvent)
	subscribeService(*listenerAddress, cs)

	waitToExit()
}

func waitToExit() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	<-c
	time.Sleep(200 * time.Millisecond)
	os.Exit(1)
}

func getIPAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal("Error finding network IP:", err)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
