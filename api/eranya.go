package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// https://jkt-2.console.eranyacloud.com:8800/api/v2/login

type (
	EranyaCloud interface {
		Authenticate(domain, username, password string) ([]byte, error)
	}

	EranyaCloudImpl struct {
		BaseURL   string `json:"base_url"`
		AuthToken string `json:"auth_token"`
	}
)

func NewEranyaCloud(baseURL string) EranyaCloud {
	return &EranyaCloudImpl{
		BaseURL: baseURL,
	}
}

func (v *EranyaCloudImpl) Authenticate(domain, username, password string) ([]byte, error) {
	url := v.BaseURL + "/v2/login"
	args := map[string]interface{}{
		"domain":                      domain,
		"username":                    username,
		"password":                    password,
		"domainAdminStartPageEnabled": false,
	}
	argsByte, _ := json.Marshal(args)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(argsByte))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Requested-With", "XMLHttpRequest")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	b, _ := io.ReadAll(response.Body)
	// {"id": "c3356974b0434053929db5695e782e16", "name": "a-lyrid", "is_group": false, "is_superuser": false, "is_enabled": true, "description": "", "external_id": null, "external_provider": null, "roles": [{"id": "login", "name": "Login", "description": "Can login in web UI", "tags": []}], "email": null, "domain_id": "7c503dbcc85b494d91834e2b0494a05e", "project_id": null, "token": "unscoped", "scoped_token": null}

	return b, nil
}
