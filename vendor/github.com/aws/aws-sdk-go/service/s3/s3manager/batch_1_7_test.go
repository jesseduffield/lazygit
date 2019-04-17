// +build go1.7

package s3manager

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
)

// #1790 bug
func TestBatchDeleteContext(t *testing.T) {
	cases := []struct {
		objects  []BatchDeleteObject
		size     int
		expected int
		ctx      aws.Context
		closeAt  int
		errCheck func(error) (string, bool)
	}{
		{
			[]BatchDeleteObject{
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("1"),
						Bucket: aws.String("bucket1"),
					},
				},
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("2"),
						Bucket: aws.String("bucket2"),
					},
				},
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("3"),
						Bucket: aws.String("bucket3"),
					},
				},
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("4"),
						Bucket: aws.String("bucket4"),
					},
				},
			},
			1,
			0,
			aws.BackgroundContext(),
			0,
			func(err error) (string, bool) {
				batchErr, ok := err.(*BatchError)
				if !ok {
					return "not BatchError type", false
				}

				errs := batchErr.Errors
				if len(errs) != 4 {
					return fmt.Sprintf("expected 1, but received %d", len(errs)), false
				}

				for _, tempErr := range errs {
					aerr, ok := tempErr.OrigErr.(awserr.Error)
					if !ok {
						return "not awserr.Error type", false
					}

					if code := aerr.Code(); code != request.CanceledErrorCode {
						return fmt.Sprintf("expected %q, but received %q", request.CanceledErrorCode, code), false
					}
				}
				return "", true
			},
		},
	}

	count := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		count++
	}))

	svc := &mockS3Client{S3: buildS3SvcClient(server.URL)}
	for i, c := range cases {
		ctx, done := context.WithCancel(c.ctx)
		defer done()
		if i == c.closeAt {
			done()
		}

		batcher := BatchDelete{
			Client:    svc,
			BatchSize: c.size,
		}

		err := batcher.Delete(ctx, &DeleteObjectsIterator{Objects: c.objects})

		if msg, ok := c.errCheck(err); !ok {
			t.Error(msg)
		}

		if count != c.expected {
			t.Errorf("Case %d: expected %d, but received %d", i, c.expected, count)
		}

		count = 0
	}
}
