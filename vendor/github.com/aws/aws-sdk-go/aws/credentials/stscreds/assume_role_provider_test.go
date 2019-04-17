package stscreds

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
)

type stubSTS struct {
	TestInput func(*sts.AssumeRoleInput)
}

func (s *stubSTS) AssumeRole(input *sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
	if s.TestInput != nil {
		s.TestInput(input)
	}
	expiry := time.Now().Add(60 * time.Minute)
	return &sts.AssumeRoleOutput{
		Credentials: &sts.Credentials{
			// Just reflect the role arn to the provider.
			AccessKeyId:     input.RoleArn,
			SecretAccessKey: aws.String("assumedSecretAccessKey"),
			SessionToken:    aws.String("assumedSessionToken"),
			Expiration:      &expiry,
		},
	}, nil
}

func TestAssumeRoleProvider(t *testing.T) {
	stub := &stubSTS{}
	p := &AssumeRoleProvider{
		Client:  stub,
		RoleARN: "roleARN",
	}

	creds, err := p.Retrieve()
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	if e, a := "roleARN", creds.AccessKeyID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "assumedSecretAccessKey", creds.SecretAccessKey; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "assumedSessionToken", creds.SessionToken; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestAssumeRoleProvider_WithTokenCode(t *testing.T) {
	stub := &stubSTS{
		TestInput: func(in *sts.AssumeRoleInput) {
			if e, a := "0123456789", *in.SerialNumber; e != a {
				t.Errorf("expect %v, got %v", e, a)
			}
			if e, a := "code", *in.TokenCode; e != a {
				t.Errorf("expect %v, got %v", e, a)
			}
		},
	}
	p := &AssumeRoleProvider{
		Client:       stub,
		RoleARN:      "roleARN",
		SerialNumber: aws.String("0123456789"),
		TokenCode:    aws.String("code"),
	}

	creds, err := p.Retrieve()
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	if e, a := "roleARN", creds.AccessKeyID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "assumedSecretAccessKey", creds.SecretAccessKey; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "assumedSessionToken", creds.SessionToken; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestAssumeRoleProvider_WithTokenProvider(t *testing.T) {
	stub := &stubSTS{
		TestInput: func(in *sts.AssumeRoleInput) {
			if e, a := "0123456789", *in.SerialNumber; e != a {
				t.Errorf("expect %v, got %v", e, a)
			}
			if e, a := "code", *in.TokenCode; e != a {
				t.Errorf("expect %v, got %v", e, a)
			}
		},
	}
	p := &AssumeRoleProvider{
		Client:       stub,
		RoleARN:      "roleARN",
		SerialNumber: aws.String("0123456789"),
		TokenProvider: func() (string, error) {
			return "code", nil
		},
	}

	creds, err := p.Retrieve()
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	if e, a := "roleARN", creds.AccessKeyID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "assumedSecretAccessKey", creds.SecretAccessKey; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "assumedSessionToken", creds.SessionToken; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestAssumeRoleProvider_WithTokenProviderError(t *testing.T) {
	stub := &stubSTS{
		TestInput: func(in *sts.AssumeRoleInput) {
			t.Errorf("API request should not of been called")
		},
	}
	p := &AssumeRoleProvider{
		Client:       stub,
		RoleARN:      "roleARN",
		SerialNumber: aws.String("0123456789"),
		TokenProvider: func() (string, error) {
			return "", fmt.Errorf("error occurred")
		},
	}

	creds, err := p.Retrieve()
	if err == nil {
		t.Errorf("expect error")
	}

	if v := creds.AccessKeyID; len(v) != 0 {
		t.Errorf("expect empty, got %v", v)
	}
	if v := creds.SecretAccessKey; len(v) != 0 {
		t.Errorf("expect empty, got %v", v)
	}
	if v := creds.SessionToken; len(v) != 0 {
		t.Errorf("expect empty, got %v", v)
	}
}

func TestAssumeRoleProvider_MFAWithNoToken(t *testing.T) {
	stub := &stubSTS{
		TestInput: func(in *sts.AssumeRoleInput) {
			t.Errorf("API request should not of been called")
		},
	}
	p := &AssumeRoleProvider{
		Client:       stub,
		RoleARN:      "roleARN",
		SerialNumber: aws.String("0123456789"),
	}

	creds, err := p.Retrieve()
	if err == nil {
		t.Errorf("expect error")
	}

	if v := creds.AccessKeyID; len(v) != 0 {
		t.Errorf("expect empty, got %v", v)
	}
	if v := creds.SecretAccessKey; len(v) != 0 {
		t.Errorf("expect empty, got %v", v)
	}
	if v := creds.SessionToken; len(v) != 0 {
		t.Errorf("expect empty, got %v", v)
	}
}

func BenchmarkAssumeRoleProvider(b *testing.B) {
	stub := &stubSTS{}
	p := &AssumeRoleProvider{
		Client:  stub,
		RoleARN: "roleARN",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := p.Retrieve(); err != nil {
			b.Fatal(err)
		}
	}
}
