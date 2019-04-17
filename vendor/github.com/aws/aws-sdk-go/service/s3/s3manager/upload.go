package s3manager

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// MaxUploadParts is the maximum allowed number of parts in a multi-part upload
// on Amazon S3.
const MaxUploadParts = 10000

// MinUploadPartSize is the minimum allowed part size when uploading a part to
// Amazon S3.
const MinUploadPartSize int64 = 1024 * 1024 * 5

// DefaultUploadPartSize is the default part size to buffer chunks of a
// payload into.
const DefaultUploadPartSize = MinUploadPartSize

// DefaultUploadConcurrency is the default number of goroutines to spin up when
// using Upload().
const DefaultUploadConcurrency = 5

// A MultiUploadFailure wraps a failed S3 multipart upload. An error returned
// will satisfy this interface when a multi part upload failed to upload all
// chucks to S3. In the case of a failure the UploadID is needed to operate on
// the chunks, if any, which were uploaded.
//
// Example:
//
//     u := s3manager.NewUploader(opts)
//     output, err := u.upload(input)
//     if err != nil {
//         if multierr, ok := err.(s3manager.MultiUploadFailure); ok {
//             // Process error and its associated uploadID
//             fmt.Println("Error:", multierr.Code(), multierr.Message(), multierr.UploadID())
//         } else {
//             // Process error generically
//             fmt.Println("Error:", err.Error())
//         }
//     }
//
type MultiUploadFailure interface {
	awserr.Error

	// Returns the upload id for the S3 multipart upload that failed.
	UploadID() string
}

// So that the Error interface type can be included as an anonymous field
// in the multiUploadError struct and not conflict with the error.Error() method.
type awsError awserr.Error

// A multiUploadError wraps the upload ID of a failed s3 multipart upload.
// Composed of BaseError for code, message, and original error
//
// Should be used for an error that occurred failing a S3 multipart upload,
// and a upload ID is available. If an uploadID is not available a more relevant
type multiUploadError struct {
	awsError

	// ID for multipart upload which failed.
	uploadID string
}

// Error returns the string representation of the error.
//
// See apierr.BaseError ErrorWithExtra for output format
//
// Satisfies the error interface.
func (m multiUploadError) Error() string {
	extra := fmt.Sprintf("upload id: %s", m.uploadID)
	return awserr.SprintError(m.Code(), m.Message(), extra, m.OrigErr())
}

// String returns the string representation of the error.
// Alias for Error to satisfy the stringer interface.
func (m multiUploadError) String() string {
	return m.Error()
}

// UploadID returns the id of the S3 upload which failed.
func (m multiUploadError) UploadID() string {
	return m.uploadID
}

// UploadOutput represents a response from the Upload() call.
type UploadOutput struct {
	// The URL where the object was uploaded to.
	Location string

	// The version of the object that was uploaded. Will only be populated if
	// the S3 Bucket is versioned. If the bucket is not versioned this field
	// will not be set.
	VersionID *string

	// The ID for a multipart upload to S3. In the case of an error the error
	// can be cast to the MultiUploadFailure interface to extract the upload ID.
	UploadID string
}

// WithUploaderRequestOptions appends to the Uploader's API request options.
func WithUploaderRequestOptions(opts ...request.Option) func(*Uploader) {
	return func(u *Uploader) {
		u.RequestOptions = append(u.RequestOptions, opts...)
	}
}

// The Uploader structure that calls Upload(). It is safe to call Upload()
// on this structure for multiple objects and across concurrent goroutines.
// Mutating the Uploader's properties is not safe to be done concurrently.
type Uploader struct {
	// The buffer size (in bytes) to use when buffering data into chunks and
	// sending them as parts to S3. The minimum allowed part size is 5MB, and
	// if this value is set to zero, the DefaultUploadPartSize value will be used.
	PartSize int64

	// The number of goroutines to spin up in parallel per call to Upload when
	// sending parts. If this is set to zero, the DefaultUploadConcurrency value
	// will be used.
	//
	// The concurrency pool is not shared between calls to Upload.
	Concurrency int

	// Setting this value to true will cause the SDK to avoid calling
	// AbortMultipartUpload on a failure, leaving all successfully uploaded
	// parts on S3 for manual recovery.
	//
	// Note that storing parts of an incomplete multipart upload counts towards
	// space usage on S3 and will add additional costs if not cleaned up.
	LeavePartsOnError bool

	// MaxUploadParts is the max number of parts which will be uploaded to S3.
	// Will be used to calculate the partsize of the object to be uploaded.
	// E.g: 5GB file, with MaxUploadParts set to 100, will upload the file
	// as 100, 50MB parts.
	// With a limited of s3.MaxUploadParts (10,000 parts).
	//
	// Defaults to package const's MaxUploadParts value.
	MaxUploadParts int

	// The client to use when uploading to S3.
	S3 s3iface.S3API

	// List of request options that will be passed down to individual API
	// operation requests made by the uploader.
	RequestOptions []request.Option
}

// NewUploader creates a new Uploader instance to upload objects to S3. Pass In
// additional functional options to customize the uploader's behavior. Requires a
// client.ConfigProvider in order to create a S3 service client. The session.Session
// satisfies the client.ConfigProvider interface.
//
// Example:
//     // The session the S3 Uploader will use
//     sess := session.Must(session.NewSession())
//
//     // Create an uploader with the session and default options
//     uploader := s3manager.NewUploader(sess)
//
//     // Create an uploader with the session and custom options
//     uploader := s3manager.NewUploader(session, func(u *s3manager.Uploader) {
//          u.PartSize = 64 * 1024 * 1024 // 64MB per part
//     })
func NewUploader(c client.ConfigProvider, options ...func(*Uploader)) *Uploader {
	u := &Uploader{
		S3:                s3.New(c),
		PartSize:          DefaultUploadPartSize,
		Concurrency:       DefaultUploadConcurrency,
		LeavePartsOnError: false,
		MaxUploadParts:    MaxUploadParts,
	}

	for _, option := range options {
		option(u)
	}

	return u
}

// NewUploaderWithClient creates a new Uploader instance to upload objects to S3. Pass in
// additional functional options to customize the uploader's behavior. Requires
// a S3 service client to make S3 API calls.
//
// Example:
//     // The session the S3 Uploader will use
//     sess := session.Must(session.NewSession())
//
//     // S3 service client the Upload manager will use.
//     s3Svc := s3.New(sess)
//
//     // Create an uploader with S3 client and default options
//     uploader := s3manager.NewUploaderWithClient(s3Svc)
//
//     // Create an uploader with S3 client and custom options
//     uploader := s3manager.NewUploaderWithClient(s3Svc, func(u *s3manager.Uploader) {
//          u.PartSize = 64 * 1024 * 1024 // 64MB per part
//     })
func NewUploaderWithClient(svc s3iface.S3API, options ...func(*Uploader)) *Uploader {
	u := &Uploader{
		S3:                svc,
		PartSize:          DefaultUploadPartSize,
		Concurrency:       DefaultUploadConcurrency,
		LeavePartsOnError: false,
		MaxUploadParts:    MaxUploadParts,
	}

	for _, option := range options {
		option(u)
	}

	return u
}

// Upload uploads an object to S3, intelligently buffering large files into
// smaller chunks and sending them in parallel across multiple goroutines. You
// can configure the buffer size and concurrency through the Uploader's parameters.
//
// Additional functional options can be provided to configure the individual
// upload. These options are copies of the Uploader instance Upload is called from.
// Modifying the options will not impact the original Uploader instance.
//
// Use the WithUploaderRequestOptions helper function to pass in request
// options that will be applied to all API operations made with this uploader.
//
// It is safe to call this method concurrently across goroutines.
//
// Example:
//     // Upload input parameters
//     upParams := &s3manager.UploadInput{
//         Bucket: &bucketName,
//         Key:    &keyName,
//         Body:   file,
//     }
//
//     // Perform an upload.
//     result, err := uploader.Upload(upParams)
//
//     // Perform upload with options different than the those in the Uploader.
//     result, err := uploader.Upload(upParams, func(u *s3manager.Uploader) {
//          u.PartSize = 10 * 1024 * 1024 // 10MB part size
//          u.LeavePartsOnError = true    // Don't delete the parts if the upload fails.
//     })
func (u Uploader) Upload(input *UploadInput, options ...func(*Uploader)) (*UploadOutput, error) {
	return u.UploadWithContext(aws.BackgroundContext(), input, options...)
}

// UploadWithContext uploads an object to S3, intelligently buffering large
// files into smaller chunks and sending them in parallel across multiple
// goroutines. You can configure the buffer size and concurrency through the
// Uploader's parameters.
//
// UploadWithContext is the same as Upload with the additional support for
// Context input parameters. The Context must not be nil. A nil Context will
// cause a panic. Use the context to add deadlining, timeouts, etc. The
// UploadWithContext may create sub-contexts for individual underlying requests.
//
// Additional functional options can be provided to configure the individual
// upload. These options are copies of the Uploader instance Upload is called from.
// Modifying the options will not impact the original Uploader instance.
//
// Use the WithUploaderRequestOptions helper function to pass in request
// options that will be applied to all API operations made with this uploader.
//
// It is safe to call this method concurrently across goroutines.
func (u Uploader) UploadWithContext(ctx aws.Context, input *UploadInput, opts ...func(*Uploader)) (*UploadOutput, error) {
	i := uploader{in: input, cfg: u, ctx: ctx}

	for _, opt := range opts {
		opt(&i.cfg)
	}
	i.cfg.RequestOptions = append(i.cfg.RequestOptions, request.WithAppendUserAgent("S3Manager"))

	return i.upload()
}

// UploadWithIterator will upload a batched amount of objects to S3. This operation uses
// the iterator pattern to know which object to upload next. Since this is an interface this
// allows for custom defined functionality.
//
// Example:
//	svc:= s3manager.NewUploader(sess)
//
//	objects := []BatchUploadObject{
//		{
//			Object:	&s3manager.UploadInput {
//				Key: aws.String("key"),
//				Bucket: aws.String("bucket"),
//			},
//		},
//	}
//
//	iter := &s3manager.UploadObjectsIterator{Objects: objects}
//	if err := svc.UploadWithIterator(aws.BackgroundContext(), iter); err != nil {
//		return err
//	}
func (u Uploader) UploadWithIterator(ctx aws.Context, iter BatchUploadIterator, opts ...func(*Uploader)) error {
	var errs []Error
	for iter.Next() {
		object := iter.UploadObject()
		if _, err := u.UploadWithContext(ctx, object.Object, opts...); err != nil {
			s3Err := Error{
				OrigErr: err,
				Bucket:  object.Object.Bucket,
				Key:     object.Object.Key,
			}

			errs = append(errs, s3Err)
		}

		if object.After == nil {
			continue
		}

		if err := object.After(); err != nil {
			s3Err := Error{
				OrigErr: err,
				Bucket:  object.Object.Bucket,
				Key:     object.Object.Key,
			}

			errs = append(errs, s3Err)
		}
	}

	if len(errs) > 0 {
		return NewBatchError("BatchedUploadIncomplete", "some objects have failed to upload.", errs)
	}
	return nil
}

// internal structure to manage an upload to S3.
type uploader struct {
	ctx aws.Context
	cfg Uploader

	in *UploadInput

	readerPos int64 // current reader position
	totalSize int64 // set to -1 if the size is not known

	bufferPool sync.Pool
}

// internal logic for deciding whether to upload a single part or use a
// multipart upload.
func (u *uploader) upload() (*UploadOutput, error) {
	u.init()

	if u.cfg.PartSize < MinUploadPartSize {
		msg := fmt.Sprintf("part size must be at least %d bytes", MinUploadPartSize)
		return nil, awserr.New("ConfigError", msg, nil)
	}

	// Do one read to determine if we have more than one part
	reader, _, part, err := u.nextReader()
	if err == io.EOF { // single part
		return u.singlePart(reader)
	} else if err != nil {
		return nil, awserr.New("ReadRequestBody", "read upload data failed", err)
	}

	mu := multiuploader{uploader: u}
	return mu.upload(reader, part)
}

// init will initialize all default options.
func (u *uploader) init() {
	if u.cfg.Concurrency == 0 {
		u.cfg.Concurrency = DefaultUploadConcurrency
	}
	if u.cfg.PartSize == 0 {
		u.cfg.PartSize = DefaultUploadPartSize
	}
	if u.cfg.MaxUploadParts == 0 {
		u.cfg.MaxUploadParts = MaxUploadParts
	}

	u.bufferPool = sync.Pool{
		New: func() interface{} { return make([]byte, u.cfg.PartSize) },
	}

	// Try to get the total size for some optimizations
	u.initSize()
}

// initSize tries to detect the total stream size, setting u.totalSize. If
// the size is not known, totalSize is set to -1.
func (u *uploader) initSize() {
	u.totalSize = -1

	switch r := u.in.Body.(type) {
	case io.Seeker:
		n, err := aws.SeekerLen(r)
		if err != nil {
			return
		}
		u.totalSize = n

		// Try to adjust partSize if it is too small and account for
		// integer division truncation.
		if u.totalSize/u.cfg.PartSize >= int64(u.cfg.MaxUploadParts) {
			// Add one to the part size to account for remainders
			// during the size calculation. e.g odd number of bytes.
			u.cfg.PartSize = (u.totalSize / int64(u.cfg.MaxUploadParts)) + 1
		}
	}
}

// nextReader returns a seekable reader representing the next packet of data.
// This operation increases the shared u.readerPos counter, but note that it
// does not need to be wrapped in a mutex because nextReader is only called
// from the main thread.
func (u *uploader) nextReader() (io.ReadSeeker, int, []byte, error) {
	type readerAtSeeker interface {
		io.ReaderAt
		io.ReadSeeker
	}
	switch r := u.in.Body.(type) {
	case readerAtSeeker:
		var err error

		n := u.cfg.PartSize
		if u.totalSize >= 0 {
			bytesLeft := u.totalSize - u.readerPos

			if bytesLeft <= u.cfg.PartSize {
				err = io.EOF
				n = bytesLeft
			}
		}

		reader := io.NewSectionReader(r, u.readerPos, n)
		u.readerPos += n

		return reader, int(n), nil, err

	default:
		part := u.bufferPool.Get().([]byte)
		n, err := readFillBuf(r, part)
		u.readerPos += int64(n)

		return bytes.NewReader(part[0:n]), n, part, err
	}
}

func readFillBuf(r io.Reader, b []byte) (offset int, err error) {
	for offset < len(b) && err == nil {
		var n int
		n, err = r.Read(b[offset:])
		offset += n
	}

	return offset, err
}

// singlePart contains upload logic for uploading a single chunk via
// a regular PutObject request. Multipart requests require at least two
// parts, or at least 5MB of data.
func (u *uploader) singlePart(buf io.ReadSeeker) (*UploadOutput, error) {
	params := &s3.PutObjectInput{}
	awsutil.Copy(params, u.in)
	params.Body = buf

	// Need to use request form because URL generated in request is
	// used in return.
	req, out := u.cfg.S3.PutObjectRequest(params)
	req.SetContext(u.ctx)
	req.ApplyOptions(u.cfg.RequestOptions...)
	if err := req.Send(); err != nil {
		return nil, err
	}

	url := req.HTTPRequest.URL.String()
	return &UploadOutput{
		Location:  url,
		VersionID: out.VersionId,
	}, nil
}

// internal structure to manage a specific multipart upload to S3.
type multiuploader struct {
	*uploader
	wg       sync.WaitGroup
	m        sync.Mutex
	err      error
	uploadID string
	parts    completedParts
}

// keeps track of a single chunk of data being sent to S3.
type chunk struct {
	buf  io.ReadSeeker
	part []byte
	num  int64
}

// completedParts is a wrapper to make parts sortable by their part number,
// since S3 required this list to be sent in sorted order.
type completedParts []*s3.CompletedPart

func (a completedParts) Len() int           { return len(a) }
func (a completedParts) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a completedParts) Less(i, j int) bool { return *a[i].PartNumber < *a[j].PartNumber }

// upload will perform a multipart upload using the firstBuf buffer containing
// the first chunk of data.
func (u *multiuploader) upload(firstBuf io.ReadSeeker, firstPart []byte) (*UploadOutput, error) {
	params := &s3.CreateMultipartUploadInput{}
	awsutil.Copy(params, u.in)

	// Create the multipart
	resp, err := u.cfg.S3.CreateMultipartUploadWithContext(u.ctx, params, u.cfg.RequestOptions...)
	if err != nil {
		return nil, err
	}
	u.uploadID = *resp.UploadId

	// Create the workers
	ch := make(chan chunk, u.cfg.Concurrency)
	for i := 0; i < u.cfg.Concurrency; i++ {
		u.wg.Add(1)
		go u.readChunk(ch)
	}

	// Send part 1 to the workers
	var num int64 = 1
	ch <- chunk{buf: firstBuf, part: firstPart, num: num}

	// Read and queue the rest of the parts
	for u.geterr() == nil && err == nil {
		num++
		// This upload exceeded maximum number of supported parts, error now.
		if num > int64(u.cfg.MaxUploadParts) || num > int64(MaxUploadParts) {
			var msg string
			if num > int64(u.cfg.MaxUploadParts) {
				msg = fmt.Sprintf("exceeded total allowed configured MaxUploadParts (%d). Adjust PartSize to fit in this limit",
					u.cfg.MaxUploadParts)
			} else {
				msg = fmt.Sprintf("exceeded total allowed S3 limit MaxUploadParts (%d). Adjust PartSize to fit in this limit",
					MaxUploadParts)
			}
			u.seterr(awserr.New("TotalPartsExceeded", msg, nil))
			break
		}

		var reader io.ReadSeeker
		var nextChunkLen int
		var part []byte
		reader, nextChunkLen, part, err = u.nextReader()

		if err != nil && err != io.EOF {
			u.seterr(awserr.New(
				"ReadRequestBody",
				"read multipart upload data failed",
				err))
			break
		}

		if nextChunkLen == 0 {
			// No need to upload empty part, if file was empty to start
			// with empty single part would of been created and never
			// started multipart upload.
			break
		}

		ch <- chunk{buf: reader, part: part, num: num}
	}

	// Close the channel, wait for workers, and complete upload
	close(ch)
	u.wg.Wait()
	complete := u.complete()

	if err := u.geterr(); err != nil {
		return nil, &multiUploadError{
			awsError: awserr.New(
				"MultipartUpload",
				"upload multipart failed",
				err),
			uploadID: u.uploadID,
		}
	}

	// Create a presigned URL of the S3 Get Object in order to have parity with
	// single part upload.
	getReq, _ := u.cfg.S3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: u.in.Bucket,
		Key:    u.in.Key,
	})
	getReq.Config.Credentials = credentials.AnonymousCredentials
	uploadLocation, _, _ := getReq.PresignRequest(1)

	return &UploadOutput{
		Location:  uploadLocation,
		VersionID: complete.VersionId,
		UploadID:  u.uploadID,
	}, nil
}

// readChunk runs in worker goroutines to pull chunks off of the ch channel
// and send() them as UploadPart requests.
func (u *multiuploader) readChunk(ch chan chunk) {
	defer u.wg.Done()
	for {
		data, ok := <-ch

		if !ok {
			break
		}

		if u.geterr() == nil {
			if err := u.send(data); err != nil {
				u.seterr(err)
			}
		}
	}
}

// send performs an UploadPart request and keeps track of the completed
// part information.
func (u *multiuploader) send(c chunk) error {
	params := &s3.UploadPartInput{
		Bucket:               u.in.Bucket,
		Key:                  u.in.Key,
		Body:                 c.buf,
		UploadId:             &u.uploadID,
		SSECustomerAlgorithm: u.in.SSECustomerAlgorithm,
		SSECustomerKey:       u.in.SSECustomerKey,
		PartNumber:           &c.num,
	}
	resp, err := u.cfg.S3.UploadPartWithContext(u.ctx, params, u.cfg.RequestOptions...)
	// put the byte array back into the pool to conserve memory
	u.bufferPool.Put(c.part)
	if err != nil {
		return err
	}

	n := c.num
	completed := &s3.CompletedPart{ETag: resp.ETag, PartNumber: &n}

	u.m.Lock()
	u.parts = append(u.parts, completed)
	u.m.Unlock()

	return nil
}

// geterr is a thread-safe getter for the error object
func (u *multiuploader) geterr() error {
	u.m.Lock()
	defer u.m.Unlock()

	return u.err
}

// seterr is a thread-safe setter for the error object
func (u *multiuploader) seterr(e error) {
	u.m.Lock()
	defer u.m.Unlock()

	u.err = e
}

// fail will abort the multipart unless LeavePartsOnError is set to true.
func (u *multiuploader) fail() {
	if u.cfg.LeavePartsOnError {
		return
	}

	params := &s3.AbortMultipartUploadInput{
		Bucket:   u.in.Bucket,
		Key:      u.in.Key,
		UploadId: &u.uploadID,
	}
	_, err := u.cfg.S3.AbortMultipartUploadWithContext(u.ctx, params, u.cfg.RequestOptions...)
	if err != nil {
		logMessage(u.cfg.S3, aws.LogDebug, fmt.Sprintf("failed to abort multipart upload, %v", err))
	}
}

// complete successfully completes a multipart upload and returns the response.
func (u *multiuploader) complete() *s3.CompleteMultipartUploadOutput {
	if u.geterr() != nil {
		u.fail()
		return nil
	}

	// Parts must be sorted in PartNumber order.
	sort.Sort(u.parts)

	params := &s3.CompleteMultipartUploadInput{
		Bucket:          u.in.Bucket,
		Key:             u.in.Key,
		UploadId:        &u.uploadID,
		MultipartUpload: &s3.CompletedMultipartUpload{Parts: u.parts},
	}
	resp, err := u.cfg.S3.CompleteMultipartUploadWithContext(u.ctx, params, u.cfg.RequestOptions...)
	if err != nil {
		u.seterr(err)
		u.fail()
	}

	return resp
}
