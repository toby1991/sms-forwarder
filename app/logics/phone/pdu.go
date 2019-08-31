package phone

import (
	"errors"
	"fmt"
	"github.com/totoval/framework/helpers/log"
	"github.com/totoval/framework/helpers/zone"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"strconv"
	"strings"
	"time"
	"totoval/app/logics/phone/utils"
)

// A：短信息中心地址长度，2位十六进制数(1字节)。
// B：短信息中心号码类型，2位十六进制数。
// C：短信息中心号码，B+C的长度将由A中的数据决定。
// D：文件头字节，2位十六进制数。
// E：信息类型，2位十六进制数。
// F：被叫号码长度，2位十六进制数。
// G：被叫号码类型，2位十六进制数，取值同B。
// H：被叫号码，长度由F中的数据决定。
// I：协议标识，2位十六进制数。
// J：数据编码方案，2位十六进制数。
// K：有效期，2位十六进制数。 // 14，时间戳
// L：用户数据长度，2位十六进制数。
// M：用户数据，其长度由L中的数据决定。J中设定采用UCS2编码，这里是中英文的Unicode字符。

type pdu struct {
	raw                                   string
	a, b, c, d, e, f, g, h, i, j, k, l, m string
	center                                string
	from                                  string
	content                               string
	length                                uint
	time                                  zone.Time
	err                                   error
}

// https://yq.aliyun.com/articles/310431
func NewPDU() *pdu {
	p := &pdu{}
	return p
}

func (p *pdu) Scan(raw string) {
	p.raw = raw

	defer func() {
		if err := recover(); err != nil {
			_ = log.Error(errors.New(fmt.Sprint(err)))
		}
	}()

	p.err = p.scan()
}
func (p *pdu) Data() (sender string, content string, err error) {
	return p.from, p.content, p.err
}

func (p *pdu) scan() error {
	p.a = p.raw[0:2]
	// fmt.Println(p.a, p.smsCenterLen())
	p.b = p.raw[2 : 2+2]
	// fmt.Println(p.b)
	afterCIndex := 2 + 2 + p.smsCenterLen() - uint(len(p.a))
	p.c = p.raw[2+2 : afterCIndex]
	cData, _ := utils.ReverseSentence(p.c)
	if p.b == "91" {
		cData = "+" + cData
	}
	if strings.HasSuffix(cData, "F") {
		cData = cData[:len(cData)-1]
	}
	// fmt.Println(p.c, cData)

	p.d = p.raw[afterCIndex : afterCIndex+2]
	// fmt.Println(p.d)

	// e is only for sender
	// 			p.e = p.raw[afterCIndex+2 :afterCIndex+2 +2]
	// 			// fmt.Println(p.e)

	p.f = p.raw[afterCIndex+2 : afterCIndex+2+2]
	// fmt.Println(p.f, p.senderLen())
	p.g = p.raw[afterCIndex+2+2 : afterCIndex+2+2+2]
	// fmt.Println(p.g)
	afterHIndex := afterCIndex + 2 + 2 + 2 + p.senderLen()
	p.h = p.raw[afterCIndex+2+2+2 : afterHIndex]
	hData, _ := utils.ReverseSentence(p.h)
	if p.g == "91" { // A1
		hData = "+" + hData
	}
	if strings.HasSuffix(hData, "F") {
		hData = hData[:len(hData)-1]
	}
	// fmt.Println(p.h, hData)

	p.i = p.raw[afterHIndex : afterHIndex+2]
	// fmt.Println(p.i)
	p.j = p.raw[afterHIndex+2 : afterHIndex+2+2]
	// fmt.Println(p.j)
	// p.k = p.raw[afterHIndex+2+2 :afterHIndex+2+2 +2]
	p.k = p.raw[afterHIndex+2+2 : afterHIndex+2+2+14]
	// fmt.Println(p.k)
	p.parseTime()

	// p.l = p.raw[afterHIndex+2+2+2 :afterHIndex+2+2+2 +2]
	p.l = p.raw[afterHIndex+2+2+14 : afterHIndex+2+2+14+2]
	// fmt.Println(p.l)

	// p.m = p.raw[afterHIndex+2+2+2+2 :]
	p.m = p.raw[afterHIndex+2+2+14+2:]

	p.center = cData
	p.from = hData

	switch p.j {

	case "00":
		//  7-bit
		if err := p.parseContent7(); err != nil {
			return err
		}
	case "04":
		// default utf8
		if err := p.parseContent8(); err != nil {
			return err
		}
	case "09":
		fallthrough
	case "19":
		fallthrough
	case "08":
		// utf16 ucs2
		if err := p.parseContent16(); err != nil {
			return err
		}
	default:
		return errors.New("encoding unknown:" + p.j)

	}
	return nil
}
func (p *pdu) create() error {
	p.a = p.raw[0:2]
	// fmt.Println(p.a, p.smsCenterLen())
	p.b = p.raw[2 : 2+2]
	// fmt.Println(p.b)
	afterCIndex := 2 + 2 + p.smsCenterLen() - uint(len(p.a))
	p.c = p.raw[2+2 : afterCIndex]
	cData, _ := utils.ReverseSentence(p.c)
	if p.b == "91" {
		cData = "+" + cData
	}
	if strings.HasSuffix(cData, "F") {
		cData = cData[:len(cData)-1]
	}
	// fmt.Println(p.c, cData)

	p.d = p.raw[afterCIndex : afterCIndex+2]
	// fmt.Println(p.d)

	// e is only for sender
	p.e = p.raw[afterCIndex+2 : afterCIndex+2+2]
	// fmt.Println(p.e)

	p.f = p.raw[afterCIndex+2+2 : afterCIndex+2+2+2]
	// fmt.Println(p.f, p.senderLen())
	p.g = p.raw[afterCIndex+2+2+2 : afterCIndex+2+2+2+2]
	// fmt.Println(p.g)
	afterHIndex := afterCIndex + 2 + 2 + 2 + 2 + p.senderLen()
	p.h = p.raw[afterCIndex+2+2+2+2 : afterHIndex]
	hData, _ := utils.ReverseSentence(p.h)
	if p.g == "91" { // A1
		hData = "+" + hData
	}
	if strings.HasSuffix(hData, "F") {
		hData = hData[:len(hData)-1]
	}
	// fmt.Println(p.h, hData)

	p.i = p.raw[afterHIndex : afterHIndex+2]
	// fmt.Println(p.i)
	p.j = p.raw[afterHIndex+2 : afterHIndex+2+2]
	// fmt.Println(p.j)
	p.k = p.raw[afterHIndex+2+2 : afterHIndex+2+2+2]
	// fmt.Println(p.k)
	p.parseTime()

	p.l = p.raw[afterHIndex+2+2+2 : afterHIndex+2+2+2+2]
	// fmt.Println(p.l)
	p.m = p.raw[afterHIndex+2+2+2+2:]

	p.center = cData
	p.from = hData
	if err := p.parseContent16(); err != nil {
		return err
	}
	return nil
}
func (p *pdu) parseTime() {
	// year
	year, err := utils.ReverseStr(p.k[0:2])
	if err != nil {
		p.err = err
		return
	}
	month, err := utils.ReverseStr(p.k[2:4])
	if err != nil {
		p.err = err
		return
	}
	day, err := utils.ReverseStr(p.k[4:6])
	if err != nil {
		p.err = err
		return
	}
	hour, err := utils.ReverseStr(p.k[6:8])
	if err != nil {
		p.err = err
		return
	}
	minute, err := utils.ReverseStr(p.k[8:10])
	if err != nil {
		p.err = err
		return
	}
	second, err := utils.ReverseStr(p.k[10:12])
	if err != nil {
		p.err = err
		return
	}
	timezone, err := utils.ReverseStr(p.k[12:14])
	if err != nil {
		p.err = err
		return
	}
	timezoneInt, err := strconv.ParseInt(timezone, 10, 32)
	if err != nil {
		p.err = err
		return
	}
	timezoneInt = timezoneInt - 24

	p.time, err = zone.Parse(time.RFC3339, fmt.Sprintf("20%s-%s-%sT%s:%s:%s+%02d:00", year, month, day, hour, minute, second, timezoneInt))
	if err != nil {
		p.err = err
		return
	}
	return
}
func (p *pdu) parseContent7() error {
	defaultList := []string{
		"@", "£", "$", "¥", "è", "é", "ù", "ì", "ò", "Ç", "\n", "Ø", "ø", "\r", "Å", "å", "\u0394", "_", "\u03a6", "\u0393", "\u039b", "\u03a9", "\u03a0", "\u03a8", "\u03a3", "\u0398", "\u039e", "€", "Æ", "æ", "ß", "É", " ", "!", "\"", "#", "¤", "%", "&", "'", "(", ")", "*", "+", ",", "-", ".", "/", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", ":", ";", "<", "=", ">", "?", "¡", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "Ä", "Ö", "Ñ", "Ü", "§", "¿", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "ä", "ö", "ñ", "ü", "à",
	}

	bin := ""
	for i := 0; i < len(p.m); i += 2 {
		hex := p.m[i : i+2]
		utf8Hex, _ := strconv.ParseInt(hex, 16, 32)
		bin += fmt.Sprintf("%08b", utf8Hex)
	}

	s := 1
	count := 0
	var octetList []string
	var restList []string
	var septetsList []string
	for i := 0; i < len(bin); i += 8 {
		octetList = append(octetList, bin[i:i+8])
		restList = append(restList, octetList[count][0:s%8])
		septetsList = append(septetsList, octetList[count][s%8:8])

		s++
		count++
		if s >= 8 {
			s = 1
		}
	}

	matchCount := 0 // AJA
	content := ""
	for i := 0; i < len(restList); i++ {
		if i%7 == 0 {
			if i != 0 {
				index, err := bin2Int(restList[i-1])
				if err != nil {
					return err
				}
				content += defaultList[index]
				matchCount++
			}

			index, err := bin2Int(septetsList[i])
			if err != nil {
				return err
			}
			content += defaultList[index]
			matchCount++
		} else {
			index, err := bin2Int(septetsList[i] + restList[i-1])
			if err != nil {
				return err
			}
			content += defaultList[index]
			matchCount++
		}
	}

	p.content = content
	p.length = uint(len(content))
	return nil
}
func bin2Int(bin string) (int64, error) {
	i, err := strconv.ParseInt(bin, 2, 32)
	if err != nil {
		return 0, err
	}
	return i, nil
}
func (p *pdu) parseContent8() error {
	// @todo +CMGR: 0,\"\",26\r\n0891683108200945F5240D91681234551820F100009180122024432306E170381C0E03
	// @todo  echo 'AT+CMGDA=6\r\n' > /dev/ttyUSB1
	b := []byte{}
	for i := 0; i < len(p.m); i += 2 {
		hex := p.m[i : i+2]
		utf8Hex, _ := strconv.ParseInt(hex, 16, 32)

		b = append(b, byte(utf8Hex))
	}
	p.content = string(b[:])
	p.length = uint(len(b))
	return nil
}
func (p *pdu) parseContent16() error {
	b := []byte{}
	for i := 0; i < len(p.m); i += 4 {
		hex1 := p.m[i : i+2]
		utf16A, _ := strconv.ParseInt(hex1, 16, 32)

		b = append(b, byte(utf16A))

		hex2 := p.m[i+2 : i+4]
		utf16B, _ := strconv.ParseInt(hex2, 16, 32)

		b = append(b, byte(utf16B))

		// utf32 := utf16A * 256+utf16B
		// fmt.Print(string(utf32))
	}

	result, length, err := transform.Bytes(unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder(), b)
	if err != nil {
		return err
	}

	p.content = string(result[:])
	p.length = uint(length)
	return nil
}

func (p *pdu) smsCenterLen() uint {
	length, err := strconv.ParseUint(p.a, 16, 32)
	if err != nil {
		p.err = err
		return 0
	}
	return uint(length) * 2
}
func (p *pdu) senderLen() uint {
	length, err := strconv.ParseUint(p.f, 16, 32)
	if err != nil {
		p.err = err
		return 0
	}
	if length%2 != 0 {
		length += 1
	}
	return uint(length)
}
func (p *pdu) contentLen() uint {
	length, err := strconv.ParseUint(p.l, 16, 32)
	if err != nil {
		p.err = err
		return 0
	}
	return uint(length)
}

func (p *pdu) Error() error {
	return p.err
}
