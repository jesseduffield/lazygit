// +build integration

package s3_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/awstesting/integration"
	"github.com/aws/aws-sdk-go/awstesting/integration/s3integ"
	"github.com/aws/aws-sdk-go/service/s3"
)

const integBucketPrefix = "aws-sdk-go-integration"

var bucketName *string
var svc *s3.S3

func TestMain(m *testing.M) {
	sess := integration.SessionWithDefaultRegion("us-west-2")
	svc = s3.New(sess)
	bucketName = aws.String(s3integ.GenerateBucketName())
	if err := s3integ.SetupTest(svc, *bucketName); err != nil {
		panic(err)
	}

	var result int
	defer func() {
		if err := s3integ.CleanupTest(svc, *bucketName); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, "S3 integrationt tests paniced,", r)
			result = 1
		}
		os.Exit(result)
	}()

	result = m.Run()
}

func putTestFile(t *testing.T, filename, key string) {
	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("failed to open testfile, %v", err)
	}
	defer f.Close()

	putTestContent(t, f, key)
}

func putTestContent(t *testing.T, reader io.ReadSeeker, key string) {
	fmt.Println(bucketName, key, svc)
	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: bucketName,
		Key:    aws.String(key),
		Body:   reader,
	})
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}
