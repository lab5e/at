package at

import "net"

// APN contains data about the APN.  We have skipped the optional
// fields to simplify matters.
type APN struct {
	ContextIdentifier int
	PDPType           string
	Name              string
	Address           string
}

// Stats contains basic operational statistics
type Stats struct {
	SignalPower int
	//  total power within receive bandwidth expressed in tenth of dBm
	TotalPower int
	// TX power expressed in tenth of dBm
	TXPower int
	// elapsed TX time since last power on event expressed in milliseconds
	TXTime int
	// elapsed RX time since last power on event expressed in milliseconds
	RXTime int
	// physical ID of the cell providing service to the module
	CellID int
	// TODO(borud): document
	ECL int
	//  last SNR value expressed in tenth of dB
	SNR int
	// TODO(borud): document
	EARFCN int
	// TODO(borud): document
	PCI int
	//  last RSRQ value expressed in tenth of dB
	RSRQ int
}

// ReceivedData contains the data received from an UDP connection
type ReceivedData struct {
	Socket    int
	IP        string
	Port      int
	Length    int
	Data      []byte
	Remaining int
}

// Device is a generic interface for mobile network devices with AT
// command interfaces.
type Device interface {

	// Start opens the serial port and starts the reader goroutine
	Start() error

	// Close the serial port
	Close()

	// SetDebug turns on debugging if debug is true
	SetDebug(debug bool)

	AT() error

	// GetIMSI reads the IMSI from the device
	GetIMSI() (string, error)

	// GetIMEI reads the IMEI from the device
	GetIMEI() (string, error)

	// GetCCID returns the CCID of the SIM
	GetCCID() (string, error)

	// SetAPN sets the APN.  Be aware that this operation performs
	// multiple transactions and reboots the device.
	SetAPN(apn string) error

	// GetAPN returns the current APN settings
	GetAPN() (*APN, error)

	// GetAddr returns the context identifier (CID) and address
	// currently allocated to the device. This (usually) invokes the AT+CGPADDR command.
	GetAddr() (int, string, error)

	// SetRadio turns the radio on if on is true and off is on is false. This (usually) invokes the AT+CFUN
	// command
	SetRadio(bool) error

	// CreateUDPSocket creates an UDP socket. If the port is non-zero,
	// receiving is enabled and +NSONMI URCs will appear for any
	// message that is received on that port.
	CreateUDPSocket(port int) (int, error)

	// SendUDP sends an UDP packet using a previously opened socket.
	// Returns number of bytes written and whether or not there was an
	// error.
	SendUDP(socket int, address net.IP, remotePort int, data []byte) (int, error)

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
	ReceiveUDP(socket int, length int) (*ReceivedData, error)

	// CloseUDPSocket - Close the specified socket. The pending messages
	// to be read (if present) will be dropped. No further +NSONMI URCs
	// will be generated. If the socket has already been closed, or was
	// never created, an error result code will be issued.
	CloseUDPSocket(socket int) error
}
