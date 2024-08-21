package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/LyridInc/cluster-api-go-sdk/api"
	"github.com/LyridInc/cluster-api-go-sdk/option"
)

// go test ./test -v -run ^TestGenerateAwsClusterTemplate$
func TestGenerateAwsClusterTemplate(t *testing.T) {
	infrastructure := "aws"
	capi, _ := api.NewClusterApiClient("", "./data/lyrid-staging.kubeconfig")

	clusterName := "eks-vpccni"
	ready, err := capi.InfrastructureReadiness(infrastructure)
	if !ready && err == nil {
		t.Log("initialize infrastructure")
		capi.InitInfrastructure(infrastructure)
	}

	t.Log("Generate workload cluster YAML")
	clusterOpt := option.GenerateAwsWorkloadClusterOption{
		ClusterName:        clusterName,
		Namespace:          "default",
		Region:             "us-west-1",
		KubernetesVersion:  "v1.30.11",
		SshKeyName:         "azhary-capa",
		Flavor:             "t3.xlarge",
		WorkerMachineCount: 2,
		URL:                "./data/aws/cluster-template-eks-managedmachinepool-vpccni.yaml",
	}
	yaml, err := capi.GenerateAwsWorkloadClusterYaml(clusterOpt)
	if err != nil {
		t.Fatal("Generate workload cluster error:", err)
	}

	if err := os.WriteFile(fmt.Sprintf("./data/%s.yaml", clusterName), []byte(yaml), 0644); err != nil {
		t.Fatal("Write yaml error:", err)
	}
}
