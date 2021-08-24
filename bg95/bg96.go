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
	cmdIF.AddErrorOutput("SEND FAIL")
	cmdIF.AddSplitChars(">")
	cmdIF.AddSuccessOutput("SEND OK")
	return &bg95{
		DefaultImplementation: at.DefaultImplementation{Cmd: cmdIF},
		cmd:                   cmdIF,
	}
}

func (d *bg95) Start() error {
	if err := d.Cmd.Start(); err != nil {
		return err
	}
	// BG95 has echo turned on by default. Turn off
	return d.Cmd.Transact("ATE0", func(s string) error {
		return nil
	})
}
