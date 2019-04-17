// +build integration

package s3_test

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

func TestInteg_SelectObjectContent(t *testing.T) {
	keyName := "selectObject.csv"
	putTestFile(t, filepath.Join("testdata", "positive_select.csv"), keyName)

	resp, err := svc.SelectObjectContent(&s3.SelectObjectContentInput{
		Bucket:         bucketName,
		Key:            &keyName,
		Expression:     aws.String("Select * from S3Object"),
		ExpressionType: aws.String(s3.ExpressionTypeSql),
		InputSerialization: &s3.InputSerialization{
			CSV: &s3.CSVInput{
				FieldDelimiter: aws.String(","),
				FileHeaderInfo: aws.String(s3.FileHeaderInfoIgnore),
			},
		},
		OutputSerialization: &s3.OutputSerialization{
			CSV: &s3.CSVOutput{
				FieldDelimiter: aws.String(","),
			},
		},
	})
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer resp.EventStream.Close()

	var sum int64
	var processed int64
	for event := range resp.EventStream.Events() {
		switch tv := event.(type) {
		case *s3.RecordsEvent:
			sum += int64(len(tv.Payload))
		case *s3.StatsEvent:
			processed = *tv.Details.BytesProcessed
		}
	}

	if sum == 0 {
		t.Errorf("expect selected content, got none")
	}

	if processed == 0 {
		t.Errorf("expect selected status bytes processed, got none")
	}

	if err := resp.EventStream.Err(); err != nil {
		t.Fatalf("exect no error, %v", err)
	}
}

func TestInteg_SelectObjectContent_Error(t *testing.T) {
	keyName := "negativeSelect.csv"

	buf := make([]byte, 0, 1024*1024*6)
	buf = append(buf, []byte("name,number\n")...)
	line := []byte("jj,0\n")
	for i := 0; i < (cap(buf)/len(line))-2; i++ {
		buf = append(buf, line...)
	}
	buf = append(buf, []byte("gg,NaN\n")...)

	putTestContent(t, bytes.NewReader(buf), keyName)

	resp, err := svc.SelectObjectContent(&s3.SelectObjectContentInput{
		Bucket:         bucketName,
		Key:            &keyName,
		Expression:     aws.String("SELECT name FROM S3Object WHERE cast(number as int) < 1"),
		ExpressionType: aws.String(s3.ExpressionTypeSql),
		InputSerialization: &s3.InputSerialization{
			CSV: &s3.CSVInput{
				FileHeaderInfo: aws.String(s3.FileHeaderInfoUse),
			},
		},
		OutputSerialization: &s3.OutputSerialization{
			CSV: &s3.CSVOutput{
				FieldDelimiter: aws.String(","),
			},
		},
	})
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer resp.EventStream.Close()

	var sum int64
	for event := range resp.EventStream.Events() {
		switch tv := event.(type) {
		case *s3.RecordsEvent:
			sum += int64(len(tv.Payload))
		}
	}

	if sum == 0 {
		t.Errorf("expect selected content")
	}

	err = resp.EventStream.Err()
	if err == nil {
		t.Fatalf("exepct error")
	}

	aerr := err.(awserr.Error)
	if a := aerr.Code(); len(a) == 0 {
		t.Errorf("expect, error code")
	}
	if a := aerr.Message(); len(a) == 0 {
		t.Errorf("expect, error message")
	}
}

func TestInteg_SelectObjectContent_Stream(t *testing.T) {
	keyName := "selectGopher.csv"

	buf := `name,number
gopher,0
ᵷodɥǝɹ,1
`
	// Put a mock CSV file to the S3 bucket so that its contents can be
	// selected.
	putTestContent(t, strings.NewReader(buf), keyName)

	// Make the Select Object Content API request using the object uploaded.
	resp, err := svc.SelectObjectContent(&s3.SelectObjectContentInput{
		Bucket:         bucketName,
		Key:            &keyName,
		Expression:     aws.String("SELECT name FROM S3Object WHERE cast(number as int) < 1"),
		ExpressionType: aws.String(s3.ExpressionTypeSql),
		InputSerialization: &s3.InputSerialization{
			CSV: &s3.CSVInput{
				FileHeaderInfo: aws.String(s3.FileHeaderInfoUse),
			},
		},
		OutputSerialization: &s3.OutputSerialization{
			CSV: &s3.CSVOutput{},
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed making API request, %v\n", err)
		return
	}
	defer resp.EventStream.Close()

	results, resultWriter := io.Pipe()
	go func() {
		defer resultWriter.Close()
		for event := range resp.EventStream.Events() {
			switch e := event.(type) {
			case *s3.RecordsEvent:
				resultWriter.Write(e.Payload)
			case *s3.StatsEvent:
				fmt.Printf("Processed %d bytes\n", *e.Details.BytesProcessed)
			}
		}
	}()

	// Printout the results
	resReader := csv.NewReader(results)
	for {
		record, err := resReader.Read()
		if err == io.EOF {
			break
		}
		fmt.Println(record)
	}

	if err := resp.EventStream.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "reading from event stream failed, %v\n", err)
	}
}
