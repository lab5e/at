package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/lab5e/at"
	"github.com/lab5e/at/bg95"
	"github.com/lab5e/at/n211"
	"github.com/lab5e/at/nrf91"
)

func main() {
	var deviceType, serialDevice string
	var port int
	var debug bool
	flag.StringVar(&deviceType, "device", "nrf91", "Device type")
	flag.StringVar(&serialDevice, "serial", "/dev/serial", "Serial device")
	flag.IntVar(&port, "port", 0, "Local port")
	flag.BoolVar(&debug, "debug", false, "Show debug messages")
	flag.Parse()
	if port == 0 {
		log.Fatalf("Must specify a port")
	}
	var device at.Device

	switch strings.ToLower(deviceType) {
	case "n211":
		device = n211.New(serialDevice, n211.DefaultBaudRate)
	case "bg95":
		device = bg95.New(serialDevice, bg95.DefaultBaudRate)
	case "nrf91":
		device = nrf91.New(serialDevice, nrf91.DefaultBaudRate)
	default:
		log.Fatalf("Unknown device type: %s", deviceType)
	}

	if err := device.Start(); err != nil {
		log.Fatalf("Error opening device: %v", err)
	}
	defer device.Close()

	// Turn on debugging so you can see the interaction with the device
	device.SetDebug(debug)

	if err := device.SetRadio(true); err != nil {
		log.Fatalf("Could not enable radio: %v", err)
	}

	imsi, err := device.GetIMSI()
	if err != nil {
		log.Fatalf("Got IMSI %s", imsi)
	}
	socket, err := device.CreateUDPSocket(port)
	if err != nil {
		log.Fatalf("Could not create UDP socket: %v", err)
	}
	defer func() {
		device.CloseUDPSocket(socket)
		log.Printf("Closed UDP socket")
	}()

	log.Printf("Waiting for data.. Running on local port %d with IMSI %s", port, imsi)

	for {
		data, err := device.ReceiveUDP(socket, 128)
		if err != nil {
			log.Printf("Got error receiving data: %v. Retry in 10 seconds", err)
			time.Sleep(10 * time.Second)
			continue
		}
		log.Printf("Recevied %d bytes from %s:%d: %v", data.Length, data.IP, data.Port, string(data.Data))
	}
}
