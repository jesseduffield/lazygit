package errors

import (
	"reflect"
	"testing"
)

var createdBy = `panic: hello!

goroutine 54 [running]:
runtime.panic(0x35ce40, 0xc208039db0)
	/0/c/go/src/pkg/runtime/panic.c:279 +0xf5
github.com/loopj/bugsnag-example-apps/go/revelapp/app/controllers.func·001()
	/0/go/src/github.com/loopj/bugsnag-example-apps/go/revelapp/app/controllers/app.go:13 +0x74
net/http.(*Server).Serve(0xc20806c780, 0x910c88, 0xc20803e168, 0x0, 0x0)
	/0/c/go/src/pkg/net/http/server.go:1698 +0x91
created by github.com/loopj/bugsnag-example-apps/go/revelapp/app/controllers.App.Index
	/0/go/src/github.com/loopj/bugsnag-example-apps/go/revelapp/app/controllers/app.go:14 +0x3e

goroutine 16 [IO wait]:
net.runtime_pollWait(0x911c30, 0x72, 0x0)
	/0/c/go/src/pkg/runtime/netpoll.goc:146 +0x66
net.(*pollDesc).Wait(0xc2080ba990, 0x72, 0x0, 0x0)
	/0/c/go/src/pkg/net/fd_poll_runtime.go:84 +0x46
net.(*pollDesc).WaitRead(0xc2080ba990, 0x0, 0x0)
	/0/c/go/src/pkg/net/fd_poll_runtime.go:89 +0x42
net.(*netFD).accept(0xc2080ba930, 0x58be30, 0x0, 0x9103f0, 0x23)
	/0/c/go/src/pkg/net/fd_unix.go:409 +0x343
net.(*TCPListener).AcceptTCP(0xc20803e168, 0x8, 0x0, 0x0)
	/0/c/go/src/pkg/net/tcpsock_posix.go:234 +0x5d
net.(*TCPListener).Accept(0xc20803e168, 0x0, 0x0, 0x0, 0x0)
	/0/c/go/src/pkg/net/tcpsock_posix.go:244 +0x4b
github.com/revel/revel.Run(0xe6d9)
	/0/go/src/github.com/revel/revel/server.go:113 +0x926
main.main()
	/0/go/src/github.com/loopj/bugsnag-example-apps/go/revelapp/app/tmp/main.go:109 +0xe1a
`

var normalSplit = `panic: hello!

goroutine 54 [running]:
runtime.panic(0x35ce40, 0xc208039db0)
	/0/c/go/src/pkg/runtime/panic.c:279 +0xf5
github.com/loopj/bugsnag-example-apps/go/revelapp/app/controllers.func·001()
	/0/go/src/github.com/loopj/bugsnag-example-apps/go/revelapp/app/controllers/app.go:13 +0x74
net/http.(*Server).Serve(0xc20806c780, 0x910c88, 0xc20803e168, 0x0, 0x0)
	/0/c/go/src/pkg/net/http/server.go:1698 +0x91

goroutine 16 [IO wait]:
net.runtime_pollWait(0x911c30, 0x72, 0x0)
	/0/c/go/src/pkg/runtime/netpoll.goc:146 +0x66
net.(*pollDesc).Wait(0xc2080ba990, 0x72, 0x0, 0x0)
	/0/c/go/src/pkg/net/fd_poll_runtime.go:84 +0x46
net.(*pollDesc).WaitRead(0xc2080ba990, 0x0, 0x0)
	/0/c/go/src/pkg/net/fd_poll_runtime.go:89 +0x42
net.(*netFD).accept(0xc2080ba930, 0x58be30, 0x0, 0x9103f0, 0x23)
	/0/c/go/src/pkg/net/fd_unix.go:409 +0x343
net.(*TCPListener).AcceptTCP(0xc20803e168, 0x8, 0x0, 0x0)
	/0/c/go/src/pkg/net/tcpsock_posix.go:234 +0x5d
net.(*TCPListener).Accept(0xc20803e168, 0x0, 0x0, 0x0, 0x0)
	/0/c/go/src/pkg/net/tcpsock_posix.go:244 +0x4b
github.com/revel/revel.Run(0xe6d9)
	/0/go/src/github.com/revel/revel/server.go:113 +0x926
main.main()
	/0/go/src/github.com/loopj/bugsnag-example-apps/go/revelapp/app/tmp/main.go:109 +0xe1a
`

var lastGoroutine = `panic: hello!

goroutine 16 [IO wait]:
net.runtime_pollWait(0x911c30, 0x72, 0x0)
	/0/c/go/src/pkg/runtime/netpoll.goc:146 +0x66
net.(*pollDesc).Wait(0xc2080ba990, 0x72, 0x0, 0x0)
	/0/c/go/src/pkg/net/fd_poll_runtime.go:84 +0x46
net.(*pollDesc).WaitRead(0xc2080ba990, 0x0, 0x0)
	/0/c/go/src/pkg/net/fd_poll_runtime.go:89 +0x42
net.(*netFD).accept(0xc2080ba930, 0x58be30, 0x0, 0x9103f0, 0x23)
	/0/c/go/src/pkg/net/fd_unix.go:409 +0x343
net.(*TCPListener).AcceptTCP(0xc20803e168, 0x8, 0x0, 0x0)
	/0/c/go/src/pkg/net/tcpsock_posix.go:234 +0x5d
net.(*TCPListener).Accept(0xc20803e168, 0x0, 0x0, 0x0, 0x0)
	/0/c/go/src/pkg/net/tcpsock_posix.go:244 +0x4b
github.com/revel/revel.Run(0xe6d9)
	/0/go/src/github.com/revel/revel/server.go:113 +0x926
main.main()
	/0/go/src/github.com/loopj/bugsnag-example-apps/go/revelapp/app/tmp/main.go:109 +0xe1a

goroutine 54 [running]:
runtime.panic(0x35ce40, 0xc208039db0)
	/0/c/go/src/pkg/runtime/panic.c:279 +0xf5
github.com/loopj/bugsnag-example-apps/go/revelapp/app/controllers.func·001()
	/0/go/src/github.com/loopj/bugsnag-example-apps/go/revelapp/app/controllers/app.go:13 +0x74
net/http.(*Server).Serve(0xc20806c780, 0x910c88, 0xc20803e168, 0x0, 0x0)
	/0/c/go/src/pkg/net/http/server.go:1698 +0x91
`

var result = []StackFrame{
	StackFrame{File: "/0/c/go/src/pkg/runtime/panic.c", LineNumber: 279, Name: "panic", Package: "runtime"},
	StackFrame{File: "/0/go/src/github.com/loopj/bugsnag-example-apps/go/revelapp/app/controllers/app.go", LineNumber: 13, Name: "func.001", Package: "github.com/loopj/bugsnag-example-apps/go/revelapp/app/controllers"},
	StackFrame{File: "/0/c/go/src/pkg/net/http/server.go", LineNumber: 1698, Name: "(*Server).Serve", Package: "net/http"},
}

var resultCreatedBy = append(result,
	StackFrame{File: "/0/go/src/github.com/loopj/bugsnag-example-apps/go/revelapp/app/controllers/app.go", LineNumber: 14, Name: "App.Index", Package: "github.com/loopj/bugsnag-example-apps/go/revelapp/app/controllers", ProgramCounter: 0x0})

func TestParsePanic(t *testing.T) {

	todo := map[string]string{
		"createdBy":     createdBy,
		"normalSplit":   normalSplit,
		"lastGoroutine": lastGoroutine,
	}

	for key, val := range todo {
		Err, err := ParsePanic(val)

		if err != nil {
			t.Fatal(err)
		}

		if Err.TypeName() != "panic" {
			t.Errorf("Wrong type: %s", Err.TypeName())
		}

		if Err.Error() != "hello!" {
			t.Errorf("Wrong message: %s", Err.TypeName())
		}

		if Err.StackFrames()[0].Func() != nil {
			t.Errorf("Somehow managed to find a func...")
		}

		result := result
		if key == "createdBy" {
			result = resultCreatedBy
		}

		if !reflect.DeepEqual(Err.StackFrames(), result) {
			t.Errorf("Wrong stack for %s: %#v", key, Err.StackFrames())
		}
	}
}
