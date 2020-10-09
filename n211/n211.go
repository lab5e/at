// Package n211 implements interface to uBlox SARA N211 module
//
package n211

import (
	"bufio"
	"context"
	"errors"
	"log"
	"time"

	"github.com/lab5e/at"
	"github.com/tarm/serial"
)

// N211 maintains the state for connection to Sara N211
type n211 struct {
	device      string
	baudRate    int
	port        *serial.Port
	inputChan   chan string
	outputChan  chan string
	lineTimeout time.Duration
	debug       bool
	ctx         context.Context
	cancel      context.CancelFunc
}

const (
	defaultLineTimeout = 5 * time.Second
)

var (
	// ErrReadTimeout read from device timed out
	ErrReadTimeout = errors.New("Read timed out")

	// ErrATError ...
	ErrATError = errors.New("Device returned ERROR")
)

// New creates a new instance of the N211 interface
func New(device string, baudRate int) at.Device {
	cctx, cancel := context.WithCancel(context.Background())
	return &n211{
		device:      device,
		baudRate:    baudRate,
		inputChan:   make(chan string, 10),
		outputChan:  make(chan string, 10),
		lineTimeout: defaultLineTimeout,
		debug:       false,
		ctx:         cctx,
		cancel:      cancel,
	}
}

func (d *n211) Start() error {
	p, err := serial.OpenPort(&serial.Config{
		Name: d.device,
		Baud: d.baudRate,
	})
	if err != nil {
		return err
	}

	d.port = p

	go d.outputReader(d.ctx)
	go d.inputReader(d.ctx)

	return nil
}

func (d *n211) Close() {
	if d.port != nil {
		d.port.Close()
	}
	d.cancel()
}

// inputReader reads from the input channel and sends the strings as
// byte arrays to the serial port.
func (d *n211) inputReader(ctx context.Context) {
	for {
		select {
		case line := <-d.inputChan:
			_, err := d.port.Write([]byte(line))
			if err != nil {
				log.Fatalf("Error writing to %s: %v", d.device, err)
			}
		case <-ctx.Done():
			log.Printf("Terminating inputReader")
			return
		}
	}
}

// outputReader reads output from the device and prints it out
func (d *n211) outputReader(ctx context.Context) {
	scanner := bufio.NewScanner(d.port)
	for scanner.Scan() {
		select {
		case d.outputChan <- scanner.Text():
		case <-ctx.Done():
			log.Printf("Terminating outputReader")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// drainOutput drains the output channel
func (d *n211) drainOutput() {
	for len(d.outputChan) > 0 {
		select {
		case s := <-d.outputChan:
			d.consumeOutput(s)
		default:
			return
		}
	}
}

// consumeOutput is used to consume response codes, both solicited and
// unsolicited, so that we can use this in a state machine later.
// This function is meant to be used wherever you process data from
// the device so that you can focus on the data you want and not have
// to track stuff that you aren't really expecting or are interested
// in.  (Like when you get URC messages).
func (d *n211) consumeOutput(s string) {
	if d.debug {
		log.Printf("CONSUME '%s'", s)
	}
}

// transact drains the output from the device, then sends the
// string appending CRLF to the device, then it reads the response and
// calls fn with each line in the response as long as fn returns
// ErrContinueReading.
func (d *n211) transact(s string, fn func(string) error) error {
	var debugLog []string

	d.drainOutput()
	d.SendCRLF(s)

	// Append the outgoing command to log
	debugLog = append(debugLog, " > "+s)

	// If we didn't get a callback function we define a default
	// function.  This makes the logic a bit more regular.
	if fn == nil {
		fn = func(s string) error { return nil }
	}

	// Loop over the response until we have ERROR or OK
	var line string
	for {
		select {
		case line = <-d.outputChan:
		case <-time.After(d.lineTimeout):
			return ErrReadTimeout
		}

		debugLog = append(debugLog, " < "+line)

		// Handle OK response which should always be the last line in
		// any successful command (except unsolicited messages)
		if line == "OK" {
			if d.debug {
				for n, lin := range debugLog {
					log.Printf("[%2d] %s", n, lin)
				}
			}
			return nil
		}

		// Handle error response
		if line == "ERROR" {
			// When we have an error we're going to log it regardless
			// of whether debug is on or off.  You kind of want to
			// know about errors.
			for n, lin := range debugLog {
				log.Printf("[%2d] %s", n, lin)
			}
			return ErrATError
		}

		d.consumeOutput(line)

		err := fn(line)
		if err != nil {
			return err
		}
	}
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
	}
	return s
}
