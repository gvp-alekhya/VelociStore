package server

import ("fmt"
"io"
"net"
"os"
"strings"
"github.com/gvp-alekhya/VelociStore/core")
// Constants defining server details
const (
    SERVER_HOST = "0.0.0.0"
    SERVER_PORT = "2929"
    SERVER_TYPE = "tcp"
)
func RunTCPServer(){

    // Attempt to listen on the specified network address and port
    server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
    if err != nil {
        // Print an error message and exit the program if listening fails
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }

    // Defer closing the server connection until the main function exits
    defer server.Close()

    concurrentClient := 0;
    fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)
    fmt.Println("Waiting for client...")

    // Loop to continuously accept incoming client connections
    for {
        // Accept a connection from a client
        connection, err := server.Accept()
		concurrentClient += 1
        if err != nil {
            // Print an error message and exit the program if accepting fails
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }

        // Print a message indicating that a client has connected
        fmt.Println("client connected")
		// add go routine for concurrent requests - not supported
  //	go func(connection net.Conn, clientAddr net.Addr){
		for{  
			command, err := readCommand(connection)
			if err != nil {
				// Print an error message and exit the program if accepting fails
				connection.Close();
				concurrentClient--;
				fmt.Println("Client disconnected: ", connection.RemoteAddr(), "concurrent clients :: " , concurrentClient)
				if err == io.EOF {
					break
				}
			}
            fmt.Printf("command %v \n:: ", command)
			err = writeCommand(connection, command)
			if err != nil {
				fmt.Printf("Write failed %v \n:: ", err)
			}
		}
	//}(connection, connection.RemoteAddr())
  }
}

// Function to process a client's connection
func readCommand(connection net.Conn) (*core.RespCmd, error) {
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
    respcmd := &core.RespCmd{
        Cmd:  strings.ToUpper(tokens[0]),
        Args: tokens[1:],
    }
	return respcmd, nil
}

func writeCommand (connection net.Conn, respcmd *core.RespCmd)  ( err error){
	 // Send a response to the client acknowledging receipt of the message
     err = core.EvaluateAndRespond(respcmd, connection)
	 if err!=nil{
        errMsg := fmt.Sprintf("-%s\r\n", err)
        connection.Write([]byte(errMsg))
     }
	 return
}
