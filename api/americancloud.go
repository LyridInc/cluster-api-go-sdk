package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// https://app.americancloud.com/api/documentation
type IAmericanCloudClient interface {
	ListClusters() ([]byte, error)
	CreateCluster(args AmericanCloudCreateClusterArgs) ([]byte, error)
	GetCluster(clusterName string) ([]byte, error)
	GetClusterKubeconfig(clusterName string) ([]byte, error)
	DeleteCluster(clusterName string) ([]byte, error)

	GetNetwork(networkName string) ([]byte, error)
	DeleteNetwork(networkName string) ([]byte, error)
}

type AmericanCloudClient struct {
	ProjectID   string
	APIToken    string
	APIEndpoint string
}

type AmericanCloudCreateClusterArgs struct {
	Name             string `json:"name,omitempty"`
	Project          string `json:"project,omitempty"`
	Zone             string `json:"zone,omitempty"`
	Version          string `json:"version,omitempty"`
	NodeSize         int    `json:"node_size,omitempty"`
	Package          string `json:"package,omitempty"`
	Autoscale        int    `json:"autoscale,omitempty"`
	MinClusterSize   int    `json:"min_cluster_size,omitempty"`
	MaxClusterSize   int    `json:"max_cluster_size,omitempty"`
	BillingPeriod    string `json:"billing_period,omitempty"`
	KeyPair          string `json:"key_pair,omitempty"`
	ControlNodes     int    `json:"control_nodes,omitempty"`
	HighAvailability bool   `json:"high_availability,omitempty"`
	CouponCode       string `json:"coupon_codes,omitempty"`
}

func NewAmericanCloudClient(projectID, token, endpoint string) IAmericanCloudClient {
	return &AmericanCloudClient{
		ProjectID:   projectID,
		APIToken:    token,
		APIEndpoint: endpoint,
	}
}

func (cl *AmericanCloudClient) doHttpRequest(request *http.Request) ([]byte, error) {
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("error %d: resource not found for %s", response.StatusCode, request.URL)
	}

	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)
	return b, err
}

func (cl *AmericanCloudClient) ListClusters() ([]byte, error) {
	request, err := http.NewRequest("GET", cl.APIEndpoint+"/clusters/"+cl.ProjectID, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+cl.APIToken)

	return cl.doHttpRequest(request)
}

func (cl *AmericanCloudClient) CreateCluster(args AmericanCloudCreateClusterArgs) ([]byte, error) {
	argsByte, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", cl.APIEndpoint+"/cluster/create", bytes.NewBuffer(argsByte))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+cl.APIToken)

	return cl.doHttpRequest(request)
}

func (cl *AmericanCloudClient) GetCluster(clusterName string) ([]byte, error) {
	request, err := http.NewRequest("GET", cl.APIEndpoint+"/cluster/"+clusterName, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+cl.APIToken)

	return cl.doHttpRequest(request)
}

func (cl *AmericanCloudClient) GetClusterKubeconfig(clusterName string) ([]byte, error) {
	request, err := http.NewRequest("GET", cl.APIEndpoint+"/cluster/config/"+clusterName, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+cl.APIToken)

	return cl.doHttpRequest(request)
}

func (cl *AmericanCloudClient) DeleteCluster(clusterName string) ([]byte, error) {
	request, err := http.NewRequest("DELETE", cl.APIEndpoint+"/cluster/"+clusterName, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+cl.APIToken)

	return cl.doHttpRequest(request)
}

func (cl *AmericanCloudClient) GetNetwork(networkName string) ([]byte, error) {
	request, err := http.NewRequest("GET", cl.APIEndpoint+"/network/"+networkName, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+cl.APIToken)

	return cl.doHttpRequest(request)
}

func (cl *AmericanCloudClient) DeleteNetwork(networkName string) ([]byte, error) {
	request, err := http.NewRequest("DELETE", cl.APIEndpoint+"/network/"+networkName, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+cl.APIToken)

	return cl.doHttpRequest(request)
}
