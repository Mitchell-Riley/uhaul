package uhaul

import (
	"fmt"
	"testing"
)

func TestPack(t *testing.T) {
	testData := [][]interface{}{
		// b'\x01\x00\x02\x00\x03\x00\x00\x00'
		{1, 2, 3},
		{10000, 20000, 30000},
		// largest uint16, uint16, and uint32
		{65535, 65535, 4294967295},
		// struct.pack("I", 200000) == b'@\r\x03\x00'
		// 64, 13, 3, 0
		{200000},
	}

	testFormats := []string{
		"ccc",
		"hhh",
		"hhl",
		"I",
	}

	for i, v := range testData {
		packed, err := Pack(testFormats[i], v...)
		if err != nil {
			t.Fail()
		}

		sum, _, err := CalcSize(testFormats[i])
		if err != nil {
			t.Fail()
		}

		if len(packed) != sum {
			t.Fail()
		}

		unpacked, _ := Unpack(testFormats[i], packed)

		for i, j := range v {
			// convert to strings here because honestly idk
			if fmt.Sprint(j) != fmt.Sprint(unpacked[i]) {
				t.Fail()
			}
		}
	}
}

func TestStringUnpack(t *testing.T) {
	stringSources := []string{
		"0x32040239",
		"89.-2309e3823uhefwo92 98y",
		"Hello, World!",
	}

	for _, v := range stringSources {
		packed, err := Pack("256s", v)
		if err != nil {
			t.Fail()
		}

		unpacked, err := Unpack("256s", packed)
		if err != nil {
			t.Fail()
		}

		for i, v := range packed[:len(v)] {
			if v != unpacked[i] {
				t.Fail()
			}
		}
	}

	for _, v := range stringSources {
		packed, err := Pack("I256s", 70, v)
		if err != nil {
			t.Fail()
		}

		unpacked, err := Unpack("I256s", packed)
		if err != nil {
			t.Fail()
		}

		// convert to strings here because honestly idk
		if fmt.Sprint(packed[0]) != fmt.Sprint(unpacked[0]) {
			t.Fail()
		}

		var built string
		for _, v := range unpacked[1:] {
			built += string(v.(byte))
		}

		if v != built {
			t.Fail()
		}
	}
}
