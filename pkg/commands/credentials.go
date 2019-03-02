package commands

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/mgutz/str"
	"github.com/mjarkk/go-ps"
)

// Listener is the the type that handles is the callback for server responses
type Listener int

// InputQuestion is what is send by the lazygit client
type InputQuestion struct {
	PublicKey string
	Question  string
	Listener  string
}

// listenerMetaType is a listener there private key and ask function
type listenerMetaType struct {
	AskFunction func(string) string
	AskedFor    struct {
		Password bool
		Username bool
	}
}

var listenerMeta = map[string]listenerMetaType{} // a list of listeners
var totalListener uint32                         // this gets used to set the key of listenerMeta

// Input interacts with the lazygit client spawned by git
func (l *Listener) Input(in InputQuestion, out *EncryptedMessage) error {
	suspiciousErr := errors.New("closing message due to suspicious behavior")

	listener, ok := listenerMeta[in.Listener]
	if !ok {
		return suspiciousErr
	}

	if !HasLGAsSubProcess() {
		return suspiciousErr
	}

	updateListenerMeta := func() {
		listenerMeta[in.Listener] = listener
	}

	updateListenerMeta()

	question := in.Question

	prompts := map[string]string{
		"password": `Password\s*for\s*'.+':`,
		"username": `Username\s*for\s*'.+':`,
	}

	var toSend string

	for askFor, pattern := range prompts {
		match, _ := regexp.MatchString(pattern, question)
		if match && ((askFor == "password" && !listener.AskedFor.Password) || (askFor == "username" && !listener.AskedFor.Username)) {
			switch askFor {
			case "password":
				listener.AskedFor.Password = true
			case "username":
				listener.AskedFor.Username = true
			}
			updateListenerMeta()
			toSend = strings.Replace(listener.AskFunction(askFor), "\n", "", -1)
			break
		}
	}

	encryptedData, err := encryptMessage(in.PublicKey, toSend)
	if err != nil {
		return suspiciousErr
	}

	*out = encryptedData

	return nil
}

// DetectUnamePass runs git commands that need credentials
// ask() gets executed when git needs credentials
// The ask argument will be "username" or "password"
func (c *OSCommand) DetectUnamePass(command string, ask func(string) string) error {
	totalListener++
	currentListener := fmt.Sprintf("%v", totalListener)

	listener := listenerMetaType{AskFunction: ask}
	listenerMeta[currentListener] = listener

	defer delete(listenerMeta, currentListener)

	end := make(chan error)
	hostPort := GetFreePort()
	serverStartedChan := make(chan struct{})

	go runGit(
		serverStartedChan,
		end,
		command,
		hostPort,
		currentListener,
	)

	go runServer(
		serverStartedChan,
		end,
		hostPort,
		currentListener,
	)

	err := <-end

	return err
}

// runGit runs the actual git command with the needed git
func runGit(serverStartedChan chan struct{}, end chan error, command, hostPort, currentListener string) {
	<-serverStartedChan

	ex, err := os.Executable()
	if err != nil {
		ex = os.Args[0]
	}

	splitCmd := str.ToArgv(command)
	cmd := exec.Command(splitCmd[0], splitCmd[1:]...)
	cmd.Env = os.Environ()
	cmd.Env = append(
		cmd.Env,
		"LAZYGIT_ASK_FOR_PASS=true",
		"LAZYGIT_HOST_PORT="+hostPort,
		"LAZYGIT_LISTENER="+currentListener, // the lisener ID

		"GIT_ASKPASS="+ex,    // tell git where lazygit is located so it can ask lazygit for credentials
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
}

// runServer starts the server that waits for events from the lazygit client
func runServer(serverStartedChan chan struct{}, end chan error, hostPort, currentListener string) {
	serverRunning := false

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

	go func() {
		<-end
		if serverRunning {
			inbound.Close()
		}
	}()

	listener := new(Listener)

	// every listener needs a different name it this is not dune rpc.RegisterName will error
	err = rpc.RegisterName("Listener"+currentListener, listener)
	if err != nil {
		end <- err
		return
	}

	serverStartedChan <- struct{}{}

	serverRunning = true
	rpc.Accept(inbound)
	serverRunning = false
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
	port := os.Getenv("LAZYGIT_HOST_PORT")
	listener := os.Getenv("LAZYGIT_LISTENER")

	privateKey, publicKey, err := generateKeyPair()
	if err != nil {
		return
	}

	client, err := rpc.Dial("tcp", "127.0.0.1:"+port)
	if err != nil {
		return
	}

	var data *EncryptedMessage
	err = client.Call("Listener"+listener+".Input", InputQuestion{
		Question:  os.Args[len(os.Args)-1],
		Listener:  listener,
		PublicKey: publicKey,
	}, &data)
	client.Close()
	if err != nil {
		return
	}

	msg, err := decryptMessage(privateKey, *data)
	if err != nil {
		return
	}

	fmt.Println(msg)
}

// HasLGAsSubProcess returns true if lazygit is a child of this process
func HasLGAsSubProcess() bool {
	if !ps.Supported() {
		return true
	}

	if runtime.GOOS == "windows" {
		return true
	}

	lgHostPid := os.Getpid()
	list, err := ps.Processes()
	if err != nil {
		return false
	}
procListLoop:
	for _, proc := range list {
		procName := proc.Executable()
		if procName != "lazygit" {
			continue
		}
		parrent := proc.PPid()
		for {
			if parrent < 30 {
				continue procListLoop
			}
			proc, err := ps.FindProcess(parrent)
			if err != nil {
				continue procListLoop
			}
			if proc.Pid() == lgHostPid {
				return true
			}
			parrent = proc.PPid()
		}
	}
	return false
}
