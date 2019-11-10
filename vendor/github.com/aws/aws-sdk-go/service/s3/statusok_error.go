package s3

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/internal/sdkio"
)

func copyMultipartStatusOKUnmarhsalError(r *request.Request) {
	b, err := ioutil.ReadAll(r.HTTPResponse.Body)
	if err != nil {
		r.Error = awserr.NewRequestFailure(
			awserr.New(request.ErrCodeSerialization, "unable to read response body", err),
			r.HTTPResponse.StatusCode,
			r.RequestID,
		)
		return
	}
	body := bytes.NewReader(b)
	r.HTTPResponse.Body = ioutil.NopCloser(body)
	defer body.Seek(0, sdkio.SeekStart)

	if body.Len() == 0 {
		// If there is no body don't attempt to parse the body.
		return
	}

	unmarshalError(r)
	if err, ok := r.Error.(awserr.Error); ok && err != nil {
		if err.Code() == request.ErrCodeSerialization {
			r.Error = nil
			return
		}
		r.HTTPResponse.StatusCode = http.StatusServiceUnavailable
	}
}
