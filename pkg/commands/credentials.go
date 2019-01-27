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

var totalListener uint32
var askFunc func(string) string

// DetectUnamePass detect a username / password question in a command
// ask is a function that gets executen when this function detect you need to fillin a password
// The ask argument will be "username" or "password" and expects the user's password or username back
func (c *OSCommand) DetectUnamePass(command string, ask func(string) string) error {
	totalListener++
	currentListener := fmt.Sprintf("%v", totalListener)
	askFunc = ask
	end := make(chan error)
	hostPort := GetFreePort()
	serverRunning := false
	serverStartedChan := make(chan struct{})
	var inbound *net.TCPListener

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
			"LAZYGIT_LISTENER="+currentListener,
			"GIT_ASKPASS="+ex,    // tell git where lazygit is located,
			"LANG=en_US.UTF-8",   // Force using EN as language
			"LC_ALL=en_US.UTF-8", // Force using EN as language
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
		addy, err := net.ResolveTCPAddr("tcp", "127.0.0.1:"+hostPort)
		if err != nil {
			end <- err
			return
		}

		in, err := net.ListenTCP("tcp", addy)
		inbound = in
		if err != nil {
			end <- err
			return
		}

		listener := new(Listener)

		// every listener needs a different name it this is not dune rpc.RegisterName will error
		err = rpc.RegisterName("Listener"+currentListener, listener)
		if err != nil {
			end <- err
			return
		}

		serverStartedChan <- struct{}{}
		rpc.Accept(inbound)

		serverRunning = false
	}()

	err := <-end
	if serverRunning {
		inbound.Close()
	}
	askFunc = func(i string) string { return "" } // make sure that the program doesn't popup a input for credentials if not needed

	return err
}

// GetFreePort returns a free port that can be used by lazygit
func GetFreePort() string {
	checkFrom := 5000
	toReturn := ""
	for {
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
	ListenerNumber := os.Getenv("LAZYGIT_LISTENER")

	client, err := rpc.Dial("tcp", "127.0.0.1:"+hostPort)
	if err != nil {
		return
	}

	var rply *string
	err = client.Call("Listener"+ListenerNumber+".Input", os.Args[1], &rply)
	if err != nil {
		return
	}

	fmt.Println(*rply)
}
