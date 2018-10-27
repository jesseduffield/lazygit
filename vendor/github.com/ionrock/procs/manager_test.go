package procs_test

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/ionrock/procs"
)

func TestManagerStartHelper(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	port := flag.String("p", "12212", "port to serve on")
	directory := flag.String("d", ".", "the directory of static file to host")
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*directory)))

	log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
	os.Exit(0)
}

func TestManagerStart(t *testing.T) {
	m := procs.NewManager()

	err := m.Start("test", fmt.Sprintf("%s -test.run=TestManagerStartHelper", os.Args[0]))
	if err != nil {
		t.Errorf("failed to start test process: %s", err)
	}

	if len(m.Processes) != 1 {
		t.Error("failed to add process")
	}

	err = m.Stop("test")
	if err != nil {
		t.Errorf("error stopping process: %s", err)
	}

	err = m.Remove("test")
	if err != nil {
		t.Errorf("error removing process: %s", err)
	}

	if len(m.Processes) != 0 {
		t.Error("failed to remove processes")
	}
}
