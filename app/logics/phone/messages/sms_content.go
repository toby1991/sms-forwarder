package messages

import (
	"errors"
	"regexp"
)

const REGEXP_CMGR = `\+CMGR\:\ \d+\,\"\"\,\d+\r\n(.*)`
func ParseSmsContent(msg []byte) (matched bool, pduStr string, err error){
	msgStr :=  string(msg[:])
	matched, err = regexp.MatchString(REGEXP_CMGR, msgStr)
	if err != nil {
		return matched,  "", err
	}

	if !matched {
		return matched,"",errors.New("not a message content")
	}

	pduStr, err = CMGR(msgStr)
	return matched, pduStr, err
}


func CMGR(msg string) (pduStr string, err error) {
	re, err := regexp.Compile(REGEXP_CMGR)
	if err != nil {
		return  "", err
	}

	pduStr = re.FindStringSubmatch(msg)[1]

	return pduStr, nil
}
