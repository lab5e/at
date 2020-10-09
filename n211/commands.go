package n211

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/lab5e/at"
)

var (
	// IMSIRegex ...
	IMSIRegex = regexp.MustCompile("([0-9]{5,15})")

	// IMEIRegex ...
	IMEIRegex = regexp.MustCompile("\\+CGSN: ([0-9]{5,15})")

	// CCIDRegex ...
	CCIDRegex = regexp.MustCompile("\\+CCID: ([0-9]{5,15})")
)

// AT ...
func (d *N211) AT() error {
	return d.transact("AT", func(s string) error {
		return nil
	})
}

// SetDebug turns on debugging if debug is set to true and turns it
// off if set to false.
func (d *N211) SetDebug(debug bool) {
	d.debug = debug
}

// SendCRLF sends a string to the device adding CRLF to the end of the line
func (d *N211) SendCRLF(s string) {
	d.inputChan <- (s + "\r\n")
}

// GetIMSI reads the IMSI from the device
func (d *N211) GetIMSI() (string, error) {
	var imsi string

	err := d.transact("AT+CIMI", func(s string) error {
		sub := IMSIRegex.FindStringSubmatch(s)
		if len(sub) > 0 {
			imsi = sub[1]
		}
		return nil
	})

	return imsi, err
}

// GetIMEI reads the IMEI from the device
func (d *N211) GetIMEI() (string, error) {
	var imsi string

	err := d.transact("AT+CGSN=1", func(s string) error {
		sub := IMSIRegex.FindStringSubmatch(s)
		if len(sub) > 0 {
			imsi = sub[1]
		}
		return nil
	})

	return imsi, err
}

// GetCCID returns the CCID of the SIM
func (d *N211) GetCCID() (string, error) {
	var ccid string

	err := d.transact("AT+CCID", func(s string) error {
		sub := IMSIRegex.FindStringSubmatch(s)
		if len(sub) > 0 {
			ccid = sub[1]
		}
		return nil
	})
	return ccid, err
}

// SetAutoconnectOff turns off autoconnect
func (d *N211) SetAutoconnectOff() error {
	return d.transact("AT+NCONFIG=\"AUTOCONNECT\",\"FALSE\"", nil)
}

// SetAutoconnectOn turns on autoconnect
func (d *N211) SetAutoconnectOn() error {
	return d.transact("AT+NCONFIG=\"AUTOCONNECT\",\"TRUE\"", nil)
}

// Reboot device
func (d *N211) Reboot() error {
	return d.transact("AT+NRB", nil)
}

// SetAPN sets the APN.  Be aware that this operation performs
// multiple transactions and reboots the device.
//
func (d *N211) SetAPN(apn string) error {
	err := d.SetAutoconnectOff()
	if err != nil {
		return err
	}

	err = d.Reboot()
	if err != nil {
		return err
	}

	err = d.transact(fmt.Sprintf("AT+CGDCONT=0,\"IP\",\"%s\"", apn), nil)
	if err != nil {
		return err
	}

	err = d.SetAutoconnectOn()
	if err != nil {
		return err
	}

	err = d.Reboot()
	if err != nil {
		return err
	}

	return nil
}

// GetAPN returns the current APN settings
func (d *N211) GetAPN() (*at.APN, error) {
	var apn = &at.APN{}

	err := d.transact("AT+CGDCONT?", func(s string) error {
		var err error
		if st := strings.TrimPrefix(s, "+CGDCONT: "); st != s {
			parts := strings.Split(st, ",")
			if len(parts) < 4 {
				return errors.New("missing some fields in response")
			}

			apn.ContextIdentifier, err = strconv.Atoi(parts[0])
			if err != nil {
				return errors.New("invalid CID")
			}
			apn.PDPType = trimQuotes(parts[1])
			apn.Name = trimQuotes(parts[2])
			apn.Address = trimQuotes(parts[3])
			return nil
		}

		return nil
	})

	return apn, err
}

// GetAddr returns the context identifier (CID) and address
// currently allocated to the device.
func (d *N211) GetAddr() (int, string, error) {
	var cid int
	var addr string

	err := d.transact("AT+CGPADDR", func(s string) error {
		var err error
		if st := strings.TrimPrefix(s, "+CGPADDR: "); st != s {
			parts := strings.Split(st, ",")
			if len(parts) < 2 {
				return errors.New("missing field in response")
			}

			cid, err = strconv.Atoi(parts[0])
			if err != nil {
				return errors.New("invalid CID")
			}

			addr = trimQuotes(parts[1])
		}
		return nil
	})

	return cid, addr, err
}

// SetRadio turns the radio on if on is true and off is on is false.
func (d *N211) SetRadio(on bool) error {
	ind := 0
	if on {
		ind = 1
	}
	return d.transact(fmt.Sprintf("AT+CFUN=%d", ind), nil)
}

// GetStats returns the most recent operational statistics.  Note that
// this translates to the AT+NUESTATS command and doesn't specify any
// parameters.
func (d *N211) GetStats() (*at.Stats, error) {
	var stats at.Stats

	err := d.transact("AT+NUESTATS", func(s string) error {
		parts := strings.Split(s, ",")
		if len(parts) < 2 {
			return nil
		}

		switch parts[0] {
		case "\"Signal power\"":
			if n, err := strconv.Atoi(parts[1]); err == nil {
				stats.SignalPower = n
			}

		case "\"Total power\"":
			if n, err := strconv.Atoi(parts[1]); err == nil {
				stats.TotalPower = n
			}

		case "\"TX power\"":
			if n, err := strconv.Atoi(parts[1]); err == nil {
				stats.TXPower = n
			}

		case "\"TX time\"":
			if n, err := strconv.Atoi(parts[1]); err == nil {
				stats.TXTime = n
			}

		case "\"RX time\"":
			if n, err := strconv.Atoi(parts[1]); err == nil {
				stats.RXTime = n
			}

		case "\"Cell ID\"":
			if n, err := strconv.Atoi(parts[1]); err == nil {
				stats.CellID = n
			}

		case "\"ECL\"":
			if n, err := strconv.Atoi(parts[1]); err == nil {
				stats.ECL = n
			}

		case "\"SNR\"":
			if n, err := strconv.Atoi(parts[1]); err == nil {
				stats.SNR = n
			}

		case "\"EARFCN\"":
			if n, err := strconv.Atoi(parts[1]); err == nil {
				stats.EARFCN = n
			}

		case "\"PCI\"":
			if n, err := strconv.Atoi(parts[1]); err == nil {
				stats.PCI = n
			}

		case "\"RSRQ\"":
			if n, err := strconv.Atoi(parts[1]); err == nil {
				stats.RSRQ = n
			}

		}
		return nil
	})
	return &stats, err
}
