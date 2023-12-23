package server

import (
	"fmt"
	"net"
	"syscall"
	"github.com/gvp-alekhya/VelociStore/config"
	"github.com/gvp-alekhya/VelociStore/core"
)

var con_clients int = 0

func RunAsyncTCPServer() error {
	fmt.Println("starting an asynchronous TCP server on", config.SERVER_HOST, config.SERVER_PORT)

	// Create EPOLL Event Objects to hold events
	var events []syscall.EpollEvent = make([]syscall.EpollEvent, config.MAX_CLIENTS)
	serverFD,err := CreateAndBindSocket()

	defer syscall.Close(*serverFD)
	if err!=nil{
		fmt.Println(" Error in Binding Socket :: ", err)
	}
	// Start listening
	if err = syscall.Listen(*serverFD, config.MAX_CLIENTS); err != nil {
		return err
	}

		// creating EPOLL instance
	epollFD, err := syscall.EpollCreate1(0)
		if err != nil {
			fmt.Print(err)
	}
	defer syscall.Close(epollFD)
	// Adding server to Epoll to listen for events
	err = AddToEpoll(serverFD, epollFD)	
	if err != nil {
		return err
	}
	
	for {
		// see if any FD is ready for an IO
		nevents, e := syscall.EpollWait(epollFD, events[:], -1)
		if e != nil {
			fmt.Println("Error", e)
			return err
		}

		for i := 0; i < nevents; i++ {
			// if the socket server itself is ready for an IO
			if int(events[i].Fd) == *serverFD {
				// accept the incoming connection from a client
				fmt.Println("Server Ready for IO")
				//Accept connection
				fd, _, err := syscall.Accept(*serverFD)
				if err != nil {
					fmt.Println("err", err)
					continue
				}
				//Enable non block mode
				syscall.SetNonblock(*serverFD, true)
				//Add to ePoll for even notification
				err = AddToEpoll(&fd, epollFD)
				if err != nil {
					fmt.Println("Error adding cllient to epoll ", e)
				}
			} else {
				//Read from connection
				command := core.FDComm{Fd: int(events[i].Fd)}
				respCmd, err := ReadCommand(command)
				if err != nil {
					fmt.Println("Error reading from client ", e)
				}
				fmt.Println("Finished reading ")
				//Write to Connection
				WriteCommand(command, respCmd)
			}
		}
	}
}

func CreateAndBindSocket()(*int,error){
	// Create a socket
	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Println(" Error in Create Socket :: ", err)
		return &serverFD,err
	}

	// Set the Socket operate in a non-blocking mode
	if err = syscall.SetNonblock(serverFD, true); err != nil {
		fmt.Println(" Error in SetNonblock Socket :: ", err)
		return &serverFD,err
	}

	// Bind the IP and the port
	addr := syscall.SockaddrInet4{Port: config.SERVER_PORT}
	copy(addr.Addr[:], net.ParseIP(config.SERVER_HOST).To4())
	err = syscall.Bind(serverFD, &addr)
	if err != nil {
		fmt.Println("Error Binding to server:", err)
		return &serverFD,err
	}
	return &serverFD,err
}

func AddToEpoll(fd *int, epollFd int)(error){
	
		// Specify the events we want to get hints about
		// and set the socket on which
		var socketServerEvent syscall.EpollEvent = syscall.EpollEvent{
			Events: syscall.EPOLLIN,
			Fd:     int32(*fd),
		}
		var err error
		// Listen to read events on the Server itself
		if err = syscall.EpollCtl(epollFd, syscall.EPOLL_CTL_ADD, *fd, &socketServerEvent); err != nil {
			return err
		}
		return  err
}