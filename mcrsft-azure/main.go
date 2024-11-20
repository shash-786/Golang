package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/shash-786/Golang/ssh"
)

var (
	LOCATION           = "westus"
	RESOURCEGROUPNAME  = "go-azure-demo"
	VIRTUALNETWORKNAME = "sample-virtual-network"
	SUBTNETNAME        = "sample-subnet-network"
)

func main() {
	var (
		publicKey, subscriptionID string
		token                     *azidentity.AzureCLICredential
		err                       error
	)

	if publicKey, err = generatekeys(); err != nil {
		fmt.Printf("generatekeys error: %v", err)
		os.Exit(1)
	}

	ctx := context.Background()
	if token, err = generateToken(); err != nil {
		fmt.Printf("generatekeys error: %v", err)
		os.Exit(1)
	}

	subscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	if subscriptionID == "" {
		fmt.Println("No SubscriptionID found!")
		os.Exit(1)
	}

	if err = launchInstance(ctx, token, subscriptionID, publicKey); err != nil {
		fmt.Printf("launchInstance error: %v", err)
		os.Exit(1)
	}
}

func generatekeys() (string, error) {
	var priv, pub []byte
	var err error

	if priv, pub, err = ssh.GenerateKeys(); err != nil {
		fmt.Printf("GenerateKeys error: %v ", err)
		os.Exit(1)
	}

	if err = os.WriteFile("./keys/myKey.pem", priv, 0600); err != nil {
		fmt.Printf("priv writefile error: %v", err)
		os.Exit(1)
	}

	if err = os.WriteFile("./keys/myKey.pub", pub, 0644); err != nil {
		fmt.Printf("pub writefile error: %v", err)
		os.Exit(1)
	}

	return string(pub), err
}

func generateToken() (*azidentity.AzureCLICredential, error) {
	token, err := azidentity.NewAzureCLICredential(nil)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func launchInstance(ctx context.Context, cred *azidentity.AzureCLICredential, subscriptionID, publicKey string) error {
	client, err := armresources.NewResourceGroupsClient(subscriptionID, cred, nil)
	if err != nil {
		return fmt.Errorf("armresources.NewClient error: %v", err)
	}

	resource_grp_params := armresources.ResourceGroup{
		Location: to.Ptr(LOCATION),
	}

	ResourceGroupsClientCreateOrUpdateResponse, err := client.CreateOrUpdate(ctx, RESOURCEGROUPNAME, resource_grp_params, nil)
	if err != nil {
		return fmt.Errorf("armresources.client error: %v", err)
	}

	networkClientFactory, err := armnetwork.NewClientFactory(subscriptionID, cred, nil)
	if err != nil {
		log.Fatal(err)
	}
	virtualNetworksClient := networkClientFactory.NewVirtualNetworksClient()

	vnetpollerResp, err := virtualNetworksClient.BeginCreateOrUpdate(
		ctx,
		*ResourceGroupsClientCreateOrUpdateResponse.Name,
		VIRTUALNETWORKNAME,
		armnetwork.VirtualNetwork{
			Location: to.Ptr(LOCATION),
			Properties: &armnetwork.VirtualNetworkPropertiesFormat{
				AddressSpace: &armnetwork.AddressSpace{
					AddressPrefixes: []*string{
						to.Ptr("10.1.0.0/16"),
					},
				},
			},
		},
		nil)
	if err != nil {
		return err
	}

	vnet_response, err := vnetpollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}

	subnetsClient := networkClientFactory.NewSubnetsClient()

	subnetpollerResp, err := subnetsClient.BeginCreateOrUpdate(
		ctx,
		*ResourceGroupsClientCreateOrUpdateResponse.Name,
		*vnet_response.Name,
		SUBTNETNAME,
		armnetwork.Subnet{
			Properties: &armnetwork.SubnetPropertiesFormat{
				AddressPrefix: to.Ptr("10.1.0.0/24"),
			},
		},
		nil,
	)
	if err != nil {
		return err
	}

	_, err = subnetpollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}
