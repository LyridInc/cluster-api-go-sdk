package servercore

import (
	"encoding/json"
	"io"
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

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	clusterListResponse := svcmodel.ClusterListResponse{}
	if err := json.Unmarshal(body, &clusterListResponse); err != nil {
		return nil, err
	}

	return &clusterListResponse, nil
}
