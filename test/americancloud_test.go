package test

import (
	"os"
	"testing"

	"github.com/LyridInc/cluster-api-go-sdk/api"
)

// export $(< test.env)

// go test ./test -v -run ^TestListCluster$
func TestListCluster(t *testing.T) {
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

// go test ./test -v -run ^TestGetClusterKubeconfig$
func TestGetClusterKubeconfig(t *testing.T) {
	projectID := os.Getenv("AC_PROJECT_ID")
	token := os.Getenv("AC_TOKEN")
	endpoint := os.Getenv("AC_API_ENDPOINT")
	client := api.NewAmericanCloudClient(projectID, token, endpoint)

	res, err := client.GetClusterKubeconfig("lyrid-dev")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(res))
}

// go test ./test -v -run ^TestCreateCluster$
func TestCreateCluster(t *testing.T) {
	projectID := os.Getenv("AC_PROJECT_ID")
	token := os.Getenv("AC_TOKEN")
	endpoint := os.Getenv("AC_API_ENDPOINT")
	client := api.NewAmericanCloudClient(projectID, token, endpoint)

	res, err := client.CreateCluster(api.AmericanCloudCreateClusterArgs{
		Name:             "lyrid-dev",
		Project:          "handoyo-sutanto-5551290",
		Zone:             "us-west-0",
		Version:          "1.25.0",
		NodeSize:         1,
		Package:          "Basic ACKS",
		BillingPeriod:    "hourly",
		ControlNodes:     1,
		HighAvailability: false,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(res))
}

// go test ./test -v -run ^TestACDeleteCluster$
func TestACDeleteCluster(t *testing.T) {
	projectID := os.Getenv("AC_PROJECT_ID")
	token := os.Getenv("AC_TOKEN")
	endpoint := os.Getenv("AC_API_ENDPOINT")
	client := api.NewAmericanCloudClient(projectID, token, endpoint)

	res, err := client.DeleteCluster("lyrid-dev")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(res))
}
