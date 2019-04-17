package ec2metadata_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/awstesting"
	"github.com/aws/aws-sdk-go/awstesting/unit"
)

func TestClientOverrideDefaultHTTPClientTimeout(t *testing.T) {
	svc := ec2metadata.New(unit.Session)

	if e, a := http.DefaultClient, svc.Config.HTTPClient; e == a {
		t.Errorf("expect %v, not to equal %v", e, a)
	}

	if e, a := 5*time.Second, svc.Config.HTTPClient.Timeout; e != a {
		t.Errorf("expect %v to be %v", e, a)
	}
}

func TestClientNotOverrideDefaultHTTPClientTimeout(t *testing.T) {
	http.DefaultClient.Transport = &http.Transport{}
	defer func() {
		http.DefaultClient.Transport = nil
	}()

	svc := ec2metadata.New(unit.Session)

	if e, a := http.DefaultClient, svc.Config.HTTPClient; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

	tr := svc.Config.HTTPClient.Transport.(*http.Transport)
	if tr == nil {
		t.Fatalf("expect transport not to be nil")
	}
	if tr.Dial != nil {
		t.Errorf("expect dial to be nil, was not")
	}
}

func TestClientDisableOverrideDefaultHTTPClientTimeout(t *testing.T) {
	svc := ec2metadata.New(unit.Session, aws.NewConfig().WithEC2MetadataDisableTimeoutOverride(true))

	if e, a := http.DefaultClient, svc.Config.HTTPClient; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestClientOverrideDefaultHTTPClientTimeoutRace(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("us-east-1a"))
	}))

	cfg := aws.NewConfig().WithEndpoint(server.URL)
	runEC2MetadataClients(t, cfg, 100)
}

func TestClientOverrideDefaultHTTPClientTimeoutRaceWithTransport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("us-east-1a"))
	}))

	cfg := aws.NewConfig().WithEndpoint(server.URL).WithHTTPClient(&http.Client{
		Transport: http.DefaultTransport,
	})

	runEC2MetadataClients(t, cfg, 100)
}

func TestClientDisableIMDS(t *testing.T) {
	env := awstesting.StashEnv()
	defer awstesting.PopEnv(env)

	os.Setenv("AWS_EC2_METADATA_DISABLED", "true")

	svc := ec2metadata.New(unit.Session, &aws.Config{
		LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody),
	})
	resp, err := svc.Region()
	if err == nil {
		t.Fatalf("expect error, got none")
	}
	if len(resp) != 0 {
		t.Errorf("expect no response, got %v", resp)
	}

	aerr := err.(awserr.Error)
	if e, a := request.CanceledErrorCode, aerr.Code(); e != a {
		t.Errorf("expect %v error code, got %v", e, a)
	}
	if e, a := "AWS_EC2_METADATA_DISABLED", aerr.Message(); !strings.Contains(a, e) {
		t.Errorf("expect %v in error message, got %v", e, a)
	}
}

func runEC2MetadataClients(t *testing.T, cfg *aws.Config, atOnce int) {
	var wg sync.WaitGroup
	wg.Add(atOnce)
	for i := 0; i < atOnce; i++ {
		go func() {
			svc := ec2metadata.New(unit.Session, cfg)
			_, err := svc.Region()
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
