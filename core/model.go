package core

import (
	"time"
)

var RedisStore map[string]*Obj //Declaring hash map to save keys of type strings and values as pointers object type
/*Go map declaration map[keyType]valueType*/

type Obj struct {
	Value          interface{}
	ExpirationInMs int64
}

// init method in Go is called once during program execution automatically
func init() {
	RedisStore = make(map[string]*Obj)
}

func NewObj(value interface{}, expirationInMs int64) *Obj {
	var expiresAt int64 = -1
	if expirationInMs > 0 {
		expiresAt = time.Now().UnixMilli() + expirationInMs
	}
	nobj := Obj{Value: value, ExpirationInMs: expiresAt}
	return &nobj
}

func Put(key string, value *Obj) {
	RedisStore[key] = value
}

func Get(key string) *Obj {
	return RedisStore[key]
}

func Del(key string) bool {
	if RedisStore[key] != nil {
		delete(RedisStore, key)
		return true
	}
	return false
}
