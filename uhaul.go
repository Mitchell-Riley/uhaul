package uhaul

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

// define the standard size of packed values in bytes
// these can be changed by the user
var (
	CHAR  = 1
	SCHAR = 1
	UCHAR = 1

	BOOL = 1

	SHORT  = 2
	USHORT = 2

	INT  = 4
	UINT = 4

	LONG      = 4
	ULONG     = 4
	LONGLONG  = 8
	ULONGLONG = 8

	FLOAT  = 4
	DOUBLE = 8
)

func Pack(format string, vals ...interface{}) ([]byte, error) {
	endian, err := determineEndianness(format[0])
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
			binary.Write(buf, endian, uint8(v.(int)))
		case 2:
			binary.Write(buf, endian, uint16(v.(int)))
		case 4:
			binary.Write(buf, endian, uint32(v.(int)))
		case 8:
			binary.Write(buf, endian, uint64(v.(int)))
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

// returns endian followed by an error. if the error exists, then
// an endian specifier was not given.
func determineEndianness(format byte) (binary.ByteOrder, error) {
	switch format {
	case '<':
		return binary.LittleEndian, nil
	case '>', '!':
		return binary.BigEndian, nil
	default:
		// default to little endian
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
	endian, err := determineEndianness(format[0])
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
				binary.Read(bytes.NewReader([]byte{j}), endian, &val)
				data = append(data, val)
			}
		} else {
			buf := bytes.NewReader(v)

			switch sizes[i] {
			case 1:
				var val uint8
				binary.Read(buf, endian, &val)
				data = append(data, val)
			case 2:
				var val uint16
				binary.Read(buf, endian, &val)
				data = append(data, val)
			case 4:
				var val uint32
				binary.Read(buf, endian, &val)
				data = append(data, val)
			case 8:
				var val uint64
				binary.Read(buf, endian, &val)
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
		typeSize := tokenMap[rune(format[i])]
		sum += typeSize
		argCount = append(argCount, typeSize)

		if unicode.IsNumber(rune(format[i])) {
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
		}
	}
	return sum, argCount, nil
}

func tokenSplit(fmt string) []string {
	tokens := []string{}
	for i := 0; i < len(fmt); i++ {
		if unicode.IsNumber(rune(fmt[i])) {
			tokens = append(tokens, string(fmt[i:i+2]))
			i++
		} else {
			tokens = append(tokens, string(fmt[i]))
		}
	}
	return tokens
}

var tokenMap = map[rune]int{
	'c': CHAR,
	'b': SCHAR,
	'B': UCHAR,
	'?': BOOL,
	'h': SHORT,
	'H': USHORT,
	'i': INT,
	'I': UINT,
	'l': LONG,
	'L': ULONG,
	'q': LONGLONG,
	'Q': ULONGLONG,
	'f': FLOAT,
	'd': DOUBLE,
}
