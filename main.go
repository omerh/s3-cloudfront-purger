package main

import (
	"context"
	"fmt"
	"os"
	"s3-cloudfront-purger/helpers"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, s3Event events.S3Event) {
	fmt.Printf("Retrive event on bucket %v\n", s3Event.Records[0].S3.Bucket.Name)
	// Get bucket name from s3Event
	s3bucket := s3Event.Records[0].S3.Bucket.Name

	// Get all cloudfront distributions
	distributions, success := helpers.GetAllCloudfrontDistributions()
	if !success {
		os.Exit(1)
	}

	// locate all distributions that the s3 bucket name is its origin
	distributionIDs := helpers.LocateDistributionFromS3(distributions, s3bucket)
	if len(distributionIDs) > 0 {
		for _, distributionID := range distributionIDs {
			// invalidate distribution
			invalidated, invalidateID := helpers.InvalidateCloudfrontDistribution(distributionID)
			if !invalidated {
				fmt.Printf("Failed to purge the distribution list %v\n", distributionID)
				os.Exit(1)
			} else {
				fmt.Printf("Invalidation distribution id %v, ivalidate id is %v\n", distributionID, invalidateID)
			}
		}
	} else {
		fmt.Printf("Failed to find a distribution id for bucket %v\n", s3bucket)
		os.Exit(1)
	}
}

func main() {
	lambda.Start(handler)
}
