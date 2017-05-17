package uhaul

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

func Pack(format string, vals ...interface{}) ([]byte, error) {
	alignment, err := alignment(format[0])
	if err == nil {
		format = format[1:]
	}

	_, sizes, err := CalcSize(format)
	if err != nil {
		return nil, err
	}

	if len(sizes) != len(vals) {
		return nil, fmt.Errorf("pack expected %v for packing (got %v)", len(sizes), len(vals))
	}

	data := []byte{}
	buf := new(bytes.Buffer)
	for i, v := range vals {
		buf.Reset()

		switch sizes[i] {
		case 1:
			binary.Write(buf, alignment, uint8(v.(int)))
		case 2:
			binary.Write(buf, alignment, uint16(v.(int)))
		case 4:
			binary.Write(buf, alignment, uint32(v.(int)))
		case 8:
			binary.Write(buf, alignment, uint64(v.(int)))
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
	return data, nil
}

// returns alignment followed by an error. if the error exists, then
// an alignment specifier was not given.
func alignment(format byte) (binary.ByteOrder, error) {
	switch format {
	case '<':
		return binary.LittleEndian, nil
	case '>', '!':
		return binary.BigEndian, nil
	default:
		return binary.LittleEndian, errors.New("")
	}
}

// for each value n in sizes, splice the source into increments of that size
func splitSlice(source []byte, sizes []int) [][]byte {
	split := make([][]byte, len(sizes))
	for i, v := range sizes {
		split[i] = source[:v]
		source = source[v:]
	}
	return split
}

// if you mix strings and regular data types in the format string,
// there's currently no way to separate the string data from the other
// packed data
func Unpack(format string, vals []byte) ([]interface{}, error) {
	alignment, err := alignment(format[0])
	if err == nil {
		format = format[1:]
	}

	sum, sizes, err := CalcSize(format)
	if err != nil {
		return nil, err
	}

	if len(vals) != sum {
		return nil, fmt.Errorf("unpack requires a []byte of length %v", sum)
	}

	split := splitSlice(vals, sizes)

	// something's fishy about this data variable
	data := []interface{}{}
	for i, v := range split {
		// string handling
		// this only handles strings of length 8 or more
		if len(v) > 8 {
			for _, j := range v {
				if j == 0x00 {
					continue
				}
				var val byte
				binary.Read(bytes.NewReader([]byte{j}), alignment, &val)
				data = append(data, val)
			}
		} else {
			buf := bytes.NewReader(v)

			switch sizes[i] {
			case 1:
				var val uint8
				binary.Read(buf, alignment, &val)
				data = append(data, val)
			case 2:
				var val uint16
				binary.Read(buf, alignment, &val)
				data = append(data, val)
			case 4:
				var val uint32
				binary.Read(buf, alignment, &val)
				data = append(data, val)
			case 8:
				var val uint64
				binary.Read(buf, alignment, &val)
				data = append(data, val)
			default:
				return nil, errors.New("unknown format size")
			}
		}
	}
	return data, nil
}

// Returns the size of the format string in bytes,
// and a slice of the size of the interpreted format identifiers
// this only handles strings of bytes, not strings of characters
func CalcSize(format string) (int, []int, error) {
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
				return 0, nil, err
			}

			sum += s
			i += sPos
			argCount = append(argCount, s)
		default:
			return 0, nil, fmt.Errorf("Unknown formatting verb %v", string(c))
		}
	}
	return sum, argCount, nil
}
