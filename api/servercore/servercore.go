package servercore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	svcmodel "github.com/LyridInc/cluster-api-go-sdk/model/servercore"
)

type ServercoreClient struct {
	ApiKey                  string
	CloudApiUrl             string
	ManagedKubernetesApiUrl string
	AuthToken               string
}

func NewServercoreClient(config svcmodel.Config) *ServercoreClient {
	return &ServercoreClient{
		ApiKey:                  config.ApiKey,
		CloudApiUrl:             config.CloudApiUrl,
		ManagedKubernetesApiUrl: config.ManagedKubernetesApiUrl,
	}
}

func (cl *ServercoreClient) doHttpRequest(request *http.Request) ([]byte, error) {
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var apiErr struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.Unmarshal(b, &apiErr); err != nil {
			return nil, fmt.Errorf("API error (status %d): %s", response.StatusCode, string(b))
		}
		return nil, fmt.Errorf("API error: %s", apiErr.Error.Message)
	}
	return b, err
}

func (cl *ServercoreClient) Authenticate(authConfig svcmodel.AuthConfig) (*svcmodel.AuthResponse, string, error) {
	authPayload := map[string]any{
		"auth": map[string]any{
			"identity": map[string]any{
				"methods": []string{"password"},
				"password": map[string]any{
					"user": map[string]any{
						"name":     authConfig.Username,
						"domain":   map[string]string{"name": authConfig.Domain},
						"password": authConfig.Password,
					},
				},
			},
			"scope": map[string]any{
				"project": map[string]any{
					"name":   authConfig.ProjectName,
					"domain": map[string]string{"name": authConfig.Domain},
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

	authResponse := svcmodel.AuthResponse{}
	if err := json.Unmarshal(body, &authResponse); err != nil {
		return nil, "", err
	}

	token := resp.Header.Get("X-Subject-Token")
	cl.AuthToken = token

	return &authResponse, token, nil
}
