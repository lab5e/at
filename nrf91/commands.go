package nrf91

import (
	"errors"
	"fmt"
	"log"
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
	var imei string

	err := d.cmd.Transact("AT+CGSN", func(s string) error {
		if strings.TrimSpace(s) == "" {
			return nil
		}
		imei = s
		return nil
	})

	return imei, err
}

func (d *nrf91) GetCCID() (string, error) {
	iccid := ""
	err := d.cmd.Transact("AT%XICCID", func(s string) error {
		parts := strings.Split(s, ":")
		if len(parts) != 2 {
			log.Printf("Unable to parse response AT%%XICCID: %s", s)
			return errors.New("unable to parse response")
		}
		iccid = strings.TrimSpace(parts[1])
		return nil
	})
	return iccid, err
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

func (d *nrf91) SetAPN(apn string) error {

	err := d.cmd.Transact(fmt.Sprintf("AT+CGDCONT=1,\"IP\",\"%s\"", apn), nil)
	if err != nil {
		return err
	}

	err = d.cmd.Transact("AT+CGACT=1,1", func(s string) error {
		return nil
	})
	return err
}

func (d *nrf91) GetAPN() (*at.APN, error) {
	var apn = &at.APN{}

	err := d.cmd.Transact("AT+CGDCONT?", func(s string) error {
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
			apn.PDPType = at.TrimQuotes(parts[1])
			apn.Name = at.TrimQuotes(parts[2])
			apn.Address = at.TrimQuotes(parts[3])
			return nil
		}

		return nil
	})

	return apn, err
}

func (d *nrf91) Reboot() error {
	return errors.New("not supported")
}

func (d *nrf91) SetAutoconnect(bool) error {
	return errors.New("not supported")
}
