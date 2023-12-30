package core

/*Go map declaration map[keyType]valueType*/
type Obj struct {
	TypeEncoding   uint8 //First 4 bits is for type and last 4 bits for encoding
	Value          interface{}
	ExpirationInMs int64
}

var OBJ_TYPE_STRING uint8 = 0 << 1 //left shoft 4 bits to indicate the type

var OBJ_ENCODING_RAW uint8 = 0
var OBJ_ENCODING_INT uint8 = 1
var OBJ_ENCODING_EMBSTR uint8 = 8
