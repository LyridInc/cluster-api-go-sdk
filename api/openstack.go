package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/LyridInc/cluster-api-go-sdk/model"
)

type (
	OpenstackClient struct {
		NetworkEndpoint string
		AuthEndpoint    string
		AuthToken       string
		ProjectId       string
	}

	OpenstackCredential struct {
		ApplicationCredentialName   string
		ApplicationCredentialId     string
		ApplicationCredentialSecret string
	}
)

func (c *OpenstackClient) Authenticate(credential OpenstackCredential) error {
	token := os.Getenv("OS_TOKEN")
	if token != "" {
		c.AuthToken = token
		return nil
	}

	url := c.AuthEndpoint + "/v3/auth/tokens"
	requestBody := []byte(`{
		"auth": {
			"identity": {
				"methods": ["application_credential"],
				"application_credential": {
					"id": "` + credential.ApplicationCredentialId + `",
					"name": "` + credential.ApplicationCredentialName + `",
					"secret": "` + credential.ApplicationCredentialSecret + `"
				}
			}
		}
	}`)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	for key, value := range response.Header {
		if key == "X-Subject-Token" && len(value) > 0 {
			c.AuthToken = value[0]
			os.Setenv("OS_TOKEN", c.AuthToken)
			break
		}
	}

	return nil
}

func (c *OpenstackClient) GetProjectQuotas() (*model.QuotasResponse, error) {
	url := c.NetworkEndpoint + "/v2.0/quotas/" + c.ProjectId + "/details.json"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Token", c.AuthToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	quotas := model.QuotasResponse{}
	json.Unmarshal(body, &quotas)

	return &quotas, nil
}
