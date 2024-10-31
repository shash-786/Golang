package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func main() {
	var (
		instanceID, region string
		err                error
	)

	flag.StringVar(&region, "r", "", "Enter a region")
	flag.Parse()

	if region == "" {
		region = "ap-south-1"
	}

	if instanceID, err = ec2Instance(region); err != nil {
		log.Fatalf("ec2 instance creation error : %v", err)
	}

	fmt.Printf("InstanceID: %v", instanceID)
}

func ec2Instance(region string) (string, error) {
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return "", fmt.Errorf("loadconfig error: %v", err)
	}

	ec2client := ec2.NewFromConfig(config)

	// NOTE: KEY CREATION LOGIC ONLY REQUIRED IF YOU WANT TO SSH
	keypair, err := ec2client.DescribeKeyPairs(context.TODO(), &ec2.DescribeKeyPairsInput{
		KeyNames: []string{"go-aws-demo"},
	})
	if err != nil && !strings.Contains(err.Error(), "InvalidKeyPair.NotFound") {
		return "", fmt.Errorf("DescribeKeyPair error: %v", err)
	}

	if keypair == nil || len(keypair.KeyPairs) == 0 {
		new_key_pair, err := ec2client.CreateKeyPair(context.TODO(), &ec2.CreateKeyPairInput{
			KeyName: aws.String("go-aws-demo"),
		})
		if err != nil {
			return "", fmt.Errorf("CreateKeyPair error: %v", err)
		}

		if err := os.WriteFile("go-aws-demo.pem", []byte(*new_key_pair.KeyMaterial), 06000); err != nil {
			return "", fmt.Errorf("writefile error: %v", err)
		}
	} else {
		fmt.Println("Key already exists!")
	}

	image_out, err := ec2client.DescribeImages(context.TODO(), &ec2.DescribeImagesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("image-id"),
				Values: []string{"ami-09b0a86a2c84101e1"},
			},
			{
				Name:   aws.String("virtualization-type"),
				Values: []string{"hvm"},
			},
		},

		Owners: []string{"099720109477"},
	})
	if err != nil {
		return "", fmt.Errorf("DescribeImages error: %v", err)
	}

	output, err := ec2client.RunInstances(context.TODO(), &ec2.RunInstancesInput{
		InstanceType: types.InstanceTypeT3Micro,
		MaxCount:     aws.Int32(1),
		MinCount:     aws.Int32(1),
		ImageId:      image_out.Images[0].ImageId,
		KeyName:      aws.String("go-aws-demo"),
	})
	if err != nil {
		return "", fmt.Errorf("Run Instances error %v: ", err)
	}

	if len(output.Instances) == 0 {
		return "", fmt.Errorf("no instances found")
	}

	return *output.Instances[0].InstanceId, nil
}
