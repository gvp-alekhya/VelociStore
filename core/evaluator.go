package core

import (
	"errors"
	"fmt"
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
	fmt.Println("Encoded String :", string(b))
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
	fmt.Println("Key :: value", key, value)

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
		c.Write([]byte(config.NILResponse))
		return nil
	}
	c.Write(Encode(Obj.Value, true))

	return nil
}
func evalTTL(args []string, c io.ReadWriter) error {
	fmt.Print("Into evalTTL", args)
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

	var timeLeft = currentTimeInMS - Obj.ExpirationInMs
	if timeLeft < 0 {
		c.Write([]byte(config.NoKeyResponse))
		return errors.New("key expired")
	}

	c.Write(Encode(int64(timeLeft/1000), true))
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
	default:
		return evalCommand(cmd.Args, c)
	}
}
