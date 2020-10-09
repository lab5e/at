package n211

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

// CreateUDPSocket creates an UDP socket. If the port is non-zero,
// receiving is enabled and +NSONMI URCs will appear for any message
// that is received on that port.
func (d *N211) CreateUDPSocket(port int) (int, error) {
	// Section 16.24.3 states that this port is reserved
	if port == 5683 {
		return -1, errors.New("Reserved port value")
	}

	cmd := "AT+NSOCR=\"DGRAM\",17"
	if port != 0 {
		cmd = fmt.Sprintf("AT+NSOCR=\"DGRAM\",17,%d,1", port)
	}

	socket := 0
	err := d.transact(cmd, func(s string) error {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil
		}
		socket = n
		return nil
	})

	return socket, err
}

// SendUDP sends an UDP packet using a previously opened socket.
// Returns number of bytes written and whether or not there was an
// error.
func (d *N211) SendUDP(socket int, address net.IP, remotePort int, data []byte) (int, error) {
	socketReturn := 0
	lengthReturn := 0

	cmd := fmt.Sprintf("AT+NSOST=%d,\"%s\",%d,%d,\"%x\"", socket, address.String(), remotePort, len(data), data)
	err := d.transact(cmd, func(s string) error {
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

// ReceiveUDP data on a socket. When data arrives a +NSONMI URC will
// be issued indicating the socket the message was received on and the
// amount of data.
//
// This command takes a length, which is the maximum amount of data
// that will be returned. If the requested length is larger than the
// actual size of the returned data, only the length of returned data
// is provided, and the remaining length is returned as 0.
//
// If the requested length is less than the amount of data returned,
// only the requested amount of data will be returned, plus an
// indication of the number of bytes remaining. Once a message has
// been fully read, a new +NSONMI URC will be sent if there is another
// message to process.
func (d *N211) ReceiveUDP(socket int, length int) (*at.ReceivedData, error) {
	var data at.ReceivedData

	err := d.transact(fmt.Sprintf("AT+NSORF=%d,%d", socket, length), func(s string) error {
		var err error

		parts := strings.Split(s, ",")
		if len(parts) < 6 {
			return nil
		}

		data.Socket, err = strconv.Atoi(parts[0])
		if err != nil {
			return err
		}

		data.IP = trimQuotes(parts[1])

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

// CloseUDPSocket - Close the specified socket. The pending messages
// to be read (if present) will be dropped. No further +NSONMI URCs
// will be generated. If the socket has already been closed, or was
// never created, an error result code will be issued.
func (d *N211) CloseUDPSocket(socket int) error {
	return d.transact(fmt.Sprintf("AT+NSOCL=%d", socket), nil)
}
