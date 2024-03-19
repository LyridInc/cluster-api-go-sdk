package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/LyridInc/cluster-api-go-sdk/api"
	"github.com/LyridInc/cluster-api-go-sdk/option"
)

// go test ./test -v -run ^TestGenerateOciClusterTemplate$
func TestGenerateOciClusterTemplate(t *testing.T) {
	infrastructure := "oci"
	capi, _ := api.NewClusterApiClient("", "./data/az-vega.kubeconfig")

	ready, err := capi.InfrastructureReadiness(infrastructure)
	if !ready && err == nil {
		t.Log("initialize infrastructure")
		capi.InitInfrastructure(infrastructure)
	}

	clusterName := "capi-oci-test"
	clusterOpt := option.GenerateOciWorkloadClusterOption{
		CompartmentID:     "ocid1.compartment.oc1..aaaaaaaabgpgr4qf7zzyfoscxn3426wuw4vs5tlw2jsdasjvwt5fdavj3ura",
		ClusterName:       clusterName,
		ImageID:           "",
		Shape:             "VM.Standard.E3.Flex",
		MachineTypeOCPU:   4,
		SSHKey:            "",
		Region:            "ap-singapore-1",
		WorkloadRegion:    "ap-singapore-1",
		Namespace:         "default",
		KubernetesVersion: "v1.28.2",
		MachineCount:      1,
		URL:               "./data/oci/cluster-template-managed-flannel.yaml",
	}
	yaml, err := capi.GenerateOciWorkloadClusterYaml(clusterOpt)
	if err != nil {
		t.Fatal("Generate workload cluster error:", err)
	}

	if err := os.WriteFile(fmt.Sprintf("./data/%s.yaml", clusterName), []byte(yaml), 0644); err != nil {
		t.Fatal("Write yaml error:", err)
	}
}
