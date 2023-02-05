package test

import (
	"os"
	"testing"

	"github.com/LyridInc/cluster-api-go-sdk/api"
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
