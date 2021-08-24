package at

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// DefeaultImplementation is a default implementation
type DefaultImplementation struct {
	Cmd *CommandInterface
}

func (d *DefaultImplementation) Start() error {
	return d.Cmd.Start()
}

func (d *DefaultImplementation) Close() {
	d.Cmd.Close()
}

func (d *DefaultImplementation) AT() error {
	return d.Cmd.Transact("AT", func(s string) error {
		return nil
	})
}

func (d *DefaultImplementation) SetDebug(debug bool) {
	d.Cmd.SetDebug(debug)
}

func (d *DefaultImplementation) GetIMSI() (string, error) {
	var imsi string

	err := d.Cmd.Transact("AT+CIMI", func(s string) error {
		sub := IMSIRegex.FindStringSubmatch(s)
		if len(sub) > 0 {
			imsi = sub[1]
		}
		return nil
	})

	return imsi, err
}

func (d *DefaultImplementation) GetIMEI() (string, error) {
	var imei string

	err := d.Cmd.Transact("AT+CGSN", func(s string) error {
		if strings.TrimSpace(s) == "" {
			return nil
		}
		imei = s
		return nil
	})

	return imei, err
}

func (d *DefaultImplementation) SetRadio(on bool) error {
	ind := 0
	if on {
		ind = 1
	}
	return d.Cmd.Transact(fmt.Sprintf("AT+CFUN=%d", ind), nil)
}

func (d *DefaultImplementation) GetAddr() (int, string, error) {
	var cid int
	var addr string

	err := d.Cmd.Transact("AT+CGPADDR", func(s string) error {
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

			addr = TrimQuotes(parts[1])
		}
		return nil
	})

	return cid, addr, err
}

func (d *DefaultImplementation) SetAPN(apn string) error {

	err := d.Cmd.Transact(fmt.Sprintf("AT+CGDCONT=1,\"IP\",\"%s\"", apn), nil)
	if err != nil {
		return err
	}

	err = d.Cmd.Transact("AT+CGACT=1,1", func(s string) error {
		return nil
	})
	return err
}

func (d *DefaultImplementation) GetAPN() (*APN, error) {
	var apn = &APN{}

	err := d.Cmd.Transact("AT+CGDCONT?", func(s string) error {
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
			apn.PDPType = TrimQuotes(parts[1])
			apn.Name = TrimQuotes(parts[2])
			apn.Address = TrimQuotes(parts[3])
			return nil
		}

		return nil
	})

	return apn, err
}
