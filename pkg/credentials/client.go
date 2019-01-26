package credentials

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
)

// ClientListener is the the type that handles is the callback for server responses
type ClientListener int

// CInput wait for the server question
func (l *ServerListener) CInput(in string) error {
	fmt.Println(in)
	end <- 0
	return nil
}

var end = make(chan int)

// SetupClient sets up the client
// This will be called if lazygit is called through git
func SetupClient() int {
	clientPort := os.Getenv("LAZYGIT_CLIENT_PORT")
	hostPort := os.Getenv("LAZYGIT_HOST_PORT")

	go func() {
		addy, err := net.ResolveTCPAddr("tcp", "127.0.0.1:"+clientPort)
		if err != nil {
			end <- 1
			return
		}

		inbound, err := net.ListenTCP("tcp", addy)
		if err != nil {
			end <- 1
			return
		}

		listener := new(ClientListener)

		rpc.Register(listener)
		rpc.Accept(inbound)

		end <- 0
	}()

	client, err := rpc.Dial("tcp", "127.0.0.1:"+hostPort)
	if err != nil {
		return 1
	}

	var reply bool
	err = client.Call("ServerListener.SInput", []byte(os.Args[1]), &reply)
	if err != nil {
		return 1
	}

	return <-end
}
