// +build example

package main

import (
	"log"
	"os"
	"sync/atomic"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type CustomReader struct {
	fp   *os.File
	size int64
	read int64
}

func (r *CustomReader) Read(p []byte) (int, error) {
	return r.fp.Read(p)
}

func (r *CustomReader) ReadAt(p []byte, off int64) (int, error) {
	n, err := r.fp.ReadAt(p, off)
	if err != nil {
		return n, err
	}

	// Got the length have read( or means has uploaded), and you can construct your message
	atomic.AddInt64(&r.read, int64(n))

	// I have no idea why the read length need to be div 2,
	// maybe the request read once when Sign and actually send call ReadAt again
	// It works for me
	log.Printf("total read:%d    progress:%d%%\n", r.read/2, int(float32(r.read*100/2)/float32(r.size)))

	return n, err
}

func (r *CustomReader) Seek(offset int64, whence int) (int64, error) {
	return r.fp.Seek(offset, whence)
}

func main() {
	if len(os.Args) < 4 {
		log.Println("USAGE ERROR: AWS_REGION=us-east-1 go run putObjWithProcess.go <credential> <bucket> <key for object> <local file name>")
		return
	}

	credential := os.Args[1]
	bucket := os.Args[2]
	key := os.Args[3]
	fileName := os.Args[4]

	creds := credentials.NewSharedCredentials(credential, "default")
	if _, err := creds.Get(); err != nil {
		log.Println("ERROR:", err)
		return
	}

	sess := session.New(&aws.Config{
		Credentials: creds,
	})

	file, err := os.Open(fileName)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println("ERROR:", err)
		return
	}

	reader := &CustomReader{
		fp:   file,
		size: fileInfo.Size(),
	}

	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024
		u.LeavePartsOnError = true
	})

	output, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   reader,
	})

	if err != nil {
		log.Println("ERROR:", err)
		return
	}

	log.Println(output.Location)
}
