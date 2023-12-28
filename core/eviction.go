package core

func Evict() {
	evictFirst()
}

func evictFirst() {
	for key := range RedisStore {
		delete(RedisStore, key)
		return
	}
}
