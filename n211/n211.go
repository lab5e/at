// Package n211 implements interface to uBlox SARA N211 module
//
package n211

import (
	"github.com/lab5e/at"
)

// N211 maintains the state for connection to Sara N211
type n211 struct {
	cmd *at.CommandInterface
}

// New creates a new instance of the N211 interface
func New(device string, baudRate int) at.Device {
	return &n211{
		cmd: at.NewCommandInterface(device, baudRate),
	}
}

func (d *n211) Start() error {
	return d.cmd.Start()
}

func (d *n211) Close() {
	d.cmd.Close()
}
