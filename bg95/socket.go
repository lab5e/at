package bg95

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/lab5e/at"
)

/*OK
[14:00:37][AT] Context activated OK
[14:00:37][AT] AT+QIOPEN=1,1,"UDP","172.16.15.14", 1234, 3030,0
OK

+QIOPEN: 1,0
[14:00:37][AT] Socket opened OK
[14:00:37][AT] Sending UDP message...
[14:00:37][AT] AT+QISEND=1,22
>
[14:00:37][INFO] Sending 22 bytes to 172.16.15.14 on port 1234
[14:00:37][AT] UDP message sent ok
[14:00:38][AT] AT+QICLOSE=1
OK
*/

var socketno = 0

// Note: There is only a single socket in the LTE modem. The port parameter may be 0, then the
// socket won't be bound to a port.
func (d *bg95) CreateUDPSocket(port int) (int, error) {
	socketno++
	if socketno > 11 {
		return 0, errors.New("sockets exhausted")
	}
	err := d.cmd.Transact(fmt.Sprintf(`AT+QIOPEN=1,%d,"UDP SERVICE","0.0.0.0",%d,1`, socketno, port), func(s string) error {
		// The module will return <connection id>,<error> and <error> should - obviously be 0
		if strings.HasPrefix(s, "+QIOPEN") {
			elems := strings.Split(s, ":")
			if len(elems) != 2 {
				log.Printf("Error parsing result from AT+QIOPEN: %s", s)
				return errors.New("could not parse return value from AT+QIOPEN command")
			}
			fields := strings.Split(elems[1], ",")
			if len(fields) != 2 {
				log.Printf("Expected 2 fields returned but got %d: %s", len(fields), s)
				return errors.New("could not parse response fields from AT+QIOPEN")
			}
			connID, err := strconv.Atoi(strings.TrimSpace(fields[0]))
			if err != nil {
				log.Printf("Invalid connection ID (field #1) from module: %s", s)
				return errors.New("could not parse connection ID")
			}
			if strings.TrimSpace(fields[1]) != "0" {
				log.Printf("Got error response from AT+QIOPEN: %s", s)
				return errors.New("error code returned from module")
			}
			socketno = connID
			return nil
		}
		return nil
	})
	return socketno, err
}

var errDummy = errors.New("not an error")

// This might not work with arbitrary binary data since the data is sent as a string
func (d *bg95) SendUDP(socket int, address net.IP, remotePort int, data []byte) (int, error) {
	// Start by sending AT+SEND and when we receive the '>' character send the payload. We should get a "SEND OK" or "SEND ERROR" back
	err := d.cmd.Transact(fmt.Sprintf(`AT+QISEND=%d,%d,"%s",%d\r\n%s`,
		socket, len(data), address.String(), remotePort, string(data)), func(s string) error {
		if s == "SEND ERROR" {
			return errors.New("send error")
		}
		if s == "SEND OK" {
			return errDummy
		}
		return nil
	})
	if err == errDummy {
		return len(data), nil
	}
	return 0, err
}

// Note: Receive has a 10 second timeout
func (d *bg95) ReceiveUDP(socket int, length int) (*at.ReceivedData, error) {
	return nil, errors.New("not implemented")
}

func (d *bg95) CloseUDPSocket(socket int) error {
	return d.cmd.Transact(fmt.Sprintf("AT+QICLOSE=%d", socket), func(s string) error {
		return nil
	})
}
