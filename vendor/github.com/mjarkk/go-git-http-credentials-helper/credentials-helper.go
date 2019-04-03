package gitcredentialhelper

import (
	"errors"
	"fmt"
	"net/rpc"
	"os"
	"os/exec"
)

/*
Notes about the code comments
- When referred to the "Host" it's the aplication that calles spawns git
- When refrered to the "Client" it's the program created by git
*/

type endRun struct {
	Err error
	Out []byte
}

// Options are the argument options for run
type Options struct {
	AppName string
}

// Run runs git commands that need credentials
// ask() gets executed when git needs credentials
// The ask argument will be "username" or "password"
func Run(cmd *exec.Cmd, ask func(string) string, potentialOptions ...Options) ([]byte, error) {
	options := Options{}

	if len(potentialOptions) == 1 {
		options = potentialOptions[0]
	}

	if len(potentialOptions) > 1 {
		return nil, errors.New("The options aregument can only be 1 or 0 arguments")
	}

	totalListener++
	currentListener := fmt.Sprintf("%v", totalListener)

	listener := listenerMetaType{
		AskFunction: ask,
		Options:     options,
	}
	listenerMeta[currentListener] = listener

	defer delete(listenerMeta, currentListener)

	end := make(chan endRun)
	hostPort := getFreePort()
	serverStartedChan := make(chan struct{})

	go func() {
		<-serverStartedChan

		ex, err := os.Executable()
		if err != nil {
			ex = os.Args[0]
		}

		if len(cmd.Env) == 0 {
			cmd.Env = os.Environ()
		}

		cmd.Env = append(
			cmd.Env,
			"GIT_CREDENTIALS_HOST_PORT="+hostPort,
			"GIT_CREDENTIALS_LISTENER="+currentListener,

			"GIT_ASKPASS="+ex,
			"LANG=en_US.UTF-8",
			"LC_ALL=en_US.UTF-8",
		)
		out, err := cmd.CombinedOutput()
		end <- endRun{
			Err: err,
			Out: out,
		}
	}()

	go runServer(
		serverStartedChan,
		end,
		hostPort,
		currentListener,
	)

	endData := <-end
	return endData.Out, endData.Err
}

// SetupClient sets up the client
// BEFORE THIS FUNCTION RUNS THERE CAN'T BE ANY PRINTING OTHERWHISE THIS LIBARY WON'T WORK!
func SetupClient(logger ...func(error)) {
	log := func(err error) {}
	if len(logger) > 0 {
		log = logger[0]
	}

	port := os.Getenv("GIT_CREDENTIALS_HOST_PORT")
	listener := os.Getenv("GIT_CREDENTIALS_LISTENER")

	if len(port) == 0 || len(listener) == 0 {
		return
	}

	privateKey, publicKey, err := generateKeyPair()
	if err != nil {
		log(err)
		os.Exit(0)
		return
	}

	client, err := rpc.Dial("tcp", "127.0.0.1:"+port)
	if err != nil {
		log(err)
		os.Exit(0)
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
		log(err)
		os.Exit(0)
		return
	}

	msg, err := decryptMessage(privateKey, *data)
	if err != nil {
		log(err)
		os.Exit(0)
		return
	}

	fmt.Println(msg)
	os.Exit(0)
}
