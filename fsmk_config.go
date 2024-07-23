package mqueue

import "github.com/aws/aws-sdk-go-v2/aws"

// MKAWSConfig defines the configuration for the MK cluster.
var MKAWSConfig = func() aws.Config {
	return *aws.NewConfig()
}()

// FSAWSConfig defines the configuration for the FS cluster.
var FSAWSConfig = func() aws.Config {
	return *aws.NewConfig()
}()

// MSMConfig defines the configuration for the MSM cluster.
var MSMConfig = func() aws.Config {
	return *aws.NewConfig()
}()

// LocalStackConfig defines the configuration for the MSM cluster.
var LocalStackConfig = func() aws.Config {
	return *aws.NewConfig()
}()
