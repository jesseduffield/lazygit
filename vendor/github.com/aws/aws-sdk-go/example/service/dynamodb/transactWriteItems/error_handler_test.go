// +build example

package transaction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/awstesting/unit"
)

const errStatusCode = 400
const requestId = "requestId1"

func TestNewTxErrorAwareDynamoDBClient(t *testing.T) {
	sess := unit.Session
	svc := NewTxErrorAwareDynamoDBClient(sess)

	if svc.Handlers.UnmarshalError.Len() != 1 {
		t.Errorf("expected 1 UnmarshallErrorHandler, got %v", svc.Handlers.UnmarshalError.Len())
	}
	if svc.Handlers.UnmarshalError.Swap(TxAwareErrorUnmarshallerName, request.NamedHandler{}) == false {
		t.Errorf("expected to contain %s, got none", TxAwareErrorUnmarshallerName)
	}
}

func TestTxAwareUnmarshalError(t *testing.T) {
	input := map[string]struct {
		enc interface{}
		err TxRequestFailure
	}{
		"Error response without CancellationReasons": {
			jsonErrorResponse{
				Code:    "com.amazonaws.dynamodb.v20120810#ResourceNotFoundException",
				Message: "Requested resource not found",
			},
			txRequestError{
				awsError:            awserr.New("ResourceNotFoundException", "Requested resource not found", nil),
				statusCode:          errStatusCode,
				requestID:           requestId,
				cancellationReasons: nil,
			},
		},
		"Error response with empty CancellationReasons": {
			jsonTxErrorResponse{
				Code:                "com.amazonaws.dynamodb.v20120810#TransactionCanceledException",
				Message:             "Transaction cancelled, please refer cancellation reasons for specific reasons [ConditionalCheckFailed, None, None]",
				CancellationReasons: []dynamodb.CancellationReason{},
			},
			txRequestError{
				awsError:            awserr.New("TransactionCanceledException", "Transaction cancelled, please refer cancellation reasons for specific reasons [ConditionalCheckFailed, None, None]", nil),
				statusCode:          errStatusCode,
				requestID:           requestId,
				cancellationReasons: []dynamodb.CancellationReason{},
			},
		},
		"Error response with non-empty CancellationReasons": {
			jsonTxErrorResponse{
				Code:    "com.amazonaws.dynamodb.v20120810#TransactionCanceledException",
				Message: "Transaction cancelled, please refer cancellation reasons for specific reasons [ConditionalCheckFailed, None, None]",
				CancellationReasons: []dynamodb.CancellationReason{
					{
						Code: aws.String("ConditionalCheckFailed"),
						Item: map[string]*dynamodb.AttributeValue{
							"hk":   {S: aws.String("hkVal1")},
							"attr": {S: aws.String("attrVal1")},
						},
						Message: aws.String("The conditional request failed"),
					},
				},
			},
			txRequestError{
				awsError:   awserr.New("TransactionCanceledException", "Transaction cancelled, please refer cancellation reasons for specific reasons [ConditionalCheckFailed, None, None]", nil),
				statusCode: errStatusCode,
				requestID:  requestId,
				cancellationReasons: []dynamodb.CancellationReason{
					{
						Code: aws.String("ConditionalCheckFailed"),
						Item: map[string]*dynamodb.AttributeValue{
							"hk":   {S: aws.String("hkVal1")},
							"attr": {S: aws.String("attrVal1")},
						},
						Message: aws.String("The conditional request failed"),
					},
				},
			},
		},
	}

	for name, in := range input {
		t.Run(name, func(t *testing.T) {
			if err := validateUnmarshallError(in.enc, in.err); err != nil {
				t.Errorf("%s: expected nil, got %v", name, err)
			}
		})
	}
}

func validateUnmarshallError(enc interface{}, err TxRequestFailure) error {
	req := &request.Request{
		HTTPResponse: &http.Response{
			StatusCode: errStatusCode,
			Body:       newBufferCloser(encode(enc)),
		},
		RequestID: requestId,
	}
	TxAwareUnmarshalError(req)

	if aerr, ok := req.Error.(TxRequestFailure); ok {
		if err.RequestID() != aerr.RequestID() {
			return fmt.Errorf("expected %v, got %v", err.RequestID(), aerr.RequestID())
		}
		if err.StatusCode() != aerr.StatusCode() {
			return fmt.Errorf("expected %v, got %v", err.StatusCode(), aerr.StatusCode())
		}
		if err.Message() != aerr.Message() {
			return fmt.Errorf("expected %v, got %v", err.Message(), aerr.Message())
		}
		if err.Code() != aerr.Code() {
			return fmt.Errorf("expected %v, got %v", err.Code(), aerr.Code())
		}
		if err.OrigErr() != aerr.OrigErr() {
			return fmt.Errorf("expected %v, got %v", err.OrigErr(), aerr.OrigErr())
		}
		if !reflect.DeepEqual(err.Error(), aerr.Error()) {
			return fmt.Errorf("expected %v, got %v", err.Error(), aerr.Error())
		}
		if !reflect.DeepEqual(err.CancellationReasons(), aerr.CancellationReasons()) {
			return fmt.Errorf("expected %v, got %v", err.CancellationReasons(), aerr.CancellationReasons())
		}
	} else {
		return fmt.Errorf("expected type 'TxRequestFailure', got %T", req.Error)
	}
	return nil
}

func encode(v interface{}) []byte {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(&v)
	return buf.Bytes()
}

// Implementation of io.ReadCloser backed by bytes.Buffer
type bufferCloser struct {
	bytes.Buffer
}

func newBufferCloser(data []byte) *bufferCloser {
	return &bufferCloser{*bytes.NewBuffer(data)}
}

func (b *bufferCloser) Close() error {
	b.Reset()
	return nil
}

// Define error response without the cancellation reason
type jsonErrorResponse struct {
	Code    string `json:"__type"`
	Message string `json:"message"`
}
