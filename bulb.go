package main

import (
	"fmt"
	"log"

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

		d.Device.Bulb(i.Model, "dim", fmt.Sprintf("%d", value), false)
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
