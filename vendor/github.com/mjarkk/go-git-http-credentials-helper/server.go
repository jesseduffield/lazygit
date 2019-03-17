package gitcredentialhelper

import (
	"errors"
	"net"
	"net/rpc"
	"os"
	"regexp"
	"strings"

	"github.com/mjarkk/go-ps"
)

// Listener represents all server functions
type Listener int

// InputQuestion is the data that is send to Host from the Client
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
	Options Options
}

// prompt the the meta data of of a request
type prompt struct {
	Pattern  string
	AskedFor *bool
}

var listenerMeta = map[string]listenerMetaType{} // a list of listeners
var totalListener uint32                         // this gets used to set the key of listenerMeta

// Input is the function that response on a credentials question from the client
func (l *Listener) Input(in InputQuestion, out *EncryptedMessage) error {
	suspiciousErr := errors.New("closing message due to suspicious behavior")

	listener, ok := listenerMeta[in.Listener]
	if !ok {
		return suspiciousErr
	}

	if !hasItselfAsSubProcess(listener.Options) {
		return suspiciousErr
	}

	updateListenerMeta := func() {
		listenerMeta[in.Listener] = listener
	}

	updateListenerMeta()

	question := in.Question

	prompts := map[string]prompt{
		"password": {Pattern: `Password\s*for\s*'.+':`, AskedFor: &listener.AskedFor.Password},
		"username": {Pattern: `Username\s*for\s*'.+':`, AskedFor: &listener.AskedFor.Username},
	}

	var toSend string

	for askFor, prompt := range prompts {
		match, _ := regexp.MatchString(prompt.Pattern, question)
		if match && !*prompt.AskedFor {
			*prompt.AskedFor = true
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

// runServer starts the server that waits for events from the client
func runServer(serverStartedChan chan struct{}, end chan endRun, hostPort, currentListener string) {
	serverRunning := false

	addy, err := net.ResolveTCPAddr("tcp", "127.0.0.1:"+hostPort)
	if err != nil {
		end <- endRun{
			Err: err,
			Out: []byte{},
		}
		return
	}

	inbound, err := net.ListenTCP("tcp", addy)
	if err != nil {
		end <- endRun{
			Err: err,
			Out: []byte{},
		}
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
		end <- endRun{
			Err: err,
			Out: []byte{},
		}
		return
	}

	serverStartedChan <- struct{}{}

	serverRunning = true
	rpc.Accept(inbound)
	serverRunning = false
}

// hasItselfAsSubProcess returns true if
// the current proccess has itself as sub process
func hasItselfAsSubProcess(options Options) bool {
	if !ps.Supported() {
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
		toMatch := os.Args[0]
		if len(options.AppName) > 0 {
			toMatch = options.AppName
		}
		if !strings.Contains(strings.Replace(toMatch, ".exe", "", -1), strings.Replace(procName, ".exe", "", -1)) {
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
			ex := proc.Executable()
			if !strings.Contains(ex, "git") && !strings.Contains(ex, "GIT") {
				continue procListLoop
			}
			parrent = proc.PPid()
		}
	}
	return false
}
