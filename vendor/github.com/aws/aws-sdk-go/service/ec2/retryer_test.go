package ec2

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/awstesting/unit"
)

func TestCustomRetryer(t *testing.T) {
	svc := New(unit.Session, &aws.Config{Region: aws.String("us-west-2")})
	if _, ok := svc.Client.Retryer.(retryer); !ok {
		t.Error("expected custom retryer, but received otherwise")
	}

	req, _ := svc.ModifyNetworkInterfaceAttributeRequest(&ModifyNetworkInterfaceAttributeInput{
		NetworkInterfaceId: aws.String("foo"),
	})

	duration := svc.Client.Retryer.RetryRules(req)
	if duration < time.Second*1 || duration > time.Second*2 {
		t.Errorf("expected duration to be between 1 and 2, but received %v", duration)
	}

	req.RetryCount = 15
	duration = svc.Client.Retryer.RetryRules(req)
	if duration < time.Second*5 || duration > time.Second*10 {
		t.Errorf("expected duration to be between 1 and 2, but received %v", duration)
	}

	svc = New(unit.Session, &aws.Config{Region: aws.String("us-west-2"), Retryer: client.DefaultRetryer{}})
	if _, ok := svc.Client.Retryer.(client.DefaultRetryer); !ok {
		t.Error("expected default retryer, but received otherwise")
	}
}
