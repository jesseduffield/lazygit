// +build integration

package s3integ

import (
	"fmt"

	"github.com/aws/aws-sdk-go/awstesting/integration"
	"github.com/aws/aws-sdk-go/service/s3"
)

// BucketPrefix is the root prefix of integration test buckets.
const BucketPrefix = "aws-sdk-go-integration"

// GenerateBucketName returns a unique bucket name.
func GenerateBucketName() string {
	return fmt.Sprintf("%s-%s",
		BucketPrefix, integration.UniqueID())
}

// SetupTest returns a test bucket created for the integration tests.
func SetupTest(svc *s3.S3, bucketName string) (err error) {

	fmt.Println("Setup: Creating test bucket,", bucketName)
	_, err = svc.CreateBucket(&s3.CreateBucketInput{Bucket: &bucketName})
	if err != nil {
		return fmt.Errorf("failed to create bucket %s, %v", bucketName, err)
	}

	fmt.Println("Setup: Waiting for bucket to exist,", bucketName)
	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{Bucket: &bucketName})
	if err != nil {
		return fmt.Errorf("failed waiting for bucket %s to be created, %v",
			bucketName, err)
	}

	return nil
}

// CleanupTest deletes the contents of a S3 bucket, before deleting the bucket
// it self.
func CleanupTest(svc *s3.S3, bucketName string) error {
	errs := []error{}

	fmt.Println("TearDown: Deleting objects from test bucket,", bucketName)
	err := svc.ListObjectsPages(
		&s3.ListObjectsInput{Bucket: &bucketName},
		func(page *s3.ListObjectsOutput, lastPage bool) bool {
			for _, o := range page.Contents {
				_, err := svc.DeleteObject(&s3.DeleteObjectInput{
					Bucket: &bucketName,
					Key:    o.Key,
				})
				if err != nil {
					errs = append(errs, err)
				}
			}
			return true
		},
	)
	if err != nil {
		return fmt.Errorf("failed to list objects, %s, %v", bucketName, err)
	}

	fmt.Println("TearDown: Deleting partial uploads from test bucket,", bucketName)
	err = svc.ListMultipartUploadsPages(
		&s3.ListMultipartUploadsInput{Bucket: &bucketName},
		func(page *s3.ListMultipartUploadsOutput, lastPage bool) bool {
			for _, u := range page.Uploads {
				svc.AbortMultipartUpload(&s3.AbortMultipartUploadInput{
					Bucket:   &bucketName,
					Key:      u.Key,
					UploadId: u.UploadId,
				})
			}
			return true
		},
	)
	if err != nil {
		return fmt.Errorf("failed to list multipart objects, %s, %v", bucketName, err)
	}

	if len(errs) != 0 {
		return fmt.Errorf("failed to delete objects, %s", errs)
	}

	fmt.Println("TearDown: Deleting test bucket,", bucketName)
	if _, err = svc.DeleteBucket(&s3.DeleteBucketInput{Bucket: &bucketName}); err != nil {
		return fmt.Errorf("failed to delete test bucket, %s", bucketName)
	}

	return nil
}
