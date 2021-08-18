package nrf91

import "github.com/lab5e/at"

const DefaultBaudRate = 115200

type nrf91 struct {
	cmd *at.CommandInterface
}

func New(serialDevice string, baudRate int) at.Device {
	return &nrf91{
		cmd: at.NewCommandInterface(device, baudRate),
	}
}

func (d *nrf91) Start() error {
	return d.cmd.Start()
}

func (d *nrf91) Close() {
	d.cmd.Close()
}
