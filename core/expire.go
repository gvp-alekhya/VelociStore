package core

import (
	"fmt"
	"time"
)

func DeleteSample() float64 {
	limit := 20
	expiredKeys := 0

	for key, value := range RedisStore {
		if value.ExpirationInMs != -1 {
			limit--
			if value.ExpirationInMs <= time.Now().UnixMilli() {
				expiredKeys++
				Del(key)
			}

		}
		if limit == 0 {
			break
		}

	}
	return float64(expiredKeys) / float64(limit)
}

func DeleteExpiredKeys() {
	for {
		if DeleteSample() < 0.25 {
			break
		}
	}
	fmt.Print("Deleted expired Keys, total keys now :: ", len(RedisStore))
}
