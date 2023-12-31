package core

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/gvp-alekhya/VelociStore/config"
)

func evalCommand(args []string, c io.ReadWriter) []byte {
	var b []byte

	if len(args) >= 2 {
		return Encode("(error) ERR wrong number of arguments for 'ping' command", false)
	}

	if len(args) == 0 {
		b = Encode("PONG", false)
	} else {
		b = Encode(args[0], false)
	}
	return b
}
func evalSet(args []string, c io.ReadWriter) []byte {
	if len(args) <= 1 {
		return Encode("(error) ERR wrong number of arguments for 'SET' command", false)
	}

	var key, value string
	var expirationInMs int64 = -1

	key = args[0]
	value = args[1]
	objType, objEncoding := DeduceTypeEncoding(value)
	for i := 2; i < len(args); i++ {
		lowercaseArgs := strings.ToLower(args[i])
		switch lowercaseArgs {
		case "ex":
			if i == len(args) {
				return Encode("(error) ERR syntax error", false)
			}
			i++
			expirationInSec, err := strconv.ParseInt(args[i], 10, 64) //converting decimal to 64 bit int
			if err != nil {
				return Encode("(error) ERR invalid expiration value passed in 'SET' command", false)
			}
			expirationInMs = expirationInSec * 1000

		default:
			return Encode("(error) ERR invalid arguments for 'SET' command", false)
		}
	}
	Put(key, NewObj(value, expirationInMs, objType, objEncoding))
	return Encode(config.OKResponse, false)
}
func evalGet(args []string, c io.ReadWriter) []byte {
	if len(args) != 1 {
		return Encode("(error) ERR wrong number of arguments for 'GET' command", false)
	}

	key := args[0]
	Obj := Get(key)
	if Obj == nil {
		return []byte(config.NILResponse)
	}
	if Obj.ExpirationInMs != -1 && Obj.ExpirationInMs <= time.Now().UnixMilli() {
		Del(key)
		return []byte(config.NILResponse)
	}
	return Encode(Obj.Value, true)
}
func evalTTL(args []string, c io.ReadWriter) []byte {

	if len(args) != 1 {
		return Encode("(error) ERR wrong number of arguments for 'TTL' command", false)
	}
	key := args[0]

	Obj := Get(key)
	if Obj == nil {

		return []byte(config.NoKeyResponse)
	}
	currentTimeInMS := time.Now().UnixMilli()
	if Obj.ExpirationInMs == -1 {

		return []byte(config.NoExpiryResponse)
	}

	var timeLeft = Obj.ExpirationInMs - currentTimeInMS
	if timeLeft < 0 {

		return Encode("key expired", false)
	}
	return (Encode(int64(timeLeft/1000), false))
}
func evalDel(args []string, c io.ReadWriter) []byte {

	if len(args) == 0 {
		return Encode("(error) ERR wrong number of arguments for 'DEL' command", false)
	}
	count := 0
	for _, key := range args {
		if Del(key) {
			count += 1
		}
	}
	return Encode(count, false)
}
func evalExpire(args []string, c io.ReadWriter) []byte {

	if len(args) != 2 {
		return Encode("(error) ERR wrong number of arguments for 'EXPIRE' command", false)
	}
	key := args[0]
	Obj := Get(key)
	if Obj == nil {
		return (Encode(config.ZeroResponse, false))
	} else {
		expirationInSec, err := strconv.ParseInt(args[1], 10, 64) //converting decimal to 64 bit int
		if err != nil {
			return Encode("(error) input is not an integer or out of range", false)
		}
		expirationInMs := expirationInSec * 1000
		Obj.ExpirationInMs = expirationInMs
	}
	return Encode(config.OneResponse, false)
}
func evalRewriteAOF(args []string) []byte {
	GenerateDumpAOF()
	return []byte(config.OKResponse)
}
func evalIncr(args []string, c io.ReadWriter) []byte {
	if len(args) != 1 {
		return Encode("(error) ERR wrong number of arguments for 'INCR' command", false)
	}

	var key string = args[0]
	obj := Get(key)

	if obj == nil {
		obj = NewObj("0", -1, OBJ_TYPE_STRING, OBJ_ENCODING_INT)
		Put(key, obj)
	}

	if err := assertType(obj.TypeEncoding, OBJ_TYPE_STRING); err != nil {
		return Encode("(error) ERR Type is not an integer", false)
	}

	if err := assertEncoding(obj.TypeEncoding, OBJ_ENCODING_INT); err != nil {
		return Encode("(error) ERR Encoding is not an integer", false)
	}

	// Parse the integer value
	parsedValue, err := strconv.ParseInt(obj.Value.(string), 10, 64)
	if err != nil {
		return Encode("(error) ERR value is not an integer", false)
	}
	parsedValue++
	obj.Value = strconv.FormatInt(parsedValue, 10)
	// Convert the result back to string and encode
	return Encode(parsedValue, false)

}
func evalInfo() []byte {
	var response []byte
	buf := bytes.NewBuffer(response)
	buf.WriteString("# Keyspace\r\n")
	for i := range KeySpaceStats {
		buf.WriteString(fmt.Sprintf("db%d:keys=%d,expires=0,avg_ttl=0\r\n", i, KeySpaceStats[i]["keys"]))
	}
	return Encode(buf.String(), false)
}

func EvaluateAndRespond(cmds RespCmds, c io.ReadWriter) {
	var response []byte
	buf := bytes.NewBuffer(response)
	for _, cmd := range cmds {
		fmt.Println("EvaluateAndRespond :: cmd", cmd)
		switch cmd.Cmd {
		case "PING":
			buf.Write(evalCommand(cmd.Args, c))
		case "GET":
			buf.Write(evalGet(cmd.Args, c))
		case "SET":
			buf.Write(evalSet(cmd.Args, c))
		case "TTL":
			buf.Write(evalTTL(cmd.Args, c))
		case "DEL":
			buf.Write(evalDel(cmd.Args, c))
		case "EXPIRE":
			buf.Write(evalExpire(cmd.Args, c))
		case "INCR":
			buf.Write(evalIncr(cmd.Args, c))
		case "BGREWRITEAOF":
			buf.Write(evalRewriteAOF(cmd.Args))
		case "INFO":
			buf.Write(evalInfo())
		default:
			buf.Write(evalCommand(cmd.Args, c))
		}
	}
	c.Write(buf.Bytes())
}
