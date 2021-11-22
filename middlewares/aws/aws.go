package awsCommon

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	myTweetConfig "mytweet/config"
)

func GetAwsCredential() (aws.Config, error) {

	id := myTweetConfig.GetConfig().Get("aws.AccessKeyID").(string)
	secret := myTweetConfig.GetConfig().Get("aws.SecretKey").(string)
	customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		env := myTweetConfig.GetConfig().Get("user.name").(string)
		if env == "development" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           myTweetConfig.GetConfig().Get("aws.EndpointURL").(string),
				SigningRegion: "ap-northeast-1",
			}, nil
		} else {
			return aws.Endpoint{
				PartitionID:   "aws",
				SigningRegion: "ap-northeast-1",
			}, nil
		}
	})

	appCreds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(id, secret, "dummy"))
	return config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(appCreds), config.WithEndpointResolver(customResolver))
}
