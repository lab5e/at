package nrf91

import "github.com/lab5e/at"

const DefaultBaudRate = 115200

type nrf91 struct {
	at.DefaultImplementation

	cmd *at.CommandInterface
}

func New(serialDevice string, baudRate int) at.Device {
	cmdIF := at.NewCommandInterface(serialDevice, baudRate)
	return &nrf91{
		DefaultImplementation: at.DefaultImplementation{Cmd: cmdIF},
		cmd:                   cmdIF,
	}
}
