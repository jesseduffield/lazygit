package s3

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/private/protocol/xml/xmlutil"
)

type xmlErrorResponse struct {
	XMLName xml.Name `xml:"Error"`
	Code    string   `xml:"Code"`
	Message string   `xml:"Message"`
}

func unmarshalError(r *request.Request) {
	defer r.HTTPResponse.Body.Close()
	defer io.Copy(ioutil.Discard, r.HTTPResponse.Body)

	// Bucket exists in a different region, and request needs
	// to be made to the correct region.
	if r.HTTPResponse.StatusCode == http.StatusMovedPermanently {
		msg := fmt.Sprintf(
			"incorrect region, the bucket is not in '%s' region at endpoint '%s'",
			aws.StringValue(r.Config.Region),
			aws.StringValue(r.Config.Endpoint),
		)
		if v := r.HTTPResponse.Header.Get("x-amz-bucket-region"); len(v) != 0 {
			msg += fmt.Sprintf(", bucket is in '%s' region", v)
		}
		r.Error = awserr.NewRequestFailure(
			awserr.New("BucketRegionError", msg, nil),
			r.HTTPResponse.StatusCode,
			r.RequestID,
		)
		return
	}

	// Attempt to parse error from body if it is known
	var errResp xmlErrorResponse
	err := xmlutil.UnmarshalXMLError(&errResp, r.HTTPResponse.Body)
	if err == io.EOF {
		// Only capture the error if an unmarshal error occurs that is not EOF,
		// because S3 might send an error without a error message which causes
		// the XML unmarshal to fail with EOF.
		err = nil
	}
	if err != nil {
		r.Error = awserr.NewRequestFailure(
			awserr.New(request.ErrCodeSerialization,
				"failed to unmarshal error message", err),
			r.HTTPResponse.StatusCode,
			r.RequestID,
		)
		return
	}

	// Fallback to status code converted to message if still no error code
	if len(errResp.Code) == 0 {
		statusText := http.StatusText(r.HTTPResponse.StatusCode)
		errResp.Code = strings.Replace(statusText, " ", "", -1)
		errResp.Message = statusText
	}

	r.Error = awserr.NewRequestFailure(
		awserr.New(errResp.Code, errResp.Message, err),
		r.HTTPResponse.StatusCode,
		r.RequestID,
	)
}

// A RequestFailure provides access to the S3 Request ID and Host ID values
// returned from API operation errors. Getting the error as a string will
// return the formated error with the same information as awserr.RequestFailure,
// while also adding the HostID value from the response.
type RequestFailure interface {
	awserr.RequestFailure

	// Host ID is the S3 Host ID needed for debug, and contacting support
	HostID() string
}
