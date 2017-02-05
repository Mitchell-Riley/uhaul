package uhaul

import "testing"

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
}

func TestStringUnpack(t *testing.T) {
	stringSources := []string{
		"meow",
		"llama",
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

		if packed[0] != unpacked[0] {
			t.Fail()
		}
	}
}
