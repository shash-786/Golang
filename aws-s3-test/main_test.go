package main

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type MockClient struct {
	ListBucketsOutput  *s3.ListBucketsOutput
	CreateBucketOutput *s3.CreateBucketOutput
}

type MockUploader struct {
	UploadOutput *manager.UploadOutput
}

func (m *MockClient) ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	return m.ListBucketsOutput, nil
}

func (m *MockClient) CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	fmt.Printf("%s\n\n", *params.Bucket)
	return m.CreateBucketOutput, nil
}

func (m *MockUploader) Upload(ctx context.Context, input *s3.PutObjectInput, opts ...func(*manager.Uploader)) (*manager.UploadOutput, error) {
	fmt.Printf("key: %s\n", *input.Key)
	body, err := io.ReadAll(input.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll error: %v", err)
	}
	fmt.Println(string(body))
	return m.UploadOutput, nil
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

func TestUploadfiletoS3(t *testing.T) {
	err := uploadfiletoS3(&MockUploader{})
	if err != nil {
		t.Errorf("uploadfiletoS3 error: %v", err)
	}
}
