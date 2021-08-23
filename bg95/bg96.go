package bg95

import "github.com/lab5e/at"

// DefaultBaudRate is the default baud rate for the BG95 UART
const DefaultBaudRate = 115200

type bg95 struct {
	at.DefaultImplementation

	cmd *at.CommandInterface
}

func New(serialDevice string, baudRate int) at.Device {
	cmdIF := at.NewCommandInterface(serialDevice, baudRate)
	return &bg95{
		DefaultImplementation: at.DefaultImplementation{Cmd: cmdIF},
		cmd:                   cmdIF,
	}
}
