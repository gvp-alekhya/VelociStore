package core

import "github.com/gvp-alekhya/VelociStore/config"

func Evict() {
	evictionStrategy := config.EVICTION_STRATEGY
	switch evictionStrategy {
	case "simple-first":
		evictFirst()
	case "allkeys-random":
		{
			evictAllKeysRandom()
		}
	}

}

func evictFirst() {
	for key := range RedisStore {
		Del(key)
		return
	}
}

func evictAllKeysRandom() {
	evictCount := int64(config.EVICTION_RATIO * config.MAX_KEY_LIMIT)
	for key := range RedisStore {
		if evictCount <= 0 {
			return
		}
		Del(key)
		evictCount--
	}
}
