package test

import (
	"os"
	"testing"

	"github.com/LyridInc/cluster-api-go-sdk/api"
)

// export $(< test.env)

func TestInfrastructureReadiness(t *testing.T) {
	capi := api.NewClusterApiClient("", "../local.kubeconfig")
	_, err := capi.InfrastructureReadiness("openstack")
	if err != nil {
		t.Fatalf(error.Error(err))
	}
}

func TestUnavailableInfrastructureReadiness(t *testing.T) {
	capi := api.NewClusterApiClient("", "../local.kubeconfig")
	_, err := capi.InfrastructureReadiness("oci")
	if err != nil {
		t.Fatalf(error.Error(err))
	}
}

// go test ./test -v -run ^TestGetWorkloadClusterKubeconfig$
func TestGetWorkloadClusterKubeconfig(t *testing.T) {
	capi := api.NewClusterApiClient("", "../local.kubeconfig")
	t.Run("cluster exists", func(t *testing.T) {
		_, err := capi.GetWorkloadClusterKubeconfig("capi-local-2")
		if err != nil {
			t.Fatalf(error.Error(err))
		}
	})
	t.Run("cluster doesn't exist", func(t *testing.T) {
		_, err := capi.GetWorkloadClusterKubeconfig("capi-local-99")
		if err != nil {
			t.Fatalf(error.Error(err))
		}
	})
}

// go test ./test -v -run ^TestKubectlManifest$
func TestKubectlManifest(t *testing.T) {
	capi := api.NewClusterApiClient("", "../local.kubeconfig")
	yamlByte, err := os.ReadFile("../capi-local.yaml") // workload cluster yaml manifest
	if err != nil {
		t.Fatal(error.Error(err))
	}
	yaml := string(yamlByte)

	t.Run("kubectl apply -f", func(t *testing.T) {
		cl := api.OpenstackClient{
			NetworkEndpoint: os.Getenv("OS_NETWORK_ENDPOINT"),
			AuthEndpoint:    os.Getenv("OS_AUTH_ENDPOINT"),
			AuthToken:       os.Getenv("OS_TOKEN"),
			ProjectId:       os.Getenv("OS_PROJECT_ID"),
		}

		if quotaAvailable, err := cl.ValidateQuotas(); !quotaAvailable || err != nil {
			t.Fatal("Quota problem:", error.Error(err))
		}

		if err := capi.ApplyYaml(yaml); err != nil {
			t.Fatal(error.Error(err))
		}
	})

	// go test ./test -run ^\QTestKubectlManifest\E$/^\Qkubectl_delete_-f\E$
	t.Run("kubectl delete -f", func(t *testing.T) {
		if err := capi.DeleteYaml(yaml); err != nil {
			t.Fatal(error.Error(err))
		}
	})
}
