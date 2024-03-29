package nrf91

import (
	"errors"
	"log"
	"strings"
)

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
