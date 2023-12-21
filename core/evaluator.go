package core

import (
	"errors"
	"fmt"
	"net"
)

func evalCommand(args []string, c net.Conn) error {
	var b []byte

	if len(args) >= 2 {
		return errors.New("ERR wrong number of arguments for 'ping' command")
	}

	if len(args) == 0 {
		b = EncodeString("PONG")
	} else {
		b = EncodeBinaryString(args[0])
	}
	fmt.Println("Encoded String :", string(b))
	_, err := c.Write(b)
	return err
}

func EvaluateAndRespond(cmd *RespCmd, c net.Conn) error {
	fmt.Println("EvaluateAndRespond comand::", cmd.Cmd)
	switch cmd.Cmd {
	case "PING":
		return evalCommand(cmd.Args, c)
	default:
		return evalCommand(cmd.Args, c)
	}
}