package core

import (
	"errors"
	"fmt"
	"strings"
)

func Decode(data []byte) (i interface{}, pos int, err error) {
	if len(data) == 0 {
		return nil, -1, errors.New("no data")
	}
	switch data[0] {
	case '+':
		return DecodeSimpleString(data)
	case '-':
		return DecodeError(data)
	case '$':
		return DecodeBulkString(data)
	case ':':
		return DecodeInteger(data)
	case '#':
		return DecodeBoolean(data)
	case '*':
		return DecodeArray(data)
	}
	return data, 0, nil
}

func DecodeSimpleString(data []byte) (string, int, error) {
	pos := 1
	for data[pos] != '\r' {
		pos++
	}
	return string(data[:pos]), pos + 2, nil
}

func DecodeBulkString(data []byte) (string, int, error) {
	pos := 1
	len, delta := readLength(data[1:])
	pos += delta

	fmt.Println("DecodeBulkString :: len :: pos :: data", len, pos, string(data))
	return string(data[pos:(pos + len)]), pos + len + 2, nil
}

func DecodeInteger(data []byte) (int64, int, error) {
	var n int64 = 0
	pos := 1
	for data[pos] != '\r' {
		n = n*10 + int64(data[pos]-'0')
		pos++
	}
	return n, pos + 2, nil
}

func DecodeError(data []byte) (string, int, error) {
	str, pos, err := DecodeSimpleString(data)
	return str, pos, err
}
func DecodeBoolean(data []byte) (boolean bool, pos int, err error) {
	str, _, _ := DecodeSimpleString(data)
	isBoolean := (strings.ToLower(str) == "true")
	return isBoolean, 0, nil
}
func DecodeArray(data []byte) (interface{}, int, error) {
	// first character *
	index := 1
	// reading the length
	count, currRead := readLength(data[index:])
	fmt.Println("DecodeArray :: count :: currRead", count, currRead, string(data))
	index += currRead
	var elems []interface{} = make([]interface{}, count)
	for i := range elems {
		fmt.Println("DecodeArray In loop :: data[index:] :: currRead", index, string(data[index:]))
		elem, currRead, err := Decode(data[index:])
		fmt.Println("DecodeArray In loop :: elem :: currRead", elem, currRead)
		if err != nil {
			return nil, 0, err
		}
		elems[i] = elem
		index += currRead
	}
	fmt.Println("DecodeArray Return:: elem :: index", elems, index)
	return elems, index, nil
}

func readLength(data []byte) (int, int) {
	pos, length := 0, 0
	for ; pos < len(data); pos++ {
		b := data[pos]
		if !(b >= '0' && b <= '9') {
			return length, pos + 2
		}
		fmt.Println("Read length :: len :: ", length, int(b-'0'))
		length = length*10 + int(b-'0')
	}
	return 0, 0
}

func DecodeCommands(data []byte) ([]interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("no data")
	}
	var values []interface{} = make([]interface{}, 0)
	index := 0
	for index < len(data) {
		value, delta, err := Decode(data[index:])
		if err != nil {
			return values, err
		}
		index += delta
		values = append(values, value)
	}
	fmt.Println("Decode commands :: ", values)
	return values, nil
}
