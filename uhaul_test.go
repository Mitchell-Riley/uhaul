package uhaul

import (
	"fmt"
	"testing"
)

func TestPack(t *testing.T) {
	// b'\x01\x00\x02\x00\x03\x00\x00\x00'
	packed, err := Pack("hhl", 1, 2, 3)
	if err != nil {
		t.Fail()
	}

	sum, _, err := CalcSize("hhl")
	if err != nil {
		t.Fail()
	}

	if len(packed) != sum {
		t.Fail()
	}

	// struct.pack("I", 200000) == b'@\r\x03\x00'
	// 64, 13, 3, 0
	packed, _ = Pack("I", 200000)
	unpacked, _ := Unpack("I", packed)
	// convert to strings here because honestly idk
	if fmt.Sprint(unpacked[0]) != fmt.Sprint(200000) {
		t.Fail()
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
