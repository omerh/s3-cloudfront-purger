package helpers

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

// GetAllCloudfrontDistributions to get all cloudfront distributions from aws
func GetAllCloudfrontDistributions() ([]*cloudfront.DistributionSummary, bool) {
	svc := cloudfront.New(session.New())
	input := &cloudfront.ListDistributionsInput{}

	result, err := svc.ListDistributions(input)
	if err != nil {
		fmt.Printf("Failed to get all distributions with error:\n%v\n", err)
		return nil, false
	}

	distributions := result.DistributionList.Items

	for result.DistributionList.NextMarker != nil {
		input = &cloudfront.ListDistributionsInput{
			Marker: result.DistributionList.NextMarker,
		}
		result, err = svc.ListDistributions(input)
		if err != nil {
			fmt.Printf("Failed to get all distributions with error:\n%v\n", err)
			return nil, false
		}
		for _, d := range result.DistributionList.Items {
			distributions = append(distributions, d)
		}
	}
	return distributions, true
}

// LocateDistributionFromS3 get all distributions and locate the one that use the s3 as origin
func LocateDistributionFromS3(distributions []*cloudfront.DistributionSummary, bucket string) []string {
	var foundDistributions []string
	for _, distribution := range distributions {
		if strings.HasPrefix(*distribution.Origins.Items[0].DomainName, bucket+".") {
			fmt.Printf("Found distribution id %v for bucket %v\n", *distribution.Id, bucket)
			foundDistributions = append(foundDistributions, *distribution.Id)
		}
	}
	return foundDistributions
}

// InvalidateCloudfrontDistribution purge distribution
func InvalidateCloudfrontDistribution(distibuitionID string) (bool, string) {
	svc := cloudfront.New(session.New())
	input := &cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(distibuitionID),
		InvalidationBatch: &cloudfront.InvalidationBatch{
			CallerReference: aws.String(fmt.Sprintf("Purge by s3-cloudfront-purger lambda at %v", time.Now())),
			Paths: &cloudfront.Paths{
				Quantity: aws.Int64(1),
				Items: []*string{
					aws.String("/*"),
				},
			},
		},
	}

	result, err := svc.CreateInvalidation(input)
	if err != nil {
		fmt.Printf("Problem creating an invalidation request to %v with the follwoing error:\n%v\n", distibuitionID, err)
		return false, ""
	}
	return true, *result.Invalidation.Id
}
