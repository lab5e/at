package at

import (
	"bufio"
	"context"
	"errors"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/tarm/serial"
)

var (
	// IMSIRegex is a regexp that matches an IMSI number
	IMSIRegex = regexp.MustCompile("([0-9]{5,15})")

	// IMEIRegex is a regexp that matches an IMEI number
	IMEIRegex = regexp.MustCompile(`\+CGSN: ([0-9]{5,15})`)

	// CCIDRegex is a regexp that matches an CCID
	CCIDRegex = regexp.MustCompile(`\+CCID: ([0-9]{5,15})`)
)

var (
	// ErrReadTimeout read from device timed out
	ErrReadTimeout = errors.New("read timed out")

	// ErrATError ...
	ErrATError = errors.New("device returned ERROR")
)

const DefaultLineTimeout = 5 * time.Second

// CommandInterface is a helper type for modems. It's optional to use it when implementing
// support for new modules but quite helpful.
type CommandInterface struct {
	device      string
	baudRate    int
	port        *serial.Port
	inputChan   chan string
	outputChan  chan string
	lineTimeout time.Duration
	debug       bool
	ctx         context.Context
	cancel      context.CancelFunc
	errors      []string
	successes   []string
	splits      []string
}

func NewCommandInterface(device string, baudRate int) *CommandInterface {
	ctx, cancel := context.WithCancel(context.Background())
	return &CommandInterface{
		device:      device,
		baudRate:    baudRate,
		inputChan:   make(chan string, 10),
		outputChan:  make(chan string, 10),
		lineTimeout: DefaultLineTimeout,
		debug:       false,
		ctx:         ctx,
		cancel:      cancel,
		errors:      []string{"ERROR"},
		successes:   []string{"OK"},
		splits:      []string{"\r\n"},
	}
}

// Some modems (hello BG95) sends additional text string when a command completes successful.
func (c *CommandInterface) AddSuccessOutput(newSuccess string) {
	c.successes = append(c.successes, newSuccess)
}

// Some modems (hello BG95) sends additional text strings when a command completes with an error
func (c *CommandInterface) AddErrorOutput(newError string) {
	c.errors = append(c.errors, newError)
}

// Add a command split character or set of characters. Since some modems (hello BG95) changes
// its behaviour in certain comamnds the CRLF sequence isn't always emitted when the modem is
// ready to accept new characters.
func (c *CommandInterface) AddSplitChars(newSplit string) {
	c.splits = append(c.splits, newSplit)
}

func (c *CommandInterface) SetDebug(debug bool) {
	c.debug = debug
}

func (c *CommandInterface) Start() error {
	p, err := serial.OpenPort(&serial.Config{
		Name: c.device,
		Baud: c.baudRate,
	})
	if err != nil {
		return err
	}

	c.port = p

	go c.outputReader(c.ctx)
	go c.inputReader(c.ctx)

	return nil
}

func (c *CommandInterface) Close() {
	if c.port != nil {
		c.port.Close()
	}
	c.cancel()
}

// inputReader reads from the input channel and sends the strings as
// byte arrays to the serial port.
func (c *CommandInterface) inputReader(ctx context.Context) {
	for {
		select {
		case line := <-c.inputChan:
			_, err := c.port.Write([]byte(line))
			if err != nil {
				log.Fatalf("Error writing to %s: %v", c.device, err)
			}
		case <-ctx.Done():
			log.Printf("Terminating inputReader")
			return
		}
	}
}

func (c *CommandInterface) splitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for _, v := range c.splits {
		pos := strings.Index(string(data), v)
		if pos >= 0 {
			advance = pos + len(v)
			token = data[:pos]
			err = nil
			return
		}
	}
	return 0, nil, nil
}

// outputReader reads output from the device and prints it out
func (c *CommandInterface) outputReader(ctx context.Context) {
	scanner := bufio.NewScanner(c.port)
	scanner.Split(c.splitFunc)
	for scanner.Scan() {
		select {
		case c.outputChan <- scanner.Text():
		case <-ctx.Done():
			log.Printf("Terminating outputReader")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// drainOutput drains the output channel
func (c *CommandInterface) drainOutput() {
	for len(c.outputChan) > 0 {
		select {
		case s := <-c.outputChan:
			c.consumeOutput(s)
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
func (c *CommandInterface) consumeOutput(s string) {
	if c.debug {
		log.Printf("CONSUME '%s'", s)
	}
}

// transact drains the output from the device, then sends the
// string appending CRLF to the device, then it reads the response and
// calls fn with each line in the response as long as fn returns
// ErrContinueReading.
func (c *CommandInterface) Transact(s string, fn func(string) error) error {
	var debugLog []string

	c.drainOutput()
	c.SendCRLF(s)

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
		case line = <-c.outputChan:
		case <-time.After(c.lineTimeout):
			return ErrReadTimeout
		}

		debugLog = append(debugLog, " < "+line)

		// Handle OK response which should always be the last line in
		// any successful command (except unsolicited messages)
		for _, v := range c.successes {
			if line == v {
				if c.debug {
					for n, lin := range debugLog {
						log.Printf("[%2d] %s", n, lin)
					}
				}
				return nil
			}
		}

		// Handle error response
		for _, v := range c.errors {
			if line == v {
				// When we have an error we're going to log it regardless
				// of whether debug is on or off.  You kind of want to
				// know about errors.
				for n, lin := range debugLog {
					log.Printf("[%2d] %s", n, lin)
				}
				return ErrATError

			}
		}

		c.consumeOutput(line)

		err := fn(line)
		if err != nil {
			return err
		}
	}
}

func (c *CommandInterface) SendCRLF(s string) {
	c.inputChan <- (s + "\r\n")
}

func TrimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
	}
	return s
}
