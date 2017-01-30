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
	rediscoveryPeriod := flag.Int("d", 1800, "Rediscovery period, used to look for new devices. Defaults: 30mins")
	timeout := flag.Int("t", 3600, "Wemo device subscription renewal time, 3600 second default")
	flag.Parse()

	if *rediscoveryPeriod == 0 {
		*rediscoveryPeriod = 1800
	}

	if *timeout == 0 {
		*timeout = 3600
	}

	if *listenerAddress == "" {
		log.Fatal("No IP address specified or found")
	}

	cs := make(chan wemo.SubscriptionEvent)
	go wemo.Listener(*listenerAddress, cs)
	go updateOnEvent(cs)

	//Automatic rediscovery and subscription to new devices as they appear
	timer := time.NewTimer(time.Second * time.Duration(1))
	go func() {
		for _ = range timer.C {
			timer.Reset(time.Second * time.Duration(*rediscoveryPeriod))
			log.Println("Discover Wemo devices and Subscribe")

			discover(*netInterface, *pin)
			subscribeService(*listenerAddress, cs, *timeout)
		}
	}()

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
