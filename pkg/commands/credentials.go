package commands

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"regexp"

	"github.com/mgutz/str"
)

// ServerListener is the the type that handles is the callback for server responses
type ServerListener int

// Input wait for the server question
func (l *ServerListener) Input(in string, out *string) error {
	hasMessage <- in
	toReturn := <-resMessage
	out = &toReturn
	return nil
}

var hasMessage = make(chan string)
var resMessage = make(chan string)
var hostPort = ""
var clientPort = ""

// DetectUnamePass detect a username / password question in a command
// ask is a function that gets executen when this function detect you need to fillin a password
// The ask argument will be "username" or "password" and expects the user's password or username back
func (c *OSCommand) DetectUnamePass(command string, ask func(string) string) error {
	end := make(chan error)

	hostPort := GetRandomFreePorts(1)[0]

	go func() {
		for true {
			question := <-hasMessage
			prompts := map[string]string{
				"password": `Password\s*for\s*'.+':`,
				"username": `Username\s*for\s*'.+':`,
			}

			hasResponded := false
		toCloseLoop:
			for askFor, pattern := range prompts {
				if match, _ := regexp.MatchString(pattern, question); match {
					resMessage <- ask(askFor)
					hasResponded = true
					break toCloseLoop
				}
			}
			if !hasResponded {
				resMessage <- ""
			}
		}
	}()

	go func() {
		ex, err := os.Executable() // get the executable path for git to use
		if err != nil {
			ex = os.Args[0] // fallback to the first call argument if needed
		}

		splitCmd := str.ToArgv(command)
		cmd := exec.Command(splitCmd[0], splitCmd[1:]...)
		cmd.Env = os.Environ()
		cmd.Env = append(
			cmd.Env,
			"LAZYGIT_ASK_FOR_PASS=true",   // tell the sub lazygit process that this ran from git
			"LAZYGIT_HOST_PORT="+hostPort, // The main process communication port
			"GIT_ASKPASS="+ex,             // tell git where lazygit is located,
			"LANG=en_US.UTF-8",            // Force using EN as language
			"LC_ALL=en_US.UTF-8",          // Force using EN as language
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			outString := string(out)
			if len(outString) == 0 {
				end <- err
				return
			}
			end <- errors.New(outString)
			return
		}
		end <- nil
	}()

	go func() {
		serverErr := startServer()
		if serverErr != nil {
			end <- serverErr
		}
	}()

	return <-end
}

// serverStarted tells if the the server is started yet
var serverStarted = false

// startServer starts the credentials server if it's not turned on yet
func startServer() error {
	if !serverStarted {
		defer func() {
			// if this server stops for a reason make sure it's turned on again
			serverStarted = false
		}()
		serverStarted = true
		addy, err := net.ResolveTCPAddr("tcp", "127.0.0.1:"+hostPort)
		if err != nil {
			return err
		}

		inbound, err := net.ListenTCP("tcp", addy)
		if err != nil {
			return err
		}

		listener := new(ServerListener)
		rpc.Register(listener)
		rpc.Accept(inbound)

		return nil
	}
	return nil
}

// GetRandomFreePorts returns a n number of free ports free port
func GetRandomFreePorts(amound int) []string {
	checkFrom := 5000
	toReturn := []string{}
	for amound != 0 {
		checkFrom++
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

// SetupClient sets up the client
// This will be called if lazygit is called through git
func SetupClient() int {
	hostPort := os.Getenv("LAZYGIT_HOST_PORT")

	client, err := rpc.Dial("tcp", "127.0.0.1:"+hostPort)
	if err != nil {
		return 1
	}

	var reply string
	err = client.Call("ServerListener.Input", []byte(os.Args[1]), &reply)
	if err != nil {
		return 1
	}

	fmt.Println(reply)

	return 0
}
