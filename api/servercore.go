package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/LyridInc/cluster-api-go-sdk/model/servercore"
)

type ServercoreClient struct {
	ApiKey      string
	ApiUrl      string
	CloudApiUrl string
}

func NewServercoreClient(config servercore.Config) *ServercoreClient {
	return &ServercoreClient{
		ApiKey:      config.ApiKey,
		ApiUrl:      config.ApiUrl,
		CloudApiUrl: config.CloudApiUrl,
	}
}

func (cl *ServercoreClient) doHttpRequest(request *http.Request) ([]byte, error) {
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

func (cl *ServercoreClient) Authenticate() (*servercore.AuthResponse, string, error) {
	authPayload := map[string]any{
		"auth": map[string]any{
			"identity": map[string]any{
				"methods": []string{"password"},
				"password": map[string]any{
					"user": map[string]any{
						"name":     "lyrid-integration",
						"domain":   map[string]string{"name": "447872"},
						"password": `-L#^B~7k,."|Z/0?]hj?`,
					},
				},
			},
			"scope": map[string]any{
				"project": map[string]any{
					"name":   "My First Project",
					"domain": map[string]string{"name": "447872"},
				},
			},
		},
	}

	payloadBytes, err := json.Marshal(authPayload)
	if err != nil {
		return nil, "", err
	}

	url := cl.CloudApiUrl + "/identity/v3/auth/tokens"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	authResponse := servercore.AuthResponse{}
	if err := json.Unmarshal(body, &authResponse); err != nil {
		return nil, "", err
	}

	token := resp.Header.Get("X-Subject-Token")

	return &authResponse, token, nil
}
