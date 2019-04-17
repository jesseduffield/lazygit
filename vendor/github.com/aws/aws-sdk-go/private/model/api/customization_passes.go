// +build codegen

package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type service struct {
	srcName string
	dstName string

	serviceVersion string
}

var mergeServices = map[string]service{
	"dynamodbstreams": {
		dstName: "dynamodb",
		srcName: "streams.dynamodb",
	},
	"wafregional": {
		dstName:        "waf",
		srcName:        "waf-regional",
		serviceVersion: "2015-08-24",
	},
}

var serviceAliaseNames = map[string]string{
	"costandusagereportservice": "CostandUsageReportService",
	"elasticloadbalancing":      "ELB",
	"elasticloadbalancingv2":    "ELBV2",
	"config":                    "ConfigService",
}

func (a *API) setServiceAliaseName() {
	if newName, ok := serviceAliaseNames[a.PackageName()]; ok {
		a.name = newName
	}
}

// customizationPasses Executes customization logic for the API by package name.
func (a *API) customizationPasses() {
	var svcCustomizations = map[string]func(*API){
		"s3":         s3Customizations,
		"s3control":  s3ControlCustomizations,
		"cloudfront": cloudfrontCustomizations,
		"rds":        rdsCustomizations,

		// Disable endpoint resolving for services that require customer
		// to provide endpoint them selves.
		"cloudsearchdomain": disableEndpointResolving,
		"iotdataplane":      disableEndpointResolving,

		// MTurk smoke test is invalid. The service requires AWS account to be
		// linked to Amazon Mechanical Turk Account.
		"mturk": supressSmokeTest,

		// Backfill the authentication type for cognito identity and sts.
		// Removes the need for the customizations in these services.
		"cognitoidentity": backfillAuthType("none",
			"GetId",
			"GetOpenIdToken",
			"UnlinkIdentity",
			"GetCredentialsForIdentity",
		),
		"sts": backfillAuthType("none",
			"AssumeRoleWithSAML",
			"AssumeRoleWithWebIdentity",
		),
	}

	for k := range mergeServices {
		svcCustomizations[k] = mergeServicesCustomizations
	}

	if fn := svcCustomizations[a.PackageName()]; fn != nil {
		fn(a)
	}
}

func supressSmokeTest(a *API) {
	a.SmokeTests.TestCases = []SmokeTestCase{}
}

// s3Customizations customizes the API generation to replace values specific to S3.
func s3Customizations(a *API) {
	var strExpires *Shape

	var keepContentMD5Ref = map[string]struct{}{
		"PutObjectInput":  {},
		"UploadPartInput": {},
	}

	for name, s := range a.Shapes {
		// Remove ContentMD5 members unless specified otherwise.
		if _, keep := keepContentMD5Ref[name]; !keep {
			if _, have := s.MemberRefs["ContentMD5"]; have {
				delete(s.MemberRefs, "ContentMD5")
			}
		}

		// Generate getter methods for API operation fields used by customizations.
		for _, refName := range []string{"Bucket", "SSECustomerKey", "CopySourceSSECustomerKey"} {
			if ref, ok := s.MemberRefs[refName]; ok {
				ref.GenerateGetter = true
			}
		}

		// Expires should be a string not time.Time since the format is not
		// enforced by S3, and any value can be set to this field outside of the SDK.
		if strings.HasSuffix(name, "Output") {
			if ref, ok := s.MemberRefs["Expires"]; ok {
				if strExpires == nil {
					newShape := *ref.Shape
					strExpires = &newShape
					strExpires.Type = "string"
					strExpires.refs = []*ShapeRef{}
				}
				ref.Shape.removeRef(ref)
				ref.Shape = strExpires
				ref.Shape.refs = append(ref.Shape.refs, &s.MemberRef)
			}
		}
	}
	s3CustRemoveHeadObjectModeledErrors(a)
}

// S3 HeadObject API call incorrect models NoSuchKey as valid
// error code that can be returned. This operation does not
// return error codes, all error codes are derived from HTTP
// status codes.
//
// aws/aws-sdk-go#1208
func s3CustRemoveHeadObjectModeledErrors(a *API) {
	op, ok := a.Operations["HeadObject"]
	if !ok {
		return
	}
	op.Documentation += `
//
// See http://docs.aws.amazon.com/AmazonS3/latest/API/ErrorResponses.html#RESTErrorResponses
// for more information on returned errors.`
	op.ErrorRefs = []ShapeRef{}
}

// S3 service operations with an AccountId need accessors to be generated for
// them so the fields can be dynamically accessed without reflection.
func s3ControlCustomizations(a *API) {
	for opName, op := range a.Operations {
		// Add moving AccountId into the hostname instead of header.
		if ref, ok := op.InputRef.Shape.MemberRefs["AccountId"]; ok {
			if op.Endpoint != nil {
				fmt.Fprintf(os.Stderr, "S3 Control, %s, model already defining endpoint trait, remove this customization.\n", opName)
			}

			op.Endpoint = &EndpointTrait{HostPrefix: "{AccountId}."}
			ref.HostLabel = true
		}
	}
}

// cloudfrontCustomizations customized the API generation to replace values
// specific to CloudFront.
func cloudfrontCustomizations(a *API) {
	// MaxItems members should always be integers
	for _, s := range a.Shapes {
		if ref, ok := s.MemberRefs["MaxItems"]; ok {
			ref.ShapeName = "Integer"
			ref.Shape = a.Shapes["Integer"]
		}
	}
}

// mergeServicesCustomizations references any duplicate shapes from DynamoDB
func mergeServicesCustomizations(a *API) {
	info := mergeServices[a.PackageName()]

	p := strings.Replace(a.path, info.srcName, info.dstName, -1)

	if info.serviceVersion != "" {
		index := strings.LastIndex(p, string(filepath.Separator))
		files, _ := ioutil.ReadDir(p[:index])
		if len(files) > 1 {
			panic("New version was introduced")
		}
		p = p[:index] + "/" + info.serviceVersion
	}

	file := filepath.Join(p, "api-2.json")

	serviceAPI := API{}
	serviceAPI.Attach(file)
	serviceAPI.Setup()

	for n := range a.Shapes {
		if _, ok := serviceAPI.Shapes[n]; ok {
			a.Shapes[n].resolvePkg = SDKImportRoot + "/service/" + info.dstName
		}
	}
}

// rdsCustomizations are customization for the service/rds. This adds non-modeled fields used for presigning.
func rdsCustomizations(a *API) {
	inputs := []string{
		"CopyDBSnapshotInput",
		"CreateDBInstanceReadReplicaInput",
		"CopyDBClusterSnapshotInput",
		"CreateDBClusterInput",
	}
	for _, input := range inputs {
		if ref, ok := a.Shapes[input]; ok {
			ref.MemberRefs["SourceRegion"] = &ShapeRef{
				Documentation: docstring(`SourceRegion is the source region where the resource exists. This is not sent over the wire and is only used for presigning. This value should always have the same region as the source ARN.`),
				ShapeName:     "String",
				Shape:         a.Shapes["String"],
				Ignore:        true,
			}
			ref.MemberRefs["DestinationRegion"] = &ShapeRef{
				Documentation: docstring(`DestinationRegion is used for presigning the request to a given region.`),
				ShapeName:     "String",
				Shape:         a.Shapes["String"],
			}
		}
	}
}

func disableEndpointResolving(a *API) {
	a.Metadata.NoResolveEndpoint = true
}

func backfillAuthType(typ string, opNames ...string) func(*API) {
	return func(a *API) {
		for _, opName := range opNames {
			op, ok := a.Operations[opName]
			if !ok {
				panic("unable to backfill auth-type for unknown operation " + opName)
			}
			if v := op.AuthType; len(v) != 0 {
				fmt.Fprintf(os.Stderr, "unable to backfill auth-type for %s, already set, %s", opName, v)
				continue
			}

			op.AuthType = typ
		}
	}
}
