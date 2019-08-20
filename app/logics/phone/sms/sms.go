package sms

import (
	"fmt"
	"totoval/app/logics/phone/interfaces"
)


func Retrieve(chip interfaces.Chipper, smsIndex uint) error {
	if err := chip.Write([]byte(fmt.Sprintf("AT+CMGR=%d\r\n", smsIndex))); err != nil{
		return err
	}
	//
	// //n, b, err := s.chip.Read()
	// go func() {
	// 	for{
	// 		chip.Read2() // useless
	// 	}
	// }()
	//
	// for {
	// 	select {
	// 	case <-chip.Bytes():
	// 		return nil
	// 	case err := <-chip.Error():
	// 		return err
	// 	default:
	// 		return nil
	// 	}
	// }
	return nil
}
func Delete(chip interfaces.Chipper, smsIndex uint) error {
	if err := chip.Write([]byte(fmt.Sprintf("AT+CMGD=%d\r\n", smsIndex))); err != nil{
		return err
	}
	return nil
}

