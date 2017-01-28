package main

import (
	"log"
	"strconv"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	wemo "github.com/danward79/go.wemo"
)

//wemoSwitch place to put all these wemo thingys
type wemoSwitch struct {
	device    *wemo.DeviceInfo
	accessory *accessory.Switch
	transport hc.Transport
}

func createSwitch(d *wemo.DeviceInfo, pin string) (wemoSwitch, error) {
	i := accessory.Info{
		Name:         d.FriendlyName,
		SerialNumber: d.SerialNumber,
		Manufacturer: "Belkin Switch",
		Model:        d.UDN,
	}

	acc := accessory.NewSwitch(i)

	acc.Switch.On.OnValueRemoteUpdate(func(on bool) {
		if on == true {
			d.Device.On()
		} else {
			d.Device.Off()
		}
	})

	config := hc.Config{Pin: pin}
	t, err := hc.NewIPTransport(config, acc.Accessory)
	if err != nil {
		return wemoSwitch{device: d, accessory: acc, transport: t}, err //TODO: Fix error handling
	}

	go func() {
		t.Start()
	}()

	return wemoSwitch{device: d, accessory: acc, transport: t}, nil
}

func updateSwitch(subscription *wemo.SubscriptionInfo, acc *accessory.Switch) {
	b, err := strconv.ParseBool(subscription.Deviceevent.BinaryState)
	log.Println("Binary State", b, err) //TODO: Missing device can go undetected
	acc.Switch.On.SetValue(b)
}
