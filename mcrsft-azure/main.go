package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/shash-786/Golang/ssh"
)

var (
	LOCATION             = "westus"
	RESOURCEGROUPNAME    = "go-azure-demo"
	VIRTUALNETWORKNAME   = "sample-virtual-network"
	SUBTNETNAME          = "sample-subnet-network"
	PUBLICIPADDRESSNAME  = "sample-public-ip-address-allocation"
	SECURITYGROUPNAME    = "shash-security"
	NETWORKINTERFACENAME = "salian-iiface"
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

	subnet_resp, err := subnetpollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}

	Public_IP_Address_Client, err := armnetwork.NewPublicIPAddressesClient(subscriptionID, cred, nil)
	if err != nil {
		return err
	}

	Public_IP_Address_Poller_Resp, err := Public_IP_Address_Client.BeginCreateOrUpdate(
		ctx,
		*ResourceGroupsClientCreateOrUpdateResponse.Name,
		PUBLICIPADDRESSNAME,
		armnetwork.PublicIPAddress{
			Location: &LOCATION,
			Properties: &armnetwork.PublicIPAddressPropertiesFormat{
				PublicIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodStatic),
			},
		},
		nil)
	if err != nil {
		return err
	}
	Public_IP_Address_Resp, err := Public_IP_Address_Poller_Resp.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}

	Security_Groups_client, err := armnetwork.NewSecurityGroupsClient(subscriptionID, cred, nil)
	if err != nil {
		return err
	}

	Security_Groups_Poller_Resp, err := Security_Groups_client.BeginCreateOrUpdate(
		ctx,
		*ResourceGroupsClientCreateOrUpdateResponse.Name,
		SECURITYGROUPNAME,
		armnetwork.SecurityGroup{
			Location: &LOCATION,
			Properties: &armnetwork.SecurityGroupPropertiesFormat{
				SecurityRules: []*armnetwork.SecurityRule{
					{
						Name: to.Ptr("allow-ssh"),
						Properties: &armnetwork.SecurityRulePropertiesFormat{
							SourceAddressPrefix:      to.Ptr("0.0.0.0/0"),
							SourcePortRange:          to.Ptr("*"),
							DestinationAddressPrefix: to.Ptr("0.0.0.0/0"),
							DestinationPortRange:     to.Ptr("22"),
							Protocol:                 to.Ptr(armnetwork.SecurityRuleProtocolTCP),
							Access:                   to.Ptr(armnetwork.SecurityRuleAccessAllow),
							Description:              to.Ptr("Allow ssh on port 22"),
							Direction:                to.Ptr(armnetwork.SecurityRuleDirectionInbound),
							Priority:                 to.Ptr(int32(1001)),
						},
					},
				},
			},
		},
		nil,
	)
	if err != nil {
		return err
	}

	Security_Groups_Resp, err := Security_Groups_Poller_Resp.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}

	Interface_client, err := armnetwork.NewInterfacesClient(subscriptionID, cred, nil)
	if err != nil {
		return err
	}

	Interface_Poller_Resp, err := Interface_client.BeginCreateOrUpdate(
		ctx,
		*ResourceGroupsClientCreateOrUpdateResponse.Name,
		NETWORKINTERFACENAME,
		armnetwork.Interface{
			Location: &LOCATION,
			Properties: &armnetwork.InterfacePropertiesFormat{
				NetworkSecurityGroup: &armnetwork.SecurityGroup{
					ID: Security_Groups_Resp.ID,
				},
				IPConfigurations: []*armnetwork.InterfaceIPConfiguration{
					{
						Name: to.Ptr("ipConfig"),
						Properties: &armnetwork.InterfaceIPConfigurationPropertiesFormat{
							PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
							Subnet: &armnetwork.Subnet{
								ID: subnet_resp.ID,
							},
							PublicIPAddress: &armnetwork.PublicIPAddress{
								ID: Public_IP_Address_Resp.ID,
							},
						},
					},
				},
			},
		},
		nil)
	if err != nil {
		return err
	}

	Interface_Resp, err := Interface_Poller_Resp.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}

	fmt.Println("Creating VM")

	virtualMachinesClient, err := armcompute.NewVirtualMachinesClient(
		subscriptionID,
		cred,
		nil,
	)

	parameters := armcompute.VirtualMachine{
		Location: to.Ptr(LOCATION),
		Identity: &armcompute.VirtualMachineIdentity{
			Type: to.Ptr(armcompute.ResourceIdentityTypeNone),
		},
		Properties: &armcompute.VirtualMachineProperties{
			StorageProfile: &armcompute.StorageProfile{
				ImageReference: &armcompute.ImageReference{
					// search image reference
					// az vm image list --output table
					// require ssh key for authentication on linux
					Offer:     to.Ptr("0001-com-ubuntu-server-focal"),
					Publisher: to.Ptr("canonical"),
					SKU:       to.Ptr("20_04-lts-gen2"),
					Version:   to.Ptr("latest"),
				},
				OSDisk: &armcompute.OSDisk{
					Name:         to.Ptr("TEMP"),
					CreateOption: to.Ptr(armcompute.DiskCreateOptionTypesFromImage),
					Caching:      to.Ptr(armcompute.CachingTypesReadWrite),
					ManagedDisk: &armcompute.ManagedDiskParameters{
						StorageAccountType: to.Ptr(armcompute.StorageAccountTypesStandardLRS), // OSDisk type Standard/Premium HDD/SSD
					},
					// DiskSizeGB: to.Ptr[int32](100), // default 127G
				},
			},
			HardwareProfile: &armcompute.HardwareProfile{
				VMSize: to.Ptr(armcompute.VirtualMachineSizeTypes("Standard_B1s")), // VM size include vCPUs,RAM,Data Disks,Temp storage.
			},
			OSProfile: &armcompute.OSProfile{ //
				ComputerName:  to.Ptr("sample-computer"),
				AdminUsername: to.Ptr("sample-user"),
				// require ssh key for authentication on linux
				LinuxConfiguration: &armcompute.LinuxConfiguration{
					DisablePasswordAuthentication: to.Ptr(true),
					SSH: &armcompute.SSHConfiguration{
						PublicKeys: []*armcompute.SSHPublicKey{
							{
								Path:    to.Ptr(fmt.Sprintf("/home/%s/.ssh/authorized_keys", "sample-user")),
								KeyData: to.Ptr(publicKey),
							},
						},
					},
				},
			},
			NetworkProfile: &armcompute.NetworkProfile{
				NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
					{
						ID: Interface_Resp.ID,
					},
				},
			},
		},
	}

	pollerResponse, err := virtualMachinesClient.BeginCreateOrUpdate(
		ctx,
		*ResourceGroupsClientCreateOrUpdateResponse.Name,
		"demo-virtual-machine",
		parameters,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	fmt.Println("VM Successfully Created")

	return nil
}
