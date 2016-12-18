package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	wemo "github.com/danward79/go.wemo"
)

//wemoBulb place to put all these wemo thingys
type wemoBulb struct {
	device    *wemo.DeviceInfo
	endDevice *wemo.EndDeviceInfo
	accessory *accessory.Lightbulb
	transport hc.Transport
}

func createBulb(d *wemo.DeviceInfo, index int, pin string) (wemoBulb, error) {
	i := accessory.Info{
		Name:         d.EndDevices.EndDeviceInfo[index].FriendlyName,
		Manufacturer: "Belkin Bulb",
		Model:        d.EndDevices.EndDeviceInfo[index].DeviceID,
	}

	acc := accessory.NewLightbulb(i)

	acc.Lightbulb.Brightness.OnValueRemoteUpdate(func(value int) {
		if value < 0 || value > 255 {
			log.Println("Dim Value out of bounds:", value)
			return
		}

		d.Device.Bulb(i.Model, "dim", fmt.Sprintf("%d", value/100*255), false)
	})

	acc.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
		if on == true {
			d.Device.Bulb(i.Model, "on", "", false)
		} else {
			d.Device.Bulb(i.Model, "off", "", false)
		}
	})

	config := hc.Config{Pin: pin}
	t, err := hc.NewIPTransport(config, acc.Accessory)
	if err != nil {
		return wemoBulb{device: d, endDevice: &d.EndDevices.EndDeviceInfo[index], accessory: acc, transport: t}, err //TODO: Fix error handling
	}

	go func() {
		t.Start()
	}()

	return wemoBulb{device: d, endDevice: &d.EndDevices.EndDeviceInfo[index], accessory: acc, transport: t}, nil
}

func updateBulb(subscription *wemo.SubscriptionInfo, acc *accessory.Lightbulb) {
	switch subscription.Deviceevent.StateEvent.CapabilityID {
	case "10006":
		b, _ := strconv.ParseBool(subscription.Deviceevent.StateEvent.Value)
		acc.Lightbulb.On.SetValue(b)

	case "10008":
		s := strings.Split(subscription.Deviceevent.StateEvent.Value, ":")
		i, _ := strconv.ParseInt(s[0], 10, 0)
		level := int(float32(i) / 255 * 100)

		acc.Lightbulb.On.SetValue(true)
		acc.Lightbulb.Brightness.SetValue(level)

		if i < 1 {
			acc.Lightbulb.On.SetValue(false)
		}
	}
}
