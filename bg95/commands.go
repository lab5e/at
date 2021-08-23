package bg95

import (
	"errors"
	"log"
	"strings"
)

func (d *bg95) GetCCID() (string, error) {
	iccid := ""
	err := d.cmd.Transact("AT+CCID", func(s string) error {
		parts := strings.Split(s, ":")
		if len(parts) != 2 {
			log.Printf("Unable to parse response AT+ICCID: %s", s)
			return errors.New("unable to parse response")
		}
		iccid = strings.TrimSpace(parts[1])
		return nil
	})
	return iccid, err
}
