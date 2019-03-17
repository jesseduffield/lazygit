package gitcredentialhelper

import (
	"fmt"
	"net"
)

// getFreePort returns a free system port
func getFreePort() string {
	portCheckFrom := 5000
	for {
		portCheckFrom++
		port := fmt.Sprintf("%v", portCheckFrom)
		conn, err := net.Dial("tcp", "127.0.0.1:"+port)
		if err == nil {
			go conn.Close()
			continue
		}
		return port
	}
}
