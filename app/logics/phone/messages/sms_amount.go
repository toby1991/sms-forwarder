package messages

import (
	"errors"
	"regexp"
	"strconv"
)

// +CPMS: "SM_P",5,50,"SM_P",5,50,"SM_P",5,50
const REGEXP_CPMS = `\+CPMS\:\ \"SM_P\"\,(\d+)\,(\d+)\,`
func ParseSmsAmount(msg []byte) (matched bool, amount uint, max uint, err error){
	msgStr :=  string(msg[:])
	matched, err = regexp.MatchString(REGEXP_CPMS, msgStr)
	if err != nil {
		return matched,  0, 0, err
	}

	if !matched {
		return matched,0,0,errors.New("not a message amount")
	}

	amount, max, err = CPMS(msgStr)
	return matched, amount, max, err
}


func CPMS(msg string) (amount uint, max uint,err error) {
	re, err := regexp.Compile(REGEXP_CPMS)
	if err != nil {
		return  0, 0, err
	}

	amount32, err := strconv.ParseInt(re.FindStringSubmatch(msg)[1], 10, 32)
	if err != nil{
		return 0, 0, err
	}
	max32, err := strconv.ParseInt(re.FindStringSubmatch(msg)[2], 10, 32)
	if err != nil{
		return 0, 0, err
	}

	return uint(amount32), uint(max32), nil
}
