package main

import (
	"log"
	"os"

	"github.com/lab5e/at/n211"
)

const baudRate = 9600

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage %s <serial device>", os.Args[0])
	}

	device := n211.New(os.Args[1], baudRate)
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
