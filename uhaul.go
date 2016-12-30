package uhaul

import (
	"bytes"
	"encoding/binary"
	"log"
	"strconv"
	"unicode"
)

func Pack(format string, vals ...interface{}) []byte {
	_, sizes := CalcSize(format)
	if len(sizes) != len(vals) {
		log.Fatal("argument count mismatch")
	}

	data := []byte{}
	buf := new(bytes.Buffer)
	for i, v := range vals {
		buf.Reset()

		switch sizes[i] {
		case 1:
			binary.Write(buf, binary.LittleEndian, uint8(v.(int)))
		case 2:
			binary.Write(buf, binary.LittleEndian, uint16(v.(int)))
		case 4:
			binary.Write(buf, binary.LittleEndian, uint32(v.(int)))
		case 8:
			binary.Write(buf, binary.LittleEndian, uint64(v.(int)))
		//string type
		default:
			byteCount := sizes[i]
			stringValue := v.(string)
			// string + however many empty bytes are needed to fill
			d := append([]byte(stringValue), make([]byte, byteCount-len(stringValue))...)
			data = append(data, d...)
		}
		data = append(data, buf.Bytes()...)
	}
	return data
}

// for each value n of sizes, splice the source into increments of that size
func splitSlice(source []byte, sizes []int) [][]byte {
	split := make([][]byte, len(sizes))
	for i, v := range sizes {
		split[i] = source[:v]
		source = source[v:]
	}
	return split
}

// Somehow check if the format string and vals are of the same "length"
func Unpack(format string, vals []byte) []byte {
	_, sizes := CalcSize(format)
	// if byteSize/8 != len(vals) {
	// 	fmt.Println(len(sizes), len(vals))
	// 	log.Fatal("argument count mismatch")
	// }

	split := splitSlice(vals, sizes)

	data := []byte{}
	var val byte
	for _, v := range split {

		// string handling
		// this only handles strings of length 8 or more
		if len(v) > 8 {
			for _, j := range v {
				if j == 0x00 {
					continue
				}
				binary.Read(bytes.NewReader([]byte{j}), binary.LittleEndian, &val)
				data = append(data, val)
			}
		} else {
			buf := bytes.NewReader(v)
			binary.Read(buf, binary.LittleEndian, &val)
			data = append(data, val)
		}
	}
	return data
}

// Returns the size of the format string in bytes,
// and a slice of the size of the interpreted format identifiers
// this only handles strings of bytes, not strings of characters
func CalcSize(format string) (int, []int) {
	sum := 0
	argCount := []int{}
	for i := 0; i < len(format); i++ {
		switch c := format[i]; {
		case c == 'c', c == 'b', c == 'B', c == '?':
			sum += 1
			argCount = append(argCount, 1)
		case c == 'h', c == 'H':
			sum += 2
			argCount = append(argCount, 2)
		case c == 'i', c == 'I', c == 'l', c == 'L', c == 'f':
			sum += 4
			argCount = append(argCount, 4)
		case c == 'q', c == 'Q', c == 'd':
			sum += 8
			argCount = append(argCount, 8)
		case unicode.IsNumber(rune(c)):
			var sPos int
			for j, k := range format[i:] {
				if k == 's' {
					sPos = j
					break
				}
			}

			s, err := strconv.Atoi(format[i : sPos+i])
			if err != nil {
				log.Fatal(err)
			}

			sum += s
			i += sPos
			argCount = append(argCount, s)
		}
	}
	return sum, argCount
}
