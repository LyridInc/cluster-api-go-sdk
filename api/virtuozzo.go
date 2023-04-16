package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/LyridInc/cluster-api-go-sdk/model"
)

type (
	Virtuozzo interface {
		Authenticate(project, username, password, domainId string) ([]byte, error)
		ListPublicEndpoints() ([]byte, error)
		ListImages() error
		ListKubernetesClusterTemplates() ([]byte, error)
		CreateNetwork() error
		CreateSubnet() error
		CreateLoadBalancer() error
		CreateKubernetesClusterTemplate(args model.VirtuozzoCreateKubernetesClusterTemplateArgs) ([]byte, error)
	}

	virtuozzo struct {
		Auth         model.VirtuozzoAuth `json:"auth"`
		AuthEndpoint string              `json:"auth_url"`
		AuthToken    string              `json:"token"`
	}
)

func NewVirtuozzoClient(authEndpoint string) Virtuozzo {
	return &virtuozzo{
		AuthEndpoint: authEndpoint,
	}
}

func (v *virtuozzo) Authenticate(project, username, password, domainId string) ([]byte, error) {
	body := map[string]interface{}{
		"auth": model.VirtuozzoAuth{
			Identity: model.VirtuozzoIdentity{
				Methods: []string{"password"},
				Password: model.IdentityPassword{
					User: model.UserIdentity{
						Name: username,
						Domain: map[string]interface{}{
							"id": domainId,
						},
						Password: password,
					},
				},
			},
			Scope: model.VirtuozzoScope{
				Project: model.ScopeProject{
					Name: project,
					Domain: map[string]interface{}{
						"id": domainId,
					},
				},
			},
		},
	}
	b, _ := json.Marshal(body)

	request, err := http.NewRequest("POST", v.AuthEndpoint+"/v3/auth/tokens", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	rb, _ := io.ReadAll(response.Body)

	for key, value := range response.Header {
		if key == "X-Subject-Token" && len(value) > 0 {
			v.AuthToken = value[0]
			break
		}
	}

	return rb, nil
}

func (v *virtuozzo) ListPublicEndpoints() ([]byte, error) {
	url := v.AuthEndpoint + "/v3/auth/catalog"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Token", v.AuthToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	b, _ := io.ReadAll(response.Body)

	return b, nil
}

func (v *virtuozzo) ListImages() error {
	// TODO: https://docs.virtuozzo.com/virtuozzo_hybrid_infrastructure_5_4_compute_api_reference/index.html#listing-images.html
	return nil
}

func (v *virtuozzo) ListKubernetesClusterTemplates() ([]byte, error) {
	url := "https://jkt-2.console.eranyacloud.com:9513/v1/clustertemplates"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Token", v.AuthToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	b, _ := io.ReadAll(response.Body)

	return b, nil
}

func (v *virtuozzo) CreateNetwork() error {
	return nil
}

func (v *virtuozzo) CreateSubnet() error {
	return nil
}

func (v *virtuozzo) CreateLoadBalancer() error {
	return nil
}

func (v *virtuozzo) CreateKubernetesClusterTemplate(args model.VirtuozzoCreateKubernetesClusterTemplateArgs) ([]byte, error) {
	url := "https://jkt-2.console.eranyacloud.com:9513/v1/clustertemplates"
	argsByte, _ := json.Marshal(args)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(argsByte))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("OpenStack-API-Version", "container-infra 1.8")
	request.Header.Set("X-Auth-Token", v.AuthToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	b, _ := io.ReadAll(response.Body)

	return b, nil
}
