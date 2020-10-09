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
	Start() error
	Close()
	SetDebug(debug bool)
	AT() error
	Reboot() error
	SendCRLF(s string)
	GetIMSI() (string, error)
	GetIMEI() (string, error)
	GetCCID() (string, error)
	SetAutoconnectOff() error
	SetAutoconnectOn() error
	SetAPN(apn string) error
	GetAPN() (*APN, error)
	GetAddr() (int, string, error)
	SetRadio(bool) error
	GetStats() (*Stats, error)
	CreateUDPSocket(port int) (int, error)
	SendUDP(socket int, address net.IP, remotePort int, data []byte) (int, error)
	ReceiveUDP(socket int, length int) (*ReceivedData, error)
	CloseUDPSocket(socket int) error
}
