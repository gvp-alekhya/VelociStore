package server

import (
	"io"
	"fmt"
	"strings"
	"github.com/gvp-alekhya/VelociStore/core"
)
// Function to process a client's connection
func ReadCommand(connection io.ReadWriter) (*core.RespCmd, error) {
    // Create a buffer to read data from the client
    buffer := make([]byte, 1024)

    // Read data from the client into the buffer
    mLen, err := connection.Read(buffer)
    if err != nil {
        fmt.Println("Error reading:", err.Error())
    }
	command := (buffer[:mLen])
    // Print the received data from the client
    tokens,err := core.DecodeArrayString(command)
    if err != nil {
		return nil, err
	}
	fmt.Println("strings.ToUpper(tokens[0]) :: ", strings.ToUpper(tokens[0]))
    respcmd := core.RespCmd{
        Cmd:  strings.ToUpper(tokens[0]),
        Args: tokens[1:],
    }
	return &respcmd, nil
}

func WriteCommand (connection io.ReadWriter, respcmd *core.RespCmd)  ( err error){
	 // Send a response to the client acknowledging receipt of the message
     err = core.EvaluateAndRespond(respcmd, connection)
	 if err!=nil{
        errMsg := fmt.Sprintf("-%s\r\n", err)
        connection.Write([]byte(errMsg))
     }
	 return
}
