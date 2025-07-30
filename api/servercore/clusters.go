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

func (cl *ServercoreClient) CreateCluster(payload svcmodel.CreateClusterRequest) (*svcmodel.CreateClusterResponse, error) {
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

	createClusterResponse := svcmodel.CreateClusterResponse{}
	if err := json.Unmarshal(respBody, &createClusterResponse); err != nil {
		return nil, err
	}

	return &createClusterResponse, nil
}
