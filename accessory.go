package main

import (
	"log"
	"os"
	"time"

	"golang.org/x/net/context"

	"github.com/brutella/hc"
	wemo "github.com/danward79/go.wemo"
)

var (
	wemoThings    = make(map[string]interface{})
	cs            = make(chan wemo.SubscriptionEvent)
	subscriptions = make(map[string]*wemo.SubscriptionInfo)
)

func discover(i, pin string) {
	ctx := context.Background()

	api, _ := wemo.NewByInterface(i)

	devices, _ := api.DiscoverAll(5 * time.Second)
	for _, device := range devices {
		deviceInfo, _ := device.FetchDeviceInfo(ctx)
		createAccessory(deviceInfo, pin)
	}

	terminateAccesories()
}

func createAccessory(d *wemo.DeviceInfo, pin string) {
	var err error
	log.Println("d.DeviceType", d.DeviceType)

	switch d.DeviceType {
	case wemo.Controllee:

		if _, exists := wemoThings[d.UDN]; !exists {
			wemoThings[d.UDN], err = createSwitch(d, pin)
			if err != nil {
				log.Println(err)
			}
		}

	case wemo.Bridge:

		for k, v := range d.EndDevices.EndDeviceInfo {
			if _, exists := wemoThings[v.DeviceID]; !exists {
				wemoThings[v.DeviceID], err = createBulb(d, k, pin)
				if err != nil {
					log.Println(err)
				}
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

func subscriptionExists(d *wemo.DeviceInfo, subscriptions map[string]*wemo.SubscriptionInfo) bool {
	for _, v := range subscriptions {
		if v.DeviceInfo.UDN == d.UDN {
			return true
		}
	}
	return false
}

func updateAccessory(subscription *wemo.SubscriptionInfo) {

	for _, v := range wemoThings {
		if getDevice(v).UDN == subscription.DeviceInfo.UDN {

			switch v.(type) {
			case wemoSwitch:
				updateSwitch(subscription, v.(wemoSwitch).accessory)
			case wemoBulb:
				if v.(wemoBulb).endDevice.DeviceID == subscription.Deviceevent.StateEvent.DeviceID {
					updateBulb(subscription, v.(wemoBulb).accessory)
				}
			}

		}
	}
}

func subscribeService(listenerAddress string, subsCh chan wemo.SubscriptionEvent) {
	for _, thing := range wemoThings {
		d := getDevice(thing)
		if !subscriptionExists(d, subscriptions) {
			subscribe(d, listenerAddress, subscriptions)
		}
	}
}

func subscribe(d *wemo.DeviceInfo, listenerAddress string, subscriptions map[string]*wemo.SubscriptionInfo) {
	_, err := d.Device.ManageSubscription(listenerAddress, 300, subscriptions)
	if err != 200 {
		log.Println("Initial Error Subscribing: ", err)
	}
}

func updateOnEvent(subsCh chan wemo.SubscriptionEvent) {
	for m := range subsCh {
		if _, ok := subscriptions[m.Sid]; ok {
			subscriptions[m.Sid].Deviceevent = m.Deviceevent
			log.Println("Event:", m.Deviceevent)
			updateAccessory(subscriptions[m.Sid])
		}
	}
}
