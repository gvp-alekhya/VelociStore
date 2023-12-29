package server

import (
	"fmt"
	"io"
	"strings"

	"github.com/gvp-alekhya/VelociStore/core"
)

// Function to process a client's connection
func ReadCommands(connection io.ReadWriter) (core.RespCmds, error) {
	// Create a buffer to read data from the client
	buffer := make([]byte, 1024)

	// Read data from the client into the buffer
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	command := (buffer[:mLen])
	commands, err := core.DecodeCommands(command)
	fmt.Println("ReadCommands DecodeCommands:", commands)
	if err != nil {
		return nil, err
	}
	redisCmds := make([]*core.RespCmd, 0)
	for _, value := range commands {
		tokens, err := toArrayString(value.([]interface{}))
		if err != nil {
			return nil, err
		}
		respcmd := core.RespCmd{
			Cmd:  strings.ToUpper(tokens[0]),
			Args: tokens[1:],
		}
		redisCmds = append(redisCmds, &respcmd)
	}

	return redisCmds, nil
}

func WriteCommand(connection io.ReadWriter, respcmds core.RespCmds) {
	// Send a response to the client acknowledging receipt of the message
	core.EvaluateAndRespond(respcmds, connection)

}
func toArrayString(ai []interface{}) ([]string, error) {
	as := make([]string, len(ai))
	for i := range ai {
		as[i] = ai[i].(string)
	}
	return as, nil
}
