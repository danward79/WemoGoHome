### Wemo Homekit Bridge

Simple Go app which bridges [Wemo](http://www.wemo.com) devices to HomeKit and makes them available in iOS devices, such that they can be used by HomeKit, in apps like Apple's [Home](http://www.apple.com/au/ios/home/).

#### Status
App is being developed but is stable. It is currently deployed on my Raspberry P and has been operation for 2 months.

#### Stuff to do
- Sometimes if a device is removed from the network, it is not detected and shown in Homekit for sometime.
- ~~Wemo bulbs have no colour setting, this characteristic needs removing.~~
- ~~Status incorrectly shown upon discoery.~~

#### Install
```
go get -u github.com/danward79/WemoGoHome
cd $GOPATH/src/github.com/danward79/WemoGoHome/
go install
```

#### Run
This will run with default settings, if you go bin is in your $PATH

```
WemoGoHome
```

#### Or Build and Copy
If you want to cross compile for installing manually on a Raspberry pi or some other machine.

*Note* Remove GOARM=7 if you are unsure of the arm device type... e.g. Older Raspberry Pi devices.

```
go get -u github.com/danward79/WemoGoHome
cd $GOPATH/src/github.com/danward79/WemoGoHome/
GOOS=linux GOARCH=arm GOARM=7 go build
```

Copy to the location you wish to use as your working directory and
```
cd ~/pathtoyourworkingdirectory
./WemoGoHome
```

#### Command Options
The app has a few command params that are available to tweak, the defaults should be ok for most

```
Usage of ./WemoGoHome:
  -d int
    	Rediscovery period, used to look for new devices. Defaults: 30mins (default 1800)
  -i string
    	Network Interface, default en0 (default "en0")
  -l string
    	Listener address (default "192.168.1.22:6767")
  -pin string
    	Pin number for wemo accessories, default 87654321 (default "87654321")
```

#### Libraries
This app uses two libraries outside of the standard, these are a fork of the [go.wemo](https://github.com/danward79/go.wemo) package, which I and others have updated heavily. The other package is the [HomeKit](https://github.com/brutella/hc) library by brutella.
