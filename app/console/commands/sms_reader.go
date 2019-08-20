package commands

import (
	"errors"
	"github.com/totoval/framework/cmd"
	"github.com/totoval/framework/helpers/toto"
	"github.com/totoval/framework/helpers/zone"
	"github.com/totoval/framework/helpers/log"
	"totoval/app/logics/phone"
	"totoval/app/logics/phone/notifiers"
)

func init() {
	cmd.Add(&SmsReader{})
}

type SmsReader struct {
}

func (sr *SmsReader) Command() string {
	return "sms:read {com_port}"
}

func (sr *SmsReader) Description() string {
	return "Read sms from specified com port, and send notification"
}

var (
	comPortNotSetErr = errors.New("com port is not set")
)

func (sr *SmsReader) Handler(arg *cmd.Arg) error {
	comPort, err := arg.Get("com_port")
	if err != nil {
		return err
	}

	if comPort == nil {
		return comPortNotSetErr
	}

	// /dev/ttyUSB0
	sim800c, err := phone.NewSim800c(*comPort, 115200, 5 * zone.Second)
	if err != nil{
		return err
	}
	defer func() {
		_ = log.Error(sim800c.Close())
	}()
	log.Info("Listening", toto.V{"com_port": *comPort})

	ph := phone.New(sim800c, new(notifiers.Pushover))

	ph.Listen()

	return nil
}
