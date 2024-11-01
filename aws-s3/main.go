package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const BUCKETNAME string = "go-aws-demo-bucket-shash"

func main() {
	var s3client *s3.Client
	var err error
	var region string

	flag.StringVar(&region, "r", "", "region of bucket creation")
	flag.Parse()

	if region == "" {
		region = "ap-south-1"
	}

	if s3client, err = InitClient(region); err != nil {
		log.Fatalf("error initializing the client: %v", err)
	}

	if err = createS3bucket(s3client, region); err != nil {
		log.Fatalf("error creating bucket: %v", err)
	}

	if err = uploadtoS3bucket(s3client); err != nil {
		log.Fatal(err)
	}

	if err = downloadfromS3bucket(s3client); err != nil {
		log.Fatal(err)
	}
}

func createS3bucket(client *s3.Client, region string) error {
	bucket_list, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("Error in retrieving the bucket list: %v", err)
	}

	if len(bucket_list.Buckets) != 0 {
		for _, bucket := range bucket_list.Buckets {
			if *bucket.Name == BUCKETNAME {
				fmt.Println("Bucket Already exists!")
				return nil
			}
		}
	}

	_, err = client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
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

func get_content_type(file *os.File) string {
	extension := filepath.Ext(file.Name())
	cnt_type := mime.TypeByExtension(extension)
	if cnt_type == "" {
		return ""
	}
	return cnt_type
}

func uploadtoS3bucket(client *s3.Client) error {
	var (
		err          error
		file         *os.File
		content_type string
	)

	if file, err = os.Open("sample.csv"); err != nil {
		return fmt.Errorf("error in opening the file: %v", err)
	}
	defer file.Close()

	if content_type = get_content_type(file); content_type == "" {
		return fmt.Errorf("content_type unknown")
	}
	/*
	 * USE client.PutObject for small upload sizes
	 * it is a low level function
	 * not safe/memory opotimized for large
	 * file sizes
	 */

	uploader := manager.NewUploader(client)
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(BUCKETNAME),
		Key:         aws.String(file.Name()),
		ContentType: aws.String(content_type),
		Body:        file,
	})
	if err != nil {
		return fmt.Errorf("uploader.Upload error: %v", err)
	}
	fmt.Println("upload complete!")
	return nil
}

func downloadfromS3bucket(client *s3.Client) error {
	downloader := manager.NewDownloader(client)
	newBuff := manager.NewWriteAtBuffer([]byte{})
	numBytes, err := downloader.Download(context.TODO(), newBuff, &s3.GetObjectInput{
		Bucket: aws.String(BUCKETNAME),
		Key:    aws.String("images.jpeg"),
	})
	if err != nil {
		return fmt.Errorf("downloader.Download error: %v", err)
	}

	if numBytes != int64(len(newBuff.Bytes())) {
		return fmt.Errorf("Download error numBytes not equal to buff bytes: %d vs %d", numBytes, len(newBuff.Bytes()))
	}
	// err = os.WriteFile("downloaded_image.jpeg", newBuff.Bytes(), 0644)
	// if err != nil {
	// return fmt.Errorf("os.WriteFile error: %v", err)
	// }
	fmt.Printf("Download complete read %d bytes\n", numBytes)
	return nil
}

func InitClient(region string) (*s3.Client, error) {
	var (
		cnfg   aws.Config
		err    error
		client *s3.Client
	)

	if cnfg, err = config.LoadDefaultConfig(context.TODO(), config.WithRegion(region)); err != nil {
		return nil, fmt.Errorf("loaddefaultconfig error: %v", err)
	}

	client = s3.NewFromConfig(cnfg)
	return client, nil
}
