package server

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"

	"github.com/gvp-alekhya/VelociStore/config"
)

func RunSyncTCPServer() {

	// Attempt to listen on the specified network address and port
	server, err := net.Listen(config.SERVER_TYPE, config.SERVER_HOST+":"+strconv.Itoa(config.SERVER_PORT))
	if err != nil {
		// Print an error message and exit the program if listening fails
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	// Defer closing the server connection until the main function exits
	defer server.Close()

	concurrentClient := 0
	fmt.Println("Listening on " + config.SERVER_HOST + ":" + strconv.Itoa(config.SERVER_PORT))
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
		//	go func(connection io.ReadWriter, clientAddr net.Addr){
		for {
			command, err := ReadCommands(connection)
			if err != nil {
				// Print an error message and exit the program if accepting fails
				connection.Close()
				concurrentClient--
				fmt.Println("Client disconnected: ", connection.RemoteAddr(), "concurrent clients :: ", concurrentClient)
				if err == io.EOF {
					break
				}
			}
			fmt.Printf("command %v \n:: ", command)
			WriteCommand(connection, command)

		}
		//}(connection, connection.RemoteAddr())
	}
}
