package test

import (
	"os"
	"testing"

	svcapi "github.com/LyridInc/cluster-api-go-sdk/api/servercore"
	svcmodel "github.com/LyridInc/cluster-api-go-sdk/model/servercore"
)

// export $(< ./test/data/test.env)

func setupClient(t *testing.T) (*svcapi.ServercoreClient, string) {
	t.Helper() // marks function as test helper

	cl := svcapi.NewServercoreClient(svcmodel.Config{
		ApiKey:                  os.Getenv("SERVERCORE_API_KEY"),
		ApiUrl:                  os.Getenv("SERVERCORE_API_URL"),
		CloudApiUrl:             os.Getenv("SERVERCORE_CLOUD_API_URL"),
		ManagedKubernetesApiUrl: os.Getenv("SERVERCORE_MANAGED_KUBERNETES_API_URL"),
	})

	_, token, err := cl.Authenticate(svcmodel.AuthConfig{
		Username:    os.Getenv("SERVERCORE_USER_NAME"),
		Password:    os.Getenv("SERVERCORE_USER_PASSWORD"),
		Domain:      os.Getenv("SERVERCORE_DOMAIN_NAME"),
		ProjectName: "Lyrid Development",
	})
	if err != nil {
		t.Fatal(err)
	}

	return cl, token
}

// go test ./test -v -run ^TestAuthenticateServercore$
func TestAuthenticateServercore(t *testing.T) {
	_, token := setupClient(t)
	t.Log(token)
}

// go test ./test -v -run ^TestListServercoreClusters$
func TestListServercoreClusters(t *testing.T) {
	cl, token := setupClient(t)

	t.Log(token)

	clusterListResp, err := cl.GetListClusters()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(clusterListResp)
}
