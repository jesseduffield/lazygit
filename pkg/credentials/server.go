package credentials

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"sync"

	"github.com/mgutz/str"
)

type ServerListener int

func (l *ServerListener) Input() error {

	return nil
}

var hostPort = ""
var clientPort = ""

// SetupServer sets up the server
// This creates a server where the client can connect to
func SetupServer(command string, ask func(string) string) error {

	freePortsArr := GetRandomFreePorts(2)
	hostPort = freePortsArr[0]
	clientPort = freePortsArr[1]

	var waitForServ sync.WaitGroup
	waitForServ.Add(1)

	end := make(chan error)

	go func() {
		defer waitForServ.Done()

		ex, err := os.Executable()
		if err != nil {
			ex = os.Args[0] // fallback to the first call argument if needed
		}

		splitCmd := str.ToArgv(command)
		cmd := exec.Command(splitCmd[0], splitCmd[1:]...)
		cmd.Env = os.Environ()
		cmd.Env = append(
			cmd.Env,
			"LAZYGIT_ASK_FOR_PASS=true", // tell the sub lazygit process that this ran from git
			"LAZYGIT_HOST_PORT="+hostPort,
			"LAZYGIT_CLIENT_PORT="+clientPort,
			"GIT_ASKPASS="+ex, // tell git where lazygit is located
		)
		_, err = cmd.Output()
		end <- err
	}()

	go func() {
		addy, err := net.ResolveTCPAddr("tcp", "127.0.0.1:"+hostPort)
		if err != nil {
			end <- err
			return
		}

		inbound, err := net.ListenTCP("tcp", addy)
		if err != nil {
			end <- err
			return
		}

		listener := new(ServerListener)

		go func() {
			waitForServ.Done()
			// close the server
		}()
	}()

	return <-end
}

// GetRandomFreePorts returns a n number of free ports free port
func GetRandomFreePorts(amound int) []string {
	checkFrom := 5000
	toReturn := []string{}
	for amound != 0 {
		checkFrom++
		fmt.Println(checkFrom, amound)
		toCheck := fmt.Sprintf("%v", checkFrom)
		if IsFreePort(toCheck) {
			amound = amound - 1
			toReturn = append(toReturn, toCheck)
		}
	}
	return toReturn
}

// IsFreePort return true if the port if not in use
func IsFreePort(port string) bool {
	conn, err := net.Dial("tcp", "127.0.0.1:"+port)
	if err == nil {
		go conn.Close()
		return false
	}
	return true
}
