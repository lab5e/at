package main

import (
	"log"
	"os"
	"strings"

	"github.com/lab5e/at"
	"github.com/lab5e/at/n211"
	"github.com/lab5e/at/nrf91"
)

const baudRate = 9600

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage %s <n211|nrf91|bg95> <serial device>", os.Args[0])
	}
	deviceType := os.Args[1]
	serialDevice := os.Args[2]

	var device at.Device

	switch strings.ToLower(deviceType) {
	case "n211":
		device = n211.New(serialDevice, n211.DefaultBaudRate)
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
	device.SetDebug(true)

	// Just send a blank AT command to verify the device is there
	if err := device.AT(); err != nil {
		log.Fatalf("Error speaking to device on '%s': %v", os.Args[1], err)
	}
	log.Printf("Device seems to be responsive")
}
