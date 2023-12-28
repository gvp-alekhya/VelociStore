package config

import "time"

const (
	SERVER_HOST = "0.0.0.0"
	SERVER_PORT = 2929
	SERVER_TYPE = "tcp"
	MAX_CLIENTS = 20000
)

var CRON_FREQUENCY time.Duration = 1 * time.Second
var CRON_LAST_EXECUTED_TIME time.Time = time.Now()
