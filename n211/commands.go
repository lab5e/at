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
func (d *n211) AT() error {
	return d.transact("AT", func(s string) error {
		return nil
	})
}

func (d *n211) SetDebug(debug bool) {
	d.debug = debug
}

func (d *n211) SendCRLF(s string) {
	d.inputChan <- (s + "\r\n")
}

func (d *n211) GetIMSI() (string, error) {
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

func (d *n211) GetIMEI() (string, error) {
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

func (d *n211) GetCCID() (string, error) {
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

func (d *n211) SetAutoconnect(autoconnect bool) error {
	if autoconnect {
		return d.transact("AT+NCONFIG=\"AUTOCONNECT\",\"TRUE\"", nil)
	}
	return d.transact("AT+NCONFIG=\"AUTOCONNECT\",\"FALSE\"", nil)
}

func (d *n211) Reboot() error {
	return d.transact("AT+NRB", nil)
}

func (d *n211) SetAPN(apn string) error {
	err := d.SetAutoconnect(false)
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

	err = d.SetAutoconnect(true)
	if err != nil {
		return err
	}

	err = d.Reboot()
	if err != nil {
		return err
	}

	return nil
}

func (d *n211) GetAPN() (*at.APN, error) {
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

func (d *n211) GetAddr() (int, string, error) {
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

func (d *n211) SetRadio(on bool) error {
	ind := 0
	if on {
		ind = 1
	}
	return d.transact(fmt.Sprintf("AT+CFUN=%d", ind), nil)
}

func (d *n211) GetStats() (*at.Stats, error) {
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
