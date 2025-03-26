package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/LyridInc/cluster-api-go-sdk/model"
)

type ICivoClient interface {
	ListClusters(queryParams map[string]string) ([]byte, error)
	CreateCluster(args CivoCreateClusterArgs) ([]byte, error)
	GetCluster(id string) ([]byte, error)
	GetClusterDetail(id string) ([]byte, error)
	DeleteCluster(id string) ([]byte, error)

	ListNetworks(queryParams map[string]string) ([]model.CivoNetworkResponse, error)
	ListInstanceSizes(queryParams map[string]string) ([]model.CivoInstanceSizeResponse, error)
	ListKubernetesVersions(queryParams map[string]string) ([]model.CivoKubernetesVersionResponse, error)
	ListFirewalls(queryParams map[string]string) ([]model.CivoFirewallResponse, error)
	ListFirewallRules(id string) ([]model.CivoFirewallRuleResponse, error)
	ListMarketplaceApplications(queryParams map[string]string) ([]model.CivoMarketplaceItemResponse, error)
}

type CivoClient struct {
	APIToken    string
	APIEndpoint string
}

type CivoCreateClusterArgs struct {
	Name              string           `json:"name,omitempty"`
	NetworkID         string           `json:"network_id,omitempty"`
	Region            string           `json:"region,omitempty"`
	CNIPlugin         string           `json:"cni_plugin,omitempty"`
	Pools             []model.CivoPool `json:"pools,omitempty"`
	KubernetesVersion string           `json:"kubernetes_version,omitempty"`
	Tags              string           `json:"tags,omitempty"`
	InstanceFirewall  string           `json:"instance_firewall,omitempty"`
	FirewallRule      string           `json:"firewall_rule,omitempty"`
	Applications      string           `json:"applications,omitempty"`
}

func NewCivoClient(token, endpoint string) ICivoClient {
	return &CivoClient{
		APIToken:    token,
		APIEndpoint: endpoint,
	}
}

func (cl *CivoClient) doHttpRequest(request *http.Request) ([]byte, error) {
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

func (cl *CivoClient) ListClusters(queryParams map[string]string) ([]byte, error) {
	baseURL := cl.APIEndpoint + "/v2/kubernetes/clusters"

	if len(queryParams) > 0 {
		q := url.Values{}
		for key, value := range queryParams {
			q.Add(key, value)
		}
		baseURL += "?" + q.Encode()
	}

	request, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+cl.APIToken)

	return cl.doHttpRequest(request)
}

func (cl *CivoClient) CreateCluster(args CivoCreateClusterArgs) ([]byte, error) {
	argsByte, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", cl.APIEndpoint+"/v2/kubernetes/clusters", bytes.NewBuffer(argsByte))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+cl.APIToken)

	return cl.doHttpRequest(request)
}

func (cl *CivoClient) GetCluster(id string) ([]byte, error) {
	request, err := http.NewRequest("GET", cl.APIEndpoint+"/v2/kubernetes/clusters/"+id, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+cl.APIToken)

	return cl.doHttpRequest(request)
}

func (cl *CivoClient) GetClusterDetail(id string) ([]byte, error) {
	request, err := http.NewRequest("GET", cl.APIEndpoint+"/v2/kubernetes/clusters/"+id, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+cl.APIToken)

	return cl.doHttpRequest(request)
}

func (cl *CivoClient) DeleteCluster(id string) ([]byte, error) {
	request, err := http.NewRequest("DELETE", cl.APIEndpoint+"/v2/kubernetes/cluster/"+id, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+cl.APIToken)

	return cl.doHttpRequest(request)
}

func (cl *CivoClient) ListFirewallRules(id string) ([]model.CivoFirewallRuleResponse, error) {
	baseURL := cl.APIEndpoint + "/v2/firewalls/" + id + "/rules"

	request, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+cl.APIToken)

	b, err := cl.doHttpRequest(request)
	if err != nil {
		return nil, err
	}

	var resp []model.CivoFirewallRuleResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (cl *CivoClient) ListFirewalls(queryParams map[string]string) ([]model.CivoFirewallResponse, error) {
	baseURL := cl.APIEndpoint + "/v2/firewalls"

	if len(queryParams) > 0 {
		q := url.Values{}
		for key, value := range queryParams {
			q.Add(key, value)
		}
		baseURL += "?" + q.Encode()
	}

	request, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+cl.APIToken)

	b, err := cl.doHttpRequest(request)
	if err != nil {
		return nil, err
	}

	var resp []model.CivoFirewallResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (cl *CivoClient) ListInstanceSizes(queryParams map[string]string) ([]model.CivoInstanceSizeResponse, error) {
	baseURL := cl.APIEndpoint + "/v2/sizes"

	if len(queryParams) > 0 {
		q := url.Values{}
		for key, value := range queryParams {
			q.Add(key, value)
		}
		baseURL += "?" + q.Encode()
	}

	request, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+cl.APIToken)

	b, err := cl.doHttpRequest(request)
	if err != nil {
		return nil, err
	}

	var resp []model.CivoInstanceSizeResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (cl *CivoClient) ListKubernetesVersions(queryParams map[string]string) ([]model.CivoKubernetesVersionResponse, error) {
	baseURL := cl.APIEndpoint + "/v2/kubernetes/versions"

	if len(queryParams) > 0 {
		q := url.Values{}
		for key, value := range queryParams {
			q.Add(key, value)
		}
		baseURL += "?" + q.Encode()
	}

	request, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+cl.APIToken)

	b, err := cl.doHttpRequest(request)
	if err != nil {
		return nil, err
	}

	var resp []model.CivoKubernetesVersionResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (cl *CivoClient) ListMarketplaceApplications(queryParams map[string]string) ([]model.CivoMarketplaceItemResponse, error) {
	baseURL := cl.APIEndpoint + "/v2/kubernetes/applications"

	if len(queryParams) > 0 {
		q := url.Values{}
		for key, value := range queryParams {
			q.Add(key, value)
		}
		baseURL += "?" + q.Encode()
	}

	request, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+cl.APIToken)

	b, err := cl.doHttpRequest(request)
	if err != nil {
		return nil, err
	}

	var resp []model.CivoMarketplaceItemResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (cl *CivoClient) ListNetworks(queryParams map[string]string) ([]model.CivoNetworkResponse, error) {
	baseURL := cl.APIEndpoint + "/v2/networks"

	if len(queryParams) > 0 {
		q := url.Values{}
		for key, value := range queryParams {
			q.Add(key, value)
		}
		baseURL += "?" + q.Encode()
	}

	request, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+cl.APIToken)

	b, err := cl.doHttpRequest(request)
	if err != nil {
		return nil, err
	}

	var resp []model.CivoNetworkResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
