package messages

import (
	"errors"
	"regexp"
	"strconv"
)

type MessageHandler func(msg string) (messageIndex uint, err error)
type MessageType struct {
	Regexp string
	Handler MessageHandler
}

const REGEXP_CMTI_SM = `\+CMTI\:\ \"SM\"\,(\d+)`
const REGEXP_CMTI_ME = `\+CMTI\:\ \"ME\"\,(\d+)`

var messageTypeList = []MessageType{
	{
		Regexp: REGEXP_CMTI_SM,
		Handler: CMTISM,
	},
	{
		Regexp:REGEXP_CMTI_ME,
		Handler: CMTIME,
	},
}

func ParseSmsIndex(msg []byte) (matched bool, messageIndex uint, err error){
	msgStr :=  string(msg[:])
	var handler MessageHandler
	for _, r := range messageTypeList {
		matched, err = regexp.MatchString(r.Regexp, msgStr)
		if err != nil {
			matched = false
		}
		if matched {
			handler = r.Handler
			break
		}
	}

	if !matched{
		return matched, 0, errors.New("not a new message")
	}
	messageIndex, err = handler(msgStr)
	return matched, messageIndex, err
}

func CMTIME(msg string) (messageIndex uint, err error) {
	re, err := regexp.Compile(REGEXP_CMTI_ME)
	if err != nil {
		return 0, err
	}
	_messageIndex, err := strconv.ParseUint(re.FindStringSubmatch(msg)[1], 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(_messageIndex), nil
}
func CMTISM(msg string) (messageIndex uint, err error) {
	re, err := regexp.Compile(REGEXP_CMTI_SM)
	if err != nil {
		return 0, err
	}
	_messageIndex, err := strconv.ParseUint(re.FindStringSubmatch(msg)[1], 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(_messageIndex), nil
}
