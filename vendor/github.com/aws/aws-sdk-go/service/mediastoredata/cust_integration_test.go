// +build integration

package mediastoredata_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/awstesting/integration"
	"github.com/aws/aws-sdk-go/service/mediastore"
	"github.com/aws/aws-sdk-go/service/mediastoredata"
)

func TestInteg_DescribeEndpoint(t *testing.T) {
	const containerName = "awsgosdkteamintegcontainer"

	sess := integration.Session.Copy()
	if v := aws.StringValue(sess.Config.Region); len(v) == 0 {
		sess.Config.Region = aws.String("us-east-1")
	}

	ctrlSvc := mediastore.New(sess)
	descResp, err := ctrlSvc.DescribeContainer(&mediastore.DescribeContainerInput{
		ContainerName: aws.String(containerName),
	})
	if err != nil {
		t.Fatalf("failed to get mediastore container endpoint, %v", err)
	}

	dataSvc := mediastoredata.New(sess, &aws.Config{
		Endpoint: descResp.Container.Endpoint,
	})
	_, err = dataSvc.ListItems(&mediastoredata.ListItemsInput{})
	if err != nil {
		t.Fatalf("failed to make medaistoredata API call, %v", err)
	}
}
