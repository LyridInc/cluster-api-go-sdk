package test

import (
	"os"
	"testing"

	"github.com/LyridInc/cluster-api-go-sdk/api"
)

// export $(< test.env)

// go test ./test -v -run ^TestAuthenticateServercore$
func TestAuthenticateServercore(t *testing.T) {
	projectID := os.Getenv("AC_PROJECT_ID")
	token := os.Getenv("AC_TOKEN")
	endpoint := os.Getenv("AC_API_ENDPOINT")
	client := api.NewAmericanCloudClient(projectID, token, endpoint)

	res, err := client.ListClusters()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(res))
}
