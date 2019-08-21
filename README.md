# SMS Serial Chip Forwarder

* sim800c [I bought here](https://item.taobao.com/item.htm?id=584802370108) ¥44.50

```bash
> ./sms-notifier sms:read /dev/ttyUSB1
INFO[2019-08-21T02:44:49+08:00] Send Bytes                                    bytes="ATE0\r\n" length=6
INFO[2019-08-21T02:44:49+08:00] Send Bytes                                    bytes="AT+CMEE=1\r\n" length=11
INFO[2019-08-21T02:44:49+08:00] Send Bytes                                    bytes="AT+CMGF=0\r\n" length=11
INFO[2019-08-21T02:44:49+08:00] Listening                                     com_port=/dev/ttyUSB1
INFO[2019-08-21T02:44:49+08:00] Incoming data                                 data=OK
```

## How To Use
1. Plug in your USB chip
2. find your serial port name -> `ls -l /dev/ttyUSB*`
```bash
> ls -l /dev/ttyUSB*
crw-rw----. 1 root dialout 188, 1 8月  21 02:47 /dev/ttyUSB0
```
3. Copy `.env.example.json` -> `.env.json`
4. Save your `pushover` config in `.env.json`
5. `go run artisan sms:read /dev/YOUR-SERIAL`, such as `/dev/ttyUSB0`

> Binary `WIP`

## Implement your own NOTIFIER
`/app/logics/phone/interfaces/notifier.go`  
```go
package interfaces

type Notifier interface {
	Notify(sender, content string) error
}
```

## Thanks
* http://www.sendsms.cn/pdu/
