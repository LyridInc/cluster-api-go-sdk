package servercore

import (
	"encoding/json"
	"net/http"

	svcmodel "github.com/LyridInc/cluster-api-go-sdk/model/servercore"
)

func (cl *ServercoreClient) GetListIPs() (*[]svcmodel.IP, error) {
	url := cl.ApiUrl + "/ipam/v1"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Token", cl.ApiKey)

	respBody, err := cl.doHttpRequest(req)
	if err != nil {
		return nil, err
	}

	ipListResponse := []svcmodel.IP{}
	if err := json.Unmarshal(respBody, &ipListResponse); err != nil {
		return nil, err
	}

	return &ipListResponse, nil
}
