package core

import (
	"strings"
	"fmt"
	"errors"
)

func Decode(data []byte) (i interface{}, pos int, err error) {
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

func DecodeSimpleString(data []byte) ( string,  int,  error) {
	pos:=1
	for data[pos]!='\r'{
		pos++
	}
    return string(data[:pos]), pos+2, nil
}

func DecodeBulkString(data []byte) ( string,  int,  error) {
	pos := 1
	len,delta := readLength(data[1:])
	pos += delta
	fmt.Printf("DecodeBulkString Read Length %d:: pos :: %d :: string %s\n", len, pos, string(data[pos:pos+len]))
    return string(data[pos:pos+len]), pos+1+len+2, nil
}

func DecodeInteger(data []byte) ( int64,  int,  error) {
	var n int64 = 0;
	pos := 1
	for data[pos]!='\r'{
		n = n*10 + int64(data[pos] - '0')
		pos++
	}
    return n, pos+2, nil
}

func DecodeError(data []byte) ( string,  int,  error) {
	str, pos, err := DecodeSimpleString(data)
    return str, pos, err
}
func DecodeBoolean(data []byte) (boolean bool, pos int, err error) {
	str,_,_ := DecodeSimpleString(data)
	isBoolean := (strings.ToLower(str) == "true")
    return isBoolean, 0, nil
}
func DecodeArray(data []byte) ( interface{}, int,  error) {
	// first character *
	var index int= 1
	// reading the length
	count, currRead := readLength(data[index:])
	fmt.Printf("DecodeArray Read Outside loop Length %d:: Delta :: %d\n", count, currRead)
	index += currRead
	fmt.Printf("DecodeArray Read Outside loop index %d\n", index)
	var elems []interface{} = make([]interface{}, count)
	for i := range elems {
		elem, currRead, err := Decode(data[index:])
		fmt.Printf("DecodeArray Read Element %T :: diff %d:: index :: %d\n", elem, currRead, index)
		if err != nil {
			return nil, 0, err
		}
		elems[i] = elem
		index += currRead
	}
	return elems, index, nil
}

func readLength(data []byte) (int, int) {
	pos, length := 0, 0
	for pos = range data {
		b := data[pos]
		if !(b >= '0' && b <= '9') {
			return length, pos+2
		}
		length = length*10 + int(b-'0')
	}
	return 0, 0
}


func DecodeArrayString(data []byte) ([]string, error) {
    fmt.Printf("DecodeArrayString data %s\n", string(data))
    value, _, err := Decode(data)
    if err != nil {
        return nil, err
    }
	stringValue := fmt.Sprintf("%v", value)
	fmt.Printf("Decoded data %s\n", stringValue)

    ts, ok := value.([]interface{})
    if !ok {
        return nil, errors.New("unexpected type")
    }

    tokens := make([]string, len(ts))
    for i, v := range ts {
        switch t := v.(type) {
        case string:
            tokens[i] = t
        case []uint8:
            tokens[i] = string(t)
        default:
            return nil, fmt.Errorf("unexpected type in array: %T", v)
        }
    }

    return tokens, nil
}
