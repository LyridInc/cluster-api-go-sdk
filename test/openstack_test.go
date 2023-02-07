package test

import (
	"os"
	"testing"

	"github.com/LyridInc/cluster-api-go-sdk/api"
	"github.com/LyridInc/cluster-api-go-sdk/option"
)

// export $(< test.env)

// go test ./test -v -run ^TestClientAuthentication$
func TestClientAuthentication(t *testing.T) {
	cl := api.OpenstackClient{
		NetworkEndpoint: os.Getenv("OS_NETWORK_ENDPOINT"),
		AuthEndpoint:    os.Getenv("OS_AUTH_ENDPOINT"),
		AuthToken:       os.Getenv("OS_TOKEN"),
		ProjectId:       os.Getenv("OS_PROJECT_ID"),
	}
	cl.Authenticate(api.OpenstackCredential{
		ApplicationCredentialName:   os.Getenv("OS_APPLICATION_CREDENTIAL_NAME"),
		ApplicationCredentialId:     os.Getenv("OS_APPLICATION_CREDENTIAL_ID"),
		ApplicationCredentialSecret: os.Getenv("OS_APPLICATION_CREDENTIAL_SECRET"),
	})
}

// go test ./test -v -run ^TestCheckAuthToken$
func TestCheckAuthToken(t *testing.T) {
	cl := api.OpenstackClient{
		NetworkEndpoint: os.Getenv("OS_NETWORK_ENDPOINT"),
		AuthEndpoint:    os.Getenv("OS_AUTH_ENDPOINT"),
		AuthToken:       os.Getenv("OS_TOKEN"),
		ProjectId:       os.Getenv("OS_PROJECT_ID"),
	}
	t.Run("check auth token", func(t *testing.T) {
		t.Log("TOKEN:", cl.AuthToken)
		response, err := cl.CheckAuthToken()
		if err != nil {
			t.Fatal(error.Error(err))
		}
		t.Log(response)
	})
}

// go test ./test -v -run ^TestQuota$
func TestQuota(t *testing.T) {
	cl := api.OpenstackClient{
		NetworkEndpoint: os.Getenv("OS_NETWORK_ENDPOINT"),
		AuthEndpoint:    os.Getenv("OS_AUTH_ENDPOINT"),
		AuthToken:       os.Getenv("OS_TOKEN"),
		ProjectId:       os.Getenv("OS_PROJECT_ID"),
	}
	t.Run("show quota", func(t *testing.T) {
		cl.GetProjectQuotas()
	})
}

// go test ./test -v -run ^TestUpdateYamlManifestFlannel$
func TestUpdateYamlManifestFlannel(t *testing.T) {
	cl := api.OpenstackClient{
		NetworkEndpoint: os.Getenv("OS_NETWORK_ENDPOINT"),
		AuthEndpoint:    os.Getenv("OS_AUTH_ENDPOINT"),
		AuthToken:       os.Getenv("OS_TOKEN"),
		ProjectId:       os.Getenv("OS_PROJECT_ID"),
	}

	yamlByte, err := os.ReadFile("../capi-local.yaml") // workload cluster yaml manifest
	if err != nil {
		t.Fatal(error.Error(err))
	}
	yaml := string(yamlByte)

	yamlResult, _ := cl.UpdateClusterYamlManifestFlannel(yaml, option.ManifestSpecOption{
		ClusterKindSpecOption: option.ClusterKindSpecOption{
			CidrBlocks: []string{"10.244.0.0/16"},
		},
		InfrastructureKindSpecOption: option.InfrastructureKindSpecOption{
			AllowAllInClusterTraffic: true,
		},
	})
	if err := os.WriteFile("../capi-local-flannel.yaml", []byte(yamlResult), 0644); err != nil {
		t.Fatal("Write yaml error:", error.Error(err))
	}
}
