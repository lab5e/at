package nrf91

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/lab5e/at"
)

// Note: There is only a single socket in the LTE modem. The port parameter may be 0, then the
// socket won't be bound to a port.
func (d *nrf91) CreateUDPSocket(port int) (int, error) {
	// Parameters:
	// #1: 0 - close, 1 - open ipv4, 2 - open ipv6
	// #2: 1 - TCP, 2 - UDP
	// #3: 0 - client, 1- server

	err := d.cmd.Transact("AT#XSOCKET=1,2,0", func(s string) error {
		// This will just return "OK"
		return nil
	})
	if err != nil {
		return 0, err
	}

	// Bind to a port if a port is set
	if port != 0 {
		err = d.cmd.Transact(fmt.Sprintf("AT#XBIND=%d", port), func(s string) error {
			// This returns just OK
			return nil
		})
	}
	return 1, err
}

func (d *nrf91) SendUDP(socket int, address net.IP, remotePort int, data []byte) (int, error) {
	if socket != 1 {
		return 0, errors.New("unknown socket ID")
	}
	var bytesSent = 0
	var err error
	err = d.cmd.Transact(
		fmt.Sprintf(`AT#XSENDTO="%s",%d,0,"%s"`, address.String(), remotePort, hex.EncodeToString(data)),
		func(s string) error {
			if strings.TrimSpace(s) == "" {
				return nil
			}
			parts := strings.Split(s, ":")
			if len(parts) == 2 {
				bytesSent, err = strconv.Atoi(strings.TrimSpace(parts[1]))
				if err != nil {
					log.Printf("Could not parse byte count from %s", parts[1])
					return errors.New("could not parse number of bytes")
				}
				return nil
			}
			log.Printf("Could not parse response from modem: %s", s)
			return errors.New("uknown response")
		})
	return bytesSent, err
}

// Note: Receive has a 10 second timeout
func (d *nrf91) ReceiveUDP(socket int, length int) (*at.ReceivedData, error) {
	var data at.ReceivedData

	err := d.cmd.Transact("AT#XRECVFROM=10", func(s string) error {
		// We'll recive at least two lines - first the data, then a line with #XRECVFROM=<size>,"<ip>"
		if strings.HasPrefix(s, "#XRECVFROM:") {
			parts := strings.Split(s, "FROM:")
			if len(parts) != 2 {
				return errors.New("could not parse recvfrom response")
			}
			sizeAddr := strings.Split(parts[1], ",")
			if len(sizeAddr) != 2 {
				return errors.New("could not parse size and addr in recvfrom response")
			}
			var err error
			data.Length, err = strconv.Atoi(sizeAddr[0])
			if err != nil {
				return err
			}
			data.IP = at.TrimQuotes(sizeAddr[1])
			return nil
		}
		data.Data = []byte(s)
		data.Socket = 1
		return nil
	})

	return &data, err
}

func (d *nrf91) CloseUDPSocket(socket int) error {
	return d.cmd.Transact("AT#XSOCKET=0", nil)
}
