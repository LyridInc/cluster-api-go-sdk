package servercore

import (
	"bytes"
	"encoding/json"
	"net/http"

	svcmodel "github.com/LyridInc/cluster-api-go-sdk/model/servercore"
)

func (cl *ServercoreClient) GetListFloatingIPs() (*svcmodel.FloatingIPsResponse, error) {
	url := cl.ApiUrl + "/vpc/resell/v2/floatingips"
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

	floatingIPsResponse := svcmodel.FloatingIPsResponse{}
	if err := json.Unmarshal(respBody, &floatingIPsResponse); err != nil {
		return nil, err
	}

	return &floatingIPsResponse, nil
}

func (cl *ServercoreClient) GetListFloatingIPDetail(floatingIPId string) (*svcmodel.FloatingIPResponse, error) {
	url := cl.ApiUrl + "/vpc/resell/v2/floatingips/" + floatingIPId
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

	floatingIPResponse := svcmodel.FloatingIPResponse{}
	if err := json.Unmarshal(respBody, &floatingIPResponse); err != nil {
		return nil, err
	}

	return &floatingIPResponse, nil
}

func (cl *ServercoreClient) CreateFloatingIP(payload svcmodel.CreateFloatingIPRequest) (*svcmodel.FloatingIPsResponse, error) {
	url := cl.ApiUrl + "/vpc/resell/v2/floatingips/projects/" + cl.ProjectID
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Token", cl.ApiKey)

	respBody, err := cl.doHttpRequest(req)
	if err != nil {
		return nil, err
	}

	floatingIPsResponse := svcmodel.FloatingIPsResponse{}
	if err := json.Unmarshal(respBody, &floatingIPsResponse); err != nil {
		return nil, err
	}

	return &floatingIPsResponse, nil
}

func (cl *ServercoreClient) DeleteFloatingIPByID(floatingIPId string) (*any, error) {
	url := cl.ApiUrl + "/vpc/resell/v2/floatingips/" + floatingIPId
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Token", cl.ApiKey)

	_, err = cl.doHttpRequest(req)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
