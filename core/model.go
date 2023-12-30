package core

import (
	"time"

	"github.com/gvp-alekhya/VelociStore/config"
)

var RedisStore map[string]*Obj //Declaring hash map to save keys of type strings and values as pointers object type

// init method in Go is called once during program execution automatically
func init() {
	RedisStore = make(map[string]*Obj)
}

func NewObj(value interface{}, expirationInMs int64, objType uint8, oEncoding uint8) *Obj {
	var expiresAt int64 = -1
	if expirationInMs > 0 {
		expiresAt = time.Now().UnixMilli() + expirationInMs
	}
	nobj := Obj{Value: value, ExpirationInMs: expiresAt, TypeEncoding: objType | oEncoding}
	return &nobj
}

func Put(key string, value *Obj) {
	if len(RedisStore) == config.MAX_KEY_LIMIT {
		Evict()
	}
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
