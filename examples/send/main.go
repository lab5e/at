package main

import (
	"flag"
	"log"
	"strings"

	"net"

	"github.com/lab5e/at"
	"github.com/lab5e/at/bg95"
	"github.com/lab5e/at/n211"
	"github.com/lab5e/at/nrf91"
)

func main() {
	var deviceType, serialDevice, ip, message string
	var port int
	var debug bool
	flag.StringVar(&deviceType, "device", "nrf91", "Device type")
	flag.StringVar(&serialDevice, "serial", "/dev/serial", "Serial device")
	flag.StringVar(&ip, "ip", "172.16.15.14", "IP address")
	flag.IntVar(&port, "port", 0, "Server port")
	flag.StringVar(&message, "message", "", "Message to send")
	flag.BoolVar(&debug, "debug", false, "Show debug messages")
	flag.Parse()
	if len(message) == 0 {
		log.Fatalf("Needs a message to send")
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

	socket, err := device.CreateUDPSocket(0)
	if err != nil {
		log.Fatalf("Could not create UDP socket: %v", err)
	}
	defer device.CloseUDPSocket(socket)

	n, err := device.SendUDP(socket, net.ParseIP(ip), port, []byte(message))
	if err != nil {
		log.Fatalf("Error sending UDP: %v", err)
	}
	if n != len(message) {
		log.Printf("Device sent %d bytes but expected %d", n, len(message))
	}
	log.Printf("Message sent to %s:%d", ip, port)
}
