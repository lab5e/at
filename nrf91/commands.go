package nrf91

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/lab5e/at"
)

func (d *nrf91) AT() error {
	return d.cmd.Transact("AT", func(s string) error {
		return nil
	})
}

func (d *nrf91) SetDebug(debug bool) {
	d.cmd.SetDebug(debug)
}

func (d *nrf91) GetIMSI() (string, error) {
	var imsi string

	err := d.cmd.Transact("AT+CIMI", func(s string) error {
		sub := at.IMSIRegex.FindStringSubmatch(s)
		if len(sub) > 0 {
			imsi = sub[1]
		}
		return nil
	})

	return imsi, err
}

func (d *nrf91) GetIMEI() (string, error) {
	var imsi string

	err := d.cmd.Transact("AT+CGSN", func(s string) error {
		sub := at.IMEIRegex.FindStringSubmatch(s)
		if len(sub) > 0 {
			imsi = sub[1]
		}
		return nil
	})

	return imsi, err
}

func (d *nrf91) GetCCID() (string, error) {
	// Not supported on older versions of the AT client
	return "", errors.New("not supported")
}

func (d *nrf91) GetAddr() (int, string, error) {
	var cid int
	var addr string

	err := d.cmd.Transact("AT+CGPADDR", func(s string) error {
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

			addr = at.TrimQuotes(parts[1])
		}
		return nil
	})

	return cid, addr, err
}

func (d *nrf91) SetRadio(on bool) error {
	ind := 0
	if on {
		ind = 1
	}
	return d.cmd.Transact(fmt.Sprintf("AT+CFUN=%d", ind), nil)
}

func (d *nrf91) GetStats() (*at.Stats, error) {
	return nil, errors.New("not supported")
}
