package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Client interface {
	ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
	CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
}

const (
	region     string = "ap-south-1"
	BUCKETNAME string = "go-aws-demo-shash"
	filename   string = "test.txt"
)

func main() {
	var err error
	var client S3Client

	if client, err = initClient(); err != nil {
		log.Fatalf("initClient error: %v", err)
	}

	if err = createS3Bucket(client); err != nil {
		log.Fatalf("createS3Bucket error: %v", err)
	}
}

func createS3Bucket(client S3Client) error {
	var (
		out *s3.ListBucketsOutput
		err error
	)

	if out, err = client.ListBuckets(context.Background(), &s3.ListBucketsInput{}); err != nil {
		return fmt.Errorf("ListBuckets Error: %v", err)
	}

	if len(out.Buckets) != 0 {
		for _, bucket := range out.Buckets {
			fmt.Println(*bucket.Name)
			if *bucket.Name == BUCKETNAME {
				fmt.Println("Bucket already exists!")
				return nil
			}
		}
	}

	_, err = client.CreateBucket(context.Background(), &s3.CreateBucketInput{
		Bucket: aws.String(BUCKETNAME),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		},
	})
	if err != nil {
		return fmt.Errorf("create bucket error: %v", err)
	}
	return nil
}

func initClient() (S3Client, error) {
	var (
		cnfg aws.Config
		err  error
	)
	if cnfg, err = config.LoadDefaultConfig(context.TODO(), config.WithRegion(region)); err != nil {
		return nil, fmt.Errorf("LoadDefaultConfig error: %v", err)
	}
	client := s3.NewFromConfig(cnfg)
	return client, nil
}
