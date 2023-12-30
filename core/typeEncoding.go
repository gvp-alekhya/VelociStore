package core

import (
	"errors"
	"fmt"
	"strconv"
)

func getType(te uint8) uint8 {
	return te & 0b11110000
}
func getEncoding(te uint8) uint8 {
	return te & 0b00001111
}

func assertType(te uint8, t uint8) error {
	if getType(te) != t {
		return errors.New("(error) Operation is not permitted on this type")
	}
	return nil
}
func assertEncoding(te uint8, t uint8) error {
	fmt.Println("te ::", te, t)
	if getEncoding(te) != t {
		return errors.New("(error) Operation is not permitted on this type")
	}
	return nil
}

func DeduceTypeEncoding(value interface{}) (uint8, uint8) {
	objType := OBJ_TYPE_STRING
	_, err := strconv.ParseInt(value.(string), 10, 64)
	if err == nil {
		return objType, OBJ_ENCODING_INT
	} else if len(value.(string)) <= 44 {
		return objType, OBJ_ENCODING_EMBSTR
	}

	return objType, OBJ_ENCODING_RAW
}
