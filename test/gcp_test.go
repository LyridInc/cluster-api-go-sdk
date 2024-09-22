package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/LyridInc/cluster-api-go-sdk/api"
	"github.com/LyridInc/cluster-api-go-sdk/option"
)

// go test ./test -v -run ^TestGenerateGkeClusterTemplate$
func TestGenerateGkeClusterTemplate(t *testing.T) {
	infrastructure := "gcp"
	capi, _ := api.NewClusterApiClient("", "./data/lyrid-staging.kubeconfig")

	clusterName := "capi-gke"
	ready, err := capi.InfrastructureReadiness(infrastructure)
	if !ready && err == nil {
		t.Log("initialize infrastructure")
		capi.InitInfrastructure(infrastructure)
	}

	t.Log("Generate workload cluster YAML")
	clusterOpt := option.GenerateGkeWorkloadClusterOption{
		ClusterName:             clusterName,
		Namespace:               "default",
		KubernetesVersion:       "v1.31.0",
		ControlPlaneMachineType: "n1-standard-2",
		WorkerMachineType:       "n4-standard-2",
		WorkerMachineCount:      3,
		Project:                 "hubon-gpu2",
		Region:                  "asia-southeast1",
		NetworkName:             "default",
		URL:                     "./data/gcp/cluster-template-gke.yaml",
	}
	yaml, err := capi.GenerateGkeWorkloadClusterYaml(clusterOpt)
	if err != nil {
		t.Fatal("Generate workload cluster error:", err)
	}

	if err := os.WriteFile(fmt.Sprintf("./data/%s.yaml", clusterName), []byte(yaml), 0644); err != nil {
		t.Fatal("Write yaml error:", err)
	}
}
