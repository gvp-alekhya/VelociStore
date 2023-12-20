// socket-server project main.go
package main

import (
    "fmt"
    "github.com/gvp-alekhya/VelociStore/server"
)


// Main function where the server starts
func main() {
    fmt.Println("Server Running...")
    tcp_conn.RunTCPServer()
}