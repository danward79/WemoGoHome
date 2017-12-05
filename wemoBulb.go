package main

import (
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
)

// CustomLightbulb none standard Wemo homekit device.
type CustomLightbulb struct {
	*service.Service

	On         *characteristic.On
	Brightness *characteristic.Brightness
}

// NewCustomLightbulb ...
func NewCustomLightbulb() *CustomLightbulb {
	svc := CustomLightbulb{}
	svc.Service = service.New(service.TypeLightbulb)

	svc.On = characteristic.NewOn()
	svc.AddCharacteristic(svc.On.Characteristic)

	svc.Brightness = characteristic.NewBrightness()
	svc.AddCharacteristic(svc.Brightness.Characteristic)

	return &svc
}

// Lightbulb ...
type Lightbulb struct {
	*accessory.Accessory
	Lightbulb *CustomLightbulb
}

// NewLightbulb returns a light bulb accessory which one light bulb service.
func NewLightbulb(info accessory.Info) *Lightbulb {
	acc := Lightbulb{}
	acc.Accessory = accessory.New(info, accessory.TypeLightbulb)
	acc.Lightbulb = NewCustomLightbulb()

	acc.Lightbulb.Brightness.SetValue(100)

	acc.AddService(acc.Lightbulb.Service)

	return &acc
}
