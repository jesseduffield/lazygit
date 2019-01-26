package commands

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/mgutz/str"
)

// Listener is the the type that handles is the callback for server responses
type Listener int

// Input wait for the server question
func (l *Listener) Input(in string, out *string) error {
	prompts := map[string]string{
		"password": `Password\s*for\s*'.+':`,
		"username": `Username\s*for\s*'.+':`,
	}

	for askFor, pattern := range prompts {
		if match, _ := regexp.MatchString(pattern, in); match {
			*out = strings.Replace(askFunc(askFor), "\n", "", -1)
			break
		}
	}

	return nil
}

var askFunc func(string) string
var hostPort = ""
var clientPort = ""

// DetectUnamePass detect a username / password question in a command
// ask is a function that gets executen when this function detect you need to fillin a password
// The ask argument will be "username" or "password" and expects the user's password or username back
func (c *OSCommand) DetectUnamePass(command string, ask func(string) string) error {
	askFunc = ask

	end := make(chan error)

	if len(hostPort) == 0 {
		hostPort = GetFreePort()
	}

	serverStartedChan := make(chan struct{})
	go func() {
		<-serverStartedChan

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
		serverErr := startServer(serverStartedChan)
		if serverErr != nil {
			end <- serverErr
		}
	}()

	return <-end
}

// serverStarted tells if the the server is started yet
var serverStarted = false

// startServer starts the credentials server if it's not turned on yet
func startServer(started chan struct{}) error {
	if !serverStarted {
		defer func() {
			// if this server stops for a reason this make sure it can be turned on again
			serverStarted = false
		}()
		serverStarted = true
		addy, err := net.ResolveTCPAddr("tcp", ":"+hostPort)
		if err != nil {
			return err
		}

		inbound, err := net.ListenTCP("tcp", addy)
		if err != nil {
			return err
		}

		listener := new(Listener)
		rpc.Register(listener)
		rpc.Accept(inbound)

		return nil
	}
	started <- struct{}{}
	return nil
}

// GetFreePort returns a free port that can be used by lazygit
func GetFreePort() string {
	checkFrom := 5000
	toReturn := ""
	for true {
		checkFrom++
		check := fmt.Sprintf("%v", checkFrom)
		if IsFreePort(check) {
			toReturn = check
			break
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
func SetupClient() {
	hostPort := os.Getenv("LAZYGIT_HOST_PORT")

	client, err := rpc.Dial("tcp", "127.0.0.1:"+hostPort)
	if err != nil {
		return
	}

	var rply *string
	err = client.Call("Listener.Input", os.Args[1], &rply)
	if err != nil {
		return
	}

	time.Sleep(time.Millisecond * 50)

	fmt.Println(*rply)
}
