package uhaul

import "testing"

func TestPack(t *testing.T) {
	testData := map[string][]interface{}{
		// b'\x01\x00\x02\x00\x03\x00\x00\x00'
		"ccc": {1, 2, 3},
		"hhh": {10000, 20000, 30000},
		// largest uint16, uint16, and uint32
		"hhl": {65535, 65535, 4294967295},
		// struct.pack("I", 200000) == b'@\r\x03\x00'
		// 64, 13, 3, 0
		"I": {200000},
	}

	for k, v := range testData {
		packed, err := Pack(k, v...)
		if err != nil {
			t.Fail()
		}

		sum, _, err := CalcSize(k)
		if err != nil {
			t.Fail()
		}

		if len(packed) != sum {
			t.Fail()
		}

		unpacked, err := Unpack(k, packed)
		if err != nil {
			t.Fail()
		}

		for i, j := range v {
			if !compare(j, unpacked[i]) {
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

		// are pack and unpack perfect inverses?
		packed, err = Pack("I256s", 70, v)
		if err != nil {
			t.Fail()
		}

		unpacked, err = Unpack("I256s", packed)
		if err != nil {
			t.Fail()
		}

		if !compare(packed[0], unpacked[0]) {
			t.Fail()
		}

		var built string
		for _, v := range unpacked[1:] {
			built += string(v.(byte))
		}

	}
}

// shitty helper function for comparing integer interfaces
func compare(x, y interface{}) bool {
	var (
		w int64
		z int64
	)

	switch x.(type) {
	case int32:
		w = int64(x.(int32))
	case int64:
		w = int64(x.(int64))
	case int:
		w = int64(x.(int))
	case uint8:
		w = int64(x.(uint8))
	case uint16:
		w = int64(x.(uint16))
	case uint32:
		w = int64(x.(uint32))
	default:
		panic("Unknown type")
	}

	switch y.(type) {
	case int32:
		z = int64(y.(int32))
	case int64:
		z = int64(y.(int64))
	case int:
		z = int64(y.(int))
	case uint8:
		z = int64(y.(uint8))
	case uint16:
		z = int64(y.(uint16))
	case uint32:
		z = int64(y.(uint32))

	default:
		panic("Unknown type")
	}

	return w == z
}
