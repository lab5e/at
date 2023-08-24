package main

import (
	"log"
	"os"
	"strings"

	"github.com/lab5e/at"
	"github.com/lab5e/at/bg95"
	"github.com/lab5e/at/n211"
	"github.com/lab5e/at/nrf91"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatalf("Usage %s <n211|nrf91|bg95> <serial device> <apn>", os.Args[0])
	}
	deviceType := os.Args[1]
	serialDevice := os.Args[2]
	apn := os.Args[3]

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

	device.SetDebug(true)

	err := device.SetAPN(apn)
	if err != nil {
		log.Fatal(err)
	}
}
