package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type MockClient struct {
	ListBucketsOutput  *s3.ListBucketsOutput
	CreateBucketOutput *s3.CreateBucketOutput
}

func (m *MockClient) ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	return m.ListBucketsOutput, nil
}

func (m *MockClient) CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	fmt.Println(*params.Bucket)
	return m.CreateBucketOutput, nil
}

func TestCreateS3Bucket(t *testing.T) {
	err := createS3Bucket(&MockClient{
		ListBucketsOutput: &s3.ListBucketsOutput{
			Buckets: []types.Bucket{
				{
					Name: aws.String("test-bucket1"),
				},
				{
					Name: aws.String("test-bucket2"),
				},
			},
		},
		CreateBucketOutput: &s3.CreateBucketOutput{},
	})
	if err != nil {
		t.Errorf("error in createS3Bucket: %v", err)
	}
}
