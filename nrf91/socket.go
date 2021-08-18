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

func (d *nrf91) CreateUDPSocket(port int) (int, error) {
	// Section 16.24.3 states that this port is reserved
	if port == 5683 {
		return -1, errors.New("reserved port value")
	}

	cmd := "AT+NSOCR=\"DGRAM\",17"
	if port != 0 {
		cmd = fmt.Sprintf("AT+NSOCR=\"DGRAM\",17,%d,1", port)
	}

	socket := 0
	err := d.cmd.Transact(cmd, func(s string) error {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil
		}
		socket = n
		return nil
	})

	return socket, err
}

func (d *n211) SendUDP(socket int, address net.IP, remotePort int, data []byte) (int, error) {
	socketReturn := 0
	lengthReturn := 0

	cmd := fmt.Sprintf("AT+NSOST=%d,\"%s\",%d,%d,\"%x\"", socket, address.String(), remotePort, len(data), data)
	err := d.cmd.Transact(cmd, func(s string) error {
		parts := strings.Split(s, ",")
		var err error
		if len(parts) == 2 {
			socketReturn, err = strconv.Atoi(parts[0])
			if err != nil {
				return err
			}

			lengthReturn, err = strconv.Atoi(parts[1])
			if err != nil {
				return err
			}

			// These should be identical or there is something wrong
			if socket != socketReturn {
				log.Printf("Inconsistency: socket in response did not match socket in request")
			}
		}
		return nil
	})

	return lengthReturn, err
}

func (d *n211) ReceiveUDP(socket int, length int) (*at.ReceivedData, error) {
	var data at.ReceivedData

	err := d.cmd.Transact(fmt.Sprintf("AT+NSORF=%d,%d", socket, length), func(s string) error {
		var err error

		parts := strings.Split(s, ",")
		if len(parts) < 6 {
			return nil
		}

		data.Socket, err = strconv.Atoi(parts[0])
		if err != nil {
			return err
		}

		data.IP = at.TrimQuotes(parts[1])

		data.Port, err = strconv.Atoi(parts[2])
		if err != nil {
			return err
		}

		data.Length, err = strconv.Atoi(parts[3])
		if err != nil {
			return err
		}

		// the data is in hex so we have to decode it first
		data.Data, err = hex.DecodeString(parts[4])
		if err != nil {
			return err
		}

		return nil
	})

	return &data, err
}

func (d *n211) CloseUDPSocket(socket int) error {
	return d.cmd.Transact(fmt.Sprintf("AT+NSOCL=%d", socket), nil)
}
