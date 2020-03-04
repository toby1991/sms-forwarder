package phone

import (
	"bytes"
	"github.com/tarm/serial"
	"github.com/totoval/framework/helpers/log"
	"github.com/totoval/framework/helpers/toto"
	"github.com/totoval/framework/helpers/zone"
	"io"
)

type sim800c struct {
	writer *serial.Port
	conf *serial.Config
	chB chan []byte
	chErr chan error
	isReading bool
}
func NewSim800c(comPort string, baudRate int, readTimeout zone.Duration) (*sim800c, error) {
	s := &sim800c{
		chB:make(chan []byte, 50*3), // "+CPMS: \"SM_P\",50,50,\"SM_P\",50,50,\"SM_P\",50,50"
		chErr:make(chan error, 50*3),
	}
	s.conf = &serial.Config{Name: comPort, Baud: baudRate, ReadTimeout:readTimeout}

	var err error
	if s.writer, err = serial.OpenPort(s.conf); err != nil{
		return nil, err
	}

	if err := s.init(); err != nil{
		return nil, err
	}

	return s, nil
}
func (s *sim800c) init() error {
	// echo off
	if err := s.Write([]byte("ATE0\r\n")); err != nil{
		return err
	}
	// useful error messages
	if err := s.Write([]byte("AT+CMEE=1\r\n")); err != nil{
		return err
	}
	// disable notifications
	// if  err := s.Write([]byte("AT+WIND=0\r\n")); err != nil{
	// 	return err
	// }
	// switch to TEXT:1, PDU:0 mode
	if  err := s.Write([]byte("AT+CMGF=0\r\n")); err != nil{
		return err
	}
	return nil
}
func (s *sim800c) Close() error {
	defer close(s.chB)
	defer close(s.chErr)
	return s.writer.Close()
}
func (s *sim800c) flush() error {
	return s.writer.Flush()
}
func (s *sim800c) read(b []byte) (int, error) {
	return s.writer.Read(b)
}
func (s *sim800c) write(b []byte) (int, error) {
	return s.writer.Write(b)
}
func (s *sim800c) Read2() {
	if s.isReading {
		return
	}

	s.isReading = true // only read once
	defer func(){
		s.isReading = false
	}()

	var b []byte
	for {
		_b := make([]byte, 128)

		_n, err := s.read(_b)
		if err != nil {
			if err == io.EOF {
				if len(b) <= 0 {
					// no message, continue receiving
					continue
				}

				if _n <= 0 {
					// received finished
					return
				}

				// len(b) > 0 && _n > 0 received aborted
				s.chErr <- io.EOF
				s.chB <- b
				return
			}
			// received error
			s.chErr <- err
			return
		}

		if _n > 0{
			b = append(b, _b[:_n]...)
		}

		// explode message use \r\n
		if bytes.HasPrefix(b, []byte("\r\n")) &&bytes.HasSuffix(b, []byte("\r\n")) {
			// \r\nxxxxx\r\n
			if bytes.Contains(b, []byte("\r\n\r\n")){
				// concated msg bytes
				__b := bytes.Trim(b, "\r\n")
				msgArr := bytes.Split(__b, []byte("\r\n\r\n")) // [][data]
				log.Warn(len(msgArr))
				for _, msg := range msgArr {
					log.Warn(msg)
					s.chB <- msg
				}
			}else{
				// single msg bytes
				__b := bytes.Trim(b, "\r\n") // data
				s.chB <- __b
			}

			break // msg end
		}
	}
	return
}

func (s *sim800c) Read() (n int, b []byte, err error) {
	for {
		_b := make([]byte, 128)

		_n, err := s.read(_b)
		n += _n
		if _n > 0{
			b = append(b, _b[:_n]...)
		}

		if bytes.HasSuffix(b, []byte("\r\n")){
			break
		}

		if err != nil {
			if err == io.EOF {
				if n > 0{
					break
				}
				continue
			}
			return n, nil, err
		}
	}
	return n, b[:n], nil
}


func (s *sim800c)Write( b []byte) error {
	if err := s.flush(); err != nil{
		return err
	}
	n, err := s.write(b)
	if err != nil{
		return err
	}
	log.Info("Send Bytes", toto.V{"bytes": string(b[:]), "length": n})
	return nil
}

func (s *sim800c)Error() <-chan error {
	return s.chErr
}
func (s *sim800c)Bytes() <-chan []byte {
	return s.chB
}


