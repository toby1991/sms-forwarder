package utils

import (
	"errors"
	"strings"
)

func ReverseStr(str string) (string,error) {
	if len(str) != 2{
		return "", errors.New("string length is not 2")
	}
	newStr := strings.Repeat(str, 2)
	return newStr[1:3], nil
}

func ReverseSentence(sentence string) (string, error){
	if len(sentence) %2 != 0{
		return "", errors.New("sentence length must be even")
	}

	reversedSentence := ""
	for i:=0; i<len(sentence); i+=2{
		str, err := ReverseStr(sentence[i:i+2])
		if err != nil{
			return "", err
		}
		reversedSentence += str
	}

	return reversedSentence, nil
}