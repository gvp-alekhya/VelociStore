package core

var KeySpaceStats [4]map[string]int

func UpdateKeySpaceStats(dbNum int, metric string, val int) {
	KeySpaceStats[dbNum][metric] = val
}
