package awsendpointdiscoverytest

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/awstesting/unit"
)

func TestEndpointDiscovery(t *testing.T) {
	svc := New(unit.Session, &aws.Config{
		EnableEndpointDiscovery: aws.Bool(true),
	})
	svc.Handlers.Clear()
	svc.Handlers.Send.PushBack(mockSendDescEndpoint)

	var descCount int32
	svc.Handlers.Complete.PushBack(func(r *request.Request) {
		if r.Operation.Name != opDescribeEndpoints {
			return
		}
		atomic.AddInt32(&descCount, 1)
	})

	for i := 0; i < 2; i++ {
		req, _ := svc.TestDiscoveryIdentifiersRequiredRequest(
			&TestDiscoveryIdentifiersRequiredInput{
				Sdk: aws.String("sdk"),
			},
		)
		req.Handlers.Send.PushBack(func(r *request.Request) {
			if e, a := "http://foo", r.HTTPRequest.URL.String(); e != a {
				t.Errorf("expected %q, but received %q", e, a)
			}
		})
		if err := req.Send(); err != nil {
			t.Fatal(err)
		}
	}

	if e, a := int32(1), atomic.LoadInt32(&descCount); e != a {
		t.Errorf("expect desc endpoint called %d, got %d", e, a)
	}
}

func TestAsyncEndpointDiscovery(t *testing.T) {
	t.Parallel()

	svc := New(unit.Session, &aws.Config{
		EnableEndpointDiscovery: aws.Bool(true),
	})
	svc.Handlers.Clear()

	var firstAsyncReq sync.WaitGroup
	firstAsyncReq.Add(1)
	svc.Handlers.Build.PushBack(func(r *request.Request) {
		if r.Operation.Name == opDescribeEndpoints {
			firstAsyncReq.Wait()
		}
	})
	svc.Handlers.Send.PushBack(mockSendDescEndpoint)

	req, _ := svc.TestDiscoveryOptionalRequest(&TestDiscoveryOptionalInput{
		Sdk: aws.String("sdk"),
	})
	const clientHost = "awsendpointdiscoverytestservice.mock-region.amazonaws.com"
	req.Handlers.Send.PushBack(func(r *request.Request) {
		if e, a := clientHost, r.HTTPRequest.URL.Host; e != a {
			t.Errorf("expected %q, but received %q", e, a)
		}
	})
	req.Handlers.Complete.PushBack(func(r *request.Request) {
		firstAsyncReq.Done()
	})
	if err := req.Send(); err != nil {
		t.Fatal(err)
	}

	var cacheUpdated bool
	for s := time.Now().Add(10 * time.Second); s.After(time.Now()); {
		// Wait for the cache to be updated before making second request.
		if svc.endpointCache.Has(req.Operation.Name) {
			cacheUpdated = true
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if !cacheUpdated {
		t.Fatalf("expect endpoint cache to be updated, was not")
	}

	req, _ = svc.TestDiscoveryOptionalRequest(&TestDiscoveryOptionalInput{
		Sdk: aws.String("sdk"),
	})
	req.Handlers.Send.PushBack(func(r *request.Request) {
		if e, a := "http://foo", r.HTTPRequest.URL.String(); e != a {
			t.Errorf("expected %q, but received %q", e, a)
		}
	})
	if err := req.Send(); err != nil {
		t.Fatal(err)
	}
}

func removeHandlers(h request.Handlers, removeSendHandlers bool) request.Handlers {
	if removeSendHandlers {
		h.Send.Clear()
	}
	h.Unmarshal.Clear()
	h.UnmarshalStream.Clear()
	h.UnmarshalMeta.Clear()
	h.UnmarshalError.Clear()
	h.Validate.Clear()
	h.Complete.Clear()
	h.ValidateResponse.Clear()
	return h
}

func mockSendDescEndpoint(r *request.Request) {
	if r.Operation.Name != opDescribeEndpoints {
		return
	}

	out, _ := r.Data.(*DescribeEndpointsOutput)
	out.Endpoints = []*Endpoint{
		{
			Address:              aws.String("http://foo"),
			CachePeriodInMinutes: aws.Int64(5),
		},
	}
	r.Data = out
}
