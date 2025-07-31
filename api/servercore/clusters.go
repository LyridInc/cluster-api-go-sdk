package servercore

import (
	"bytes"
	"encoding/json"
	"net/http"

	svcmodel "github.com/LyridInc/cluster-api-go-sdk/model/servercore"
)

func (cl *ServercoreClient) GetListClusters() (*svcmodel.ClusterListResponse, error) {
	url := cl.ManagedKubernetesApiUrl + "/clusters"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", cl.AuthToken)

	respBody, err := cl.doHttpRequest(req)
	if err != nil {
		return nil, err
	}

	clusterListResponse := svcmodel.ClusterListResponse{}
	if err := json.Unmarshal(respBody, &clusterListResponse); err != nil {
		return nil, err
	}

	return &clusterListResponse, nil
}

func (cl *ServercoreClient) GetClusterByID(clusterId string) (*svcmodel.ClusterResponse, error) {
	url := cl.ManagedKubernetesApiUrl + "/clusters/" + clusterId
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", cl.AuthToken)

	respBody, err := cl.doHttpRequest(req)
	if err != nil {
		return nil, err
	}

	clusterResponse := svcmodel.ClusterResponse{}
	if err := json.Unmarshal(respBody, &clusterResponse); err != nil {
		return nil, err
	}

	return &clusterResponse, nil
}

func (cl *ServercoreClient) CreateCluster(payload svcmodel.CreateClusterRequest) (*svcmodel.ClusterResponse, error) {
	url := cl.ManagedKubernetesApiUrl + "/clusters"
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", cl.AuthToken)

	respBody, err := cl.doHttpRequest(req)
	if err != nil {
		return nil, err
	}

	clusterResponse := svcmodel.ClusterResponse{}
	if err := json.Unmarshal(respBody, &clusterResponse); err != nil {
		return nil, err
	}

	return &clusterResponse, nil
}

func (cl *ServercoreClient) DeleteClusterByID(clusterId string) (*any, error) {
	url := cl.ManagedKubernetesApiUrl + "/clusters/" + clusterId
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", cl.AuthToken)

	_, err = cl.doHttpRequest(req)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (cl *ServercoreClient) GetClusterKubeconfig(clusterId string) (*string, error) {
	url := cl.ManagedKubernetesApiUrl + "/clusters/" + clusterId + "/kubeconfig"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", cl.AuthToken)

	resp, err := cl.doHttpRequest(req)
	if err != nil {
		return nil, err
	}

	kubeconfig := string(resp)
	return &kubeconfig, nil
}
