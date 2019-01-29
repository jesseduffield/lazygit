package commands

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
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

// InputQuestion is what is send by the lazygit client
type InputQuestion struct {
	ClientPublicKeyText string // the client there public key
	Message             []byte // this is encrypted using the public key of the host
	ListenerKey         string // the listener key
}

// listenerMetaType is a listener there private key and ask function
type listenerMetaType struct {
	HostPrivateKey *rsa.PrivateKey
	AskFunction    func(string) string
}

// listenerMetaType is a list of listeners
var listenerMeta = map[string]listenerMetaType{}

// totalListener is the current amound of listeners
// this is used to track the amound listeners
// note: when this number runs out of numbers the application stops working
var totalListener uint32

// Listener is the the type that handles is the callback for server responses
type Listener int

// Input wait for the server question
func (l *Listener) Input(fromClient InputQuestion, out *[]byte) error {
	listener, ok := listenerMeta[fromClient.ListenerKey]
	if !ok {
		return errors.New("Listener not defined")
	}

	message := fromClient.Message

	clientPubBlock, _ := pem.Decode([]byte(fromClient.ClientPublicKeyText))
	clientPub, err := x509.ParsePKCS1PublicKey(clientPubBlock.Bytes)
	if err != nil {
		return errors.New("Can't parse public key")
	}

	// TODO send errors to DetectUnamePass
	decryptedMessageRaw, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, listener.HostPrivateKey, message, []byte("TELL HOST"))
	if err != nil {
		return errors.New("Decryption failed")
	}

	decryptedMessage := strings.Split(string(decryptedMessageRaw), "|")
	if len(decryptedMessage) != 2 {
		return errors.New("Missing string parts")
	}
	validation := decryptedMessage[0]
	question := decryptedMessage[1]

	if fmt.Sprintf("%x", sha256.Sum256([]byte(fromClient.ClientPublicKeyText))) != validation {
		return errors.New("Mismatch hash of public key")
	}

	prompts := map[string]string{
		"password": `Password\s*for\s*'.+':`,
		"username": `Username\s*for\s*'.+':`,
	}

	toSend := ""
	for askFor, pattern := range prompts {
		if match, _ := regexp.MatchString(pattern, question); match {
			toSend = strings.Replace(listener.AskFunction(askFor), "\n", "", -1)
			break
		}
	}

	encrpytedToSend, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, clientPub, []byte(toSend), []byte("TO PRINT"))
	if err != nil {
		return errors.New("Can't encrpyt string")
	}
	*out = encrpytedToSend

	return nil
}

// DetectUnamePass detect a username / password question in a command
// ask is a function that gets executen when this function detect you need to fillin a password
// The ask argument will be "username" or "password" and expects the user's password or username back
func (c *OSCommand) DetectUnamePass(command string, ask func(string) string) error {
	totalListener++
	currentListener := fmt.Sprintf("%v", totalListener)

	hostPriv, err := rsa.GenerateKey(rand.Reader, 3072)
	if err != nil {
		return err
	}

	listener := listenerMetaType{
		AskFunction:    ask,
		HostPrivateKey: hostPriv,
	}
	listenerMeta[currentListener] = listener
	defer delete(listenerMeta, currentListener)

	pubKeyText := string(pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(&hostPriv.PublicKey)}))
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
			"LAZYGIT_ASK_FOR_PASS=true",           // tell the sub lazygit process that this ran from git
			"LAZYGIT_HOST_PORT="+hostPort,         // The main process communication port
			"LAZYGIT_HOST_PUBLIC_KEY="+pubKeyText, // the public key of the host
			"LAZYGIT_LISTENER="+currentListener,   // the lisener ID

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

	err = <-end
	if serverRunning {
		inbound.Close()
	}

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
	hostPubText := os.Getenv("LAZYGIT_HOST_PUBLIC_KEY")
	hostPort := os.Getenv("LAZYGIT_HOST_PORT")
	listenerNumber := os.Getenv("LAZYGIT_LISTENER")

	hostPubBlock, _ := pem.Decode([]byte(hostPubText))
	hostPub, err := x509.ParsePKCS1PublicKey(hostPubBlock.Bytes)
	if err != nil {
		return
	}

	clientPriv, err := rsa.GenerateKey(rand.Reader, 3072)
	if err != nil {
		return
	}
	clientPub := clientPriv.PublicKey
	clientPubText := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(&clientPub)})
	toSend := fmt.Sprintf("%x|%v", sha256.Sum256(clientPubText), os.Args[1])
	encryptedData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, hostPub, []byte(toSend), []byte("TELL HOST"))
	if err != nil {
		return
	}

	rply, err := SendToLG(hostPort, listenerNumber, "Input", InputQuestion{
		ClientPublicKeyText: string(clientPubText),
		Message:             encryptedData,
		ListenerKey:         listenerNumber,
	})
	if err != nil {
		return
	}

	msg, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, clientPriv, rply, []byte("TO PRINT"))
	if err != nil {
		return
	}

	fmt.Println(string(msg))
}

// SendToLG sends a message to the lazygit host
func SendToLG(port, listenerNumber string, selectFunction string, args interface{}) ([]byte, error) {
	client, err := rpc.Dial("tcp", "127.0.0.1:"+port)
	if err != nil {
		return nil, err
	}
	var out *[]byte
	err = client.Call("Listener"+listenerNumber+"."+selectFunction, args, &out)
	client.Close()
	return *out, err
}
