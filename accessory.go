package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/net/context"

	"github.com/brutella/hc"
	wemo "github.com/danward79/go.wemo"
)

var (
	wemoThings = make(map[string]interface{})
	cs         = make(chan wemo.SubscriptionEvent)
)

func discover(i, pin string) {
	ctx := context.Background()

	api, _ := wemo.NewByInterface(i)

	devices, _ := api.DiscoverAll(5 * time.Second)
	for _, device := range devices {
		deviceInfo, _ := device.FetchDeviceInfo(ctx)
		//log.Println(deviceInfo) //TODO: Sometimes all devices are not found. Add periodic rescan?
		go createAccessory(deviceInfo, pin)
	}

	terminateAccesories()
}

func createAccessory(d *wemo.DeviceInfo, pin string) {
	var err error

	switch d.DeviceType {
	case wemo.Controllee:
		wemoThings[d.UDN], err = createSwitch(d, pin)
		if err != nil {
			log.Println(err)
		}
	case wemo.Bridge:
		for k, v := range d.EndDevices.EndDeviceInfo {
			wemoThings[v.DeviceID], err = createBulb(d, k, pin)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func terminateAccesories() {
	hc.OnTermination(func() {
		for _, thing := range wemoThings {
			switch thing.(type) {
			case wemoSwitch:
				thing.(wemoSwitch).transport.Stop()
			case wemoBulb:
				thing.(wemoBulb).transport.Stop()
			}
		}

		time.Sleep(100 * time.Millisecond)
		os.Exit(1)
	})
}

func getDevice(thing interface{}) *wemo.DeviceInfo {
	switch thing.(type) {
	case wemoSwitch:
		return thing.(wemoSwitch).device
	case wemoBulb:
		return thing.(wemoBulb).device
	}
	return nil
}

func subscribeService(listenerAddress string, subsCh chan wemo.SubscriptionEvent) {

	subscriptions := make(map[string]*wemo.SubscriptionInfo)

	for _, thing := range wemoThings {
		fmt.Println("subService", thing)
		subscribe(getDevice(thing), listenerAddress, subscriptions)
	}

	go wemo.Listener(listenerAddress, subsCh)

	for m := range subsCh {
		if _, ok := subscriptions[m.Sid]; ok {
			subscriptions[m.Sid].State = m.State
			log.Println("---Subscriber Event: ", subscriptions[m.Sid])
		} else {
			log.Println("Does'nt exist, ", m.Sid)
		}
	}
}

func subscribe(d *wemo.DeviceInfo, listenerAddress string, subscriptions map[string]*wemo.SubscriptionInfo) {

	_, err := d.Device.ManageSubscription(listenerAddress, 300, subscriptions)
	if err != 200 {
		log.Println("Initial Error Subscribing: ", err)
	}
}
