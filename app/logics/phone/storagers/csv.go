package storagers

import (
	"encoding/csv"
	"errors"
	"github.com/totoval/framework/helpers/zone"
	"os"
)

type csvStorage struct {
	fd  *os.File
	err error
}

func NewCSV(path string) *csvStorage {
	cs := &csvStorage{}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		cs.fd, cs.err = os.Create(path)
	} else {
		cs.fd, cs.err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	}

	return cs
}
func (cs *csvStorage) Add(data map[string]string) error {
	// @todo type switch

	sender, ok := data["sender"]
	if !ok {
		return errors.New("sender not exist")
	}
	content, ok := data["content"]
	if !ok {
		return errors.New("content not exist")
	}
	raw, ok := data["raw"]
	if !ok {
		return errors.New("raw not exist")
	}

	w := csv.NewWriter(cs.fd)
	if err := w.Write([]string{zone.Now().String(), sender, content, raw}); err != nil {
		return err
	}
	w.Flush()
	return w.Error()

}
func (cs *csvStorage) Close() error {
	if cs.fd == nil {
		return nil
	}
	return cs.fd.Close()
}
