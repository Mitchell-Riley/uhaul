package uhaul

import (
	"fmt"
	"testing"
)

var testData = []struct {
	format   string        // format string needed to pack/unpack
	packed   []byte        // the data emitted from a Pack
	unpacked []interface{} // the raw data given to Pack/output of Unpack
}{
	// examples from the python documentation
	{"!hhl", []byte{0, 1, 0, 2, 0, 0, 0, 3}, []interface{}{1, 2, 3}},
	// {"llh0l", []byte{}, []interface{}{1, 2, 3}},

	// my own examples
	// "ccc": {1, 2, 3},
	// "hhh": {10000, 20000, 30000},
	// // largest uint16, uint16, and uint32
	// "hhl": {65535, 65535, 4294967295},
	// // struct.pack("I", 200000) == b'@\r\x03\x00'
	// // 64, 13, 3, 0
	// "I": {200000},
}

func TestPack(t *testing.T) {
	for _, v := range testData {
		packed, err := Pack(v.format, v.unpacked...)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		}

		// if the packed data does not match the correct slice, then fail
		for i := range packed {
			if !compare(packed[i], v.packed[i]) {
				fmt.Println("packed wrong")
				t.Fail()
			}
		}
	}
}

func TestUnpack(t *testing.T) {
	for _, v := range testData {
		unpacked, err := Unpack(v.format, v.packed)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		}

		for i := range unpacked {
			if !compare(unpacked[i], v.unpacked[i]) {
				fmt.Println("unpacked wrong")
				t.Fail()
			}
		}
	}
}

func TestCalcSize(t *testing.T) {
	t.Skip()
	var testData = []struct {
		format string
		size   int
	}{
		{"ci", 8},
		{"ic", 5},
	}

	for _, v := range testData {
		sum, _, _ := CalcSize(v.format)
		if sum != v.size {
			// t.Fail()
		}
	}
}

func TestStringUnpack(t *testing.T) {
	t.Skip()
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
		if compare(packed[0], unpacked[0]) == false {
			t.Fail()
		}

		var built string
		for _, v := range unpacked {
			built += string(v.(byte))
		}

		if built != v {
			t.Fail()
		}
	}
}

func TestAlignment(t *testing.T) {
	t.Skip()
	p, err := Pack(">I", 2032480932)
	if err != nil {
		t.Fatal(err)
	}

	record := append([]byte("raymond   "), []byte{0x32, 0x12, 0x08, 0x01, 0x08}...)
	d, err := Unpack("<10sHHb", record)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(p)
	fmt.Println(d)
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
		fmt.Printf("Unknown type: %T\n", x)
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
		fmt.Printf("Unknown type: %T\n", x)
	}

	return w == z
}
