package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lab5e/at"
	"github.com/lab5e/at/bg95"
	"github.com/lab5e/at/n211"
	"github.com/lab5e/at/nrf91"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage %s <nrf91|bg95| <serial device>", os.Args[0])
	}
	var device at.Device
	switch strings.ToLower(os.Args[1]) {
	case "nrf91":
		device = nrf91.New(os.Args[2], nrf91.DefaultBaudRate)
	case "bg95":
		device = bg95.New(os.Args[2], bg95.DefaultBaudRate)
	case "n211":
		device = n211.New(os.Args[2], n211.DefaultBaudRate)
	default:
		log.Fatalf("Unknown type: %s", os.Args[1])
	}
	if err := device.Start(); err != nil {
		log.Fatalf("Error opening device: %v", err)
	}
	defer device.Close()

	imsi, err := device.GetIMSI()
	if err != nil {
		log.Fatalf("Error getting IMSI: %v", err)
	}

	imei, err := device.GetIMEI()
	if err != nil {
		log.Fatalf("Error getting IMEI: %v", err)
	}

	// ICCID might not be supported on all versions
	ccid, err := device.GetCCID()
	if err != nil {
		log.Printf("Error getting CCID: %v", err)
	}

	fmt.Printf("IMSI = '%s'\n", imsi)
	fmt.Printf("IMEI = '%s'\n", imei)
	fmt.Printf("CCID = '%s'\n", ccid)
}
