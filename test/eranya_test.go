package test

import (
	"os"
	"testing"

	"github.com/LyridInc/cluster-api-go-sdk/api"
)

// export $(< ./test/data/test.env)

// go test ./test -v -run ^TestEranyaCloudAuthentication$
func TestEranyaCloudAuthentication(t *testing.T) {
	username := os.Getenv("VT_USERNAME")
	domain := os.Getenv("VT_DOMAIN")
	password := os.Getenv("VT_PASSWORD")
	apiEndpoint := os.Getenv("VT_API_ENDPOINT")

	vClient := api.NewEranyaCloud(apiEndpoint)
	response, err := vClient.Authenticate(domain, username, password)
	if err != nil {
		t.Fatal(error.Error(err))
	}

	t.Log(string(response))
}
