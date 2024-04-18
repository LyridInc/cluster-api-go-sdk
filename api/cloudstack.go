package api

import (
	"github.com/apache/cloudstack-go/cloudstack"
)

type ICloudStackClient interface {
	ListPublicIpAddresses(zoneID, state string) (*cloudstack.ListPublicIpAddressesResponse, error)
	CreateSSHKeypair(keypairName string) (*cloudstack.CreateSSHKeyPairResponse, error)
	GetZoneByID(zoneID string) (*cloudstack.Zone, error)
}

type CloudStackClient struct {
	Client *cloudstack.CloudStackClient
}

func NewCloudStackClient(apiURL, apiKey, secret string, verifySSL bool) ICloudStackClient {
	cs := cloudstack.NewAsyncClient(apiURL, apiKey, secret, verifySSL)
	return &CloudStackClient{
		Client: cs,
	}
}

func (cl *CloudStackClient) ListPublicIpAddresses(zoneID, state string) (*cloudstack.ListPublicIpAddressesResponse, error) {
	params := cloudstack.ListPublicIpAddressesParams{}
	params.SetZoneid(zoneID)
	params.SetState(state)

	return cl.Client.Address.ListPublicIpAddresses(&params)
}

func (cl *CloudStackClient) CreateSSHKeypair(keypairName string) (*cloudstack.CreateSSHKeyPairResponse, error) {
	params := cloudstack.CreateSSHKeyPairParams{}
	params.SetName(keypairName)
	return cl.Client.SSH.CreateSSHKeyPair(&params)
}

func (cl *CloudStackClient) GetZoneByID(zoneID string) (*cloudstack.Zone, error) {
	zone, _, err := cl.Client.Zone.GetZoneByID(zoneID)
	return zone, err
}
