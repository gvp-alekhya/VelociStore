// socket-server project main.go
package main

import (
    "fmt"
    "github.com/gvp-alekhya/VelociStore/server"
    "log"
    _ "net/http/pprof"
    "net/http"
)


// Main function where the server starts
func main() {
    fmt.Println("Server Running...")
    server.RunAsyncTCPServer()
    go func() {
        log.Println(http.ListenAndServe("localhost:2345", nil))
    }()
}