package core

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/gvp-alekhya/VelociStore/config"
)

func evalCommand(args []string, c io.ReadWriter) error {
	var b []byte

	if len(args) >= 2 {
		return errors.New("ERR wrong number of arguments for 'ping' command")
	}

	if len(args) == 0 {
		b = Encode("PONG", false)
	} else {
		b = Encode(args[0], false)
	}
	_, err := c.Write(b)
	return err
}
func evalSet(args []string, c io.ReadWriter) error {
	if len(args) <= 1 {
		return errors.New("ERR wrong number of arguments for 'SET' command")
	}

	var key, value string
	var expirationInMs int64 = -1

	key = args[0]
	value = args[1]

	for i := 2; i < len(args); i++ {
		lowercaseArgs := strings.ToLower(args[i])
		switch lowercaseArgs {
		case "ex":
			if i == len(args) {
				return errors.New("ERR syntax error")
			}
			i++
			expirationInSec, err := strconv.ParseInt(args[i], 10, 64) //converting decimal to 64 bit int
			if err != nil {
				return errors.New("ERR invalid expiration value passed in 'SET' command")
			}
			expirationInMs = expirationInSec * 1000

		default:
			return errors.New("ERR invalid arguments for 'SET' command")
		}
	}
	Put(key, NewObj(value, expirationInMs))
	c.Write([]byte(config.OKResponse))
	return nil
}
func evalGet(args []string, c io.ReadWriter) error {
	if len(args) != 1 {
		return errors.New("ERR wrong number of arguments for 'GET' command")
	}

	key := args[0]
	Obj := Get(key)
	if Obj == nil {
		c.Write([]byte(config.NILResponse))
		return errors.New("no key found")
	}
	if Obj.ExpirationInMs != -1 && Obj.ExpirationInMs <= time.Now().UnixMilli() {
		Del(key)
		c.Write([]byte(config.NILResponse))
		return nil
	}
	c.Write(Encode(Obj.Value, true))

	return nil
}
func evalTTL(args []string, c io.ReadWriter) error {

	if len(args) != 1 {
		return errors.New("ERR wrong number of arguments for 'TTL' command")
	}
	key := args[0]

	Obj := Get(key)
	if Obj == nil {
		c.Write([]byte(config.NoKeyResponse))
		return errors.New("ERR no key found")
	}
	currentTimeInMS := time.Now().UnixMilli()
	if Obj.ExpirationInMs == -1 {
		c.Write([]byte(config.NoExpiryResponse))
		return errors.New("no expiration set for key")
	}

	var timeLeft = Obj.ExpirationInMs - currentTimeInMS
	if timeLeft < 0 {
		c.Write([]byte(config.NoKeyResponse))
		return errors.New("key expired")
	}
	c.Write(Encode(int64(timeLeft/1000), false))
	return nil
}
func evalDel(args []string, c io.ReadWriter) error {

	if len(args) == 0 {
		return errors.New("ERR wrong number of arguments for 'DEL' command")
	}
	count := 0
	for _, key := range args {
		if Del(key) {
			count += 1
		}
	}
	c.Write(Encode(count, false))
	return nil
}
func evalExpire(args []string, c io.ReadWriter) error {

	if len(args) != 2 {
		return errors.New("ERR wrong number of arguments for 'EXPIRE' command")
	}
	key := args[0]
	Obj := Get(key)
	if Obj == nil {
		c.Write(Encode(config.ZeroResponse, false))
	} else {
		expirationInSec, err := strconv.ParseInt(args[1], 10, 64) //converting decimal to 64 bit int
		if err != nil {
			return errors.New("(error) input is not an integer or out of range")
		}
		expirationInMs := expirationInSec * 1000
		Obj.ExpirationInMs = expirationInMs
	}
	c.Write(Encode(config.OneResponse, false))
	return nil
}
func EvaluateAndRespond(cmd *RespCmd, c io.ReadWriter) error {
	switch cmd.Cmd {
	case "PING":
		return evalCommand(cmd.Args, c)
	case "GET":
		return evalGet(cmd.Args, c)
	case "SET":
		return evalSet(cmd.Args, c)
	case "TTL":
		return evalTTL(cmd.Args, c)
	case "DEL":
		return evalDel(cmd.Args, c)
	case "EXPIRE":
		return evalExpire(cmd.Args, c)
	default:
		return evalCommand(cmd.Args, c)
	}
}
