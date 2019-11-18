package phone

import (
	"testing"
)

func TestPdu_Scan16(t *testing.T) {
	pdu := NewPDU()
	pdu.Scan("0891683108200905F5240FA101960119145036F90008918091110431237C3010006C00750063006B0069006E00200063006F006600660065006530116211731C4F6060F3559D70B94EC04E48FF0C90014F600035002E003562985168573A996E54C15238FF0C5168573A996E54C1768653EF4F7F7528FF5E53BB004100500050002F5C0F7A0B5E8F559D4E00676F002056DE0054004490008BA2")
	if pdu.center != "+8613800290505" {
		t.Error("pdu.center error")
	}
	if pdu.from != "106910914105639" {
		t.Error("pdu.from error")
	}
	//debug.DD(pdu)
	if pdu.content != "【luckin coffee】我猜你想喝点什么，送你5.5折全场饮品券，全场饮品皆可使用～去APP/小程序喝一杯 回TD退订" {
		t.Error("pdu.content error")
	}
}
func TestPdu_Scan7(t *testing.T) {
	pdu := NewPDU()
	pdu.Scan("0891683108200945F5240D91683123151820F500009111712291242308E8329BFD3EBFC9")
	if pdu.center != "+8613800290545" {
		t.Error("pdu.center error")
	}
	if pdu.from != "+8613325181025" {
		t.Error("pdu.from error")
	}
	//debug.DD(pdu)
	if pdu.content != "hellogod" {
		t.Error("pdu.content error")
	}
}
func TestPdu_Scan8(t *testing.T) {
	pdu := NewPDU()
	pdu.Scan("0891683108200945F5240D91681234551820F100009180122024432306E170381C0E03")
	if pdu.center != "+8613800290545" {
		t.Error("pdu.center error")
	}
	if pdu.from != "+8621435581021" {
		t.Error("pdu.from error")
	}
	//debug.DD(pdu)
	if pdu.content != "aaaaaa" {
		t.Error("pdu.content error")
	}
}