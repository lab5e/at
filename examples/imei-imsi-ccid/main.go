package main

import (
	"fmt"
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

	imsi, err := device.GetIMSI()
	if err != nil {
		log.Fatalf("Error getting IMSI: %v", err)
	}

	imei, err := device.GetIMEI()
	if err != nil {
		log.Fatalf("Error getting IMEI: %v", err)
	}

	ccid, err := device.GetCCID()
	if err != nil {
		log.Fatalf("Error getting CCID: %v", err)
	}

	fmt.Printf("IMSI = '%s'\n", imsi)
	fmt.Printf("IMEI = '%s'\n", imei)
	fmt.Printf("CCID = '%s'\n", ccid)
}
