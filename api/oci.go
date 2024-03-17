package api

import (
	"fmt"
	"os"
)

var OCI_CLUSTER_TEMPLATE_URL map[string]string = map[string]string{
	"managed_flannel":    "https://github.com/oracle/cluster-api-provider-oci/releases/download/v0.11.5/cluster-template-managed-flannel.yaml",
	"alternative_region": "https://github.com/oracle/cluster-api-provider-oci/releases/download/v0.11.5/cluster-template-alternative-region.yaml",
}

type (
	OciClient struct {
		CompartmentID string
		Region        string
	}

	CreateManagedClusterOption struct {
		ClusterName       string
		ImageID           string
		Shape             string
		MachineTypeOCPU   int
		SSHKey            string
		WorkloadRegion    string
		Namespace         string
		KubernetesVersion string
		MachineCount      int
	}
)

func (c *OciClient) CreateManagedCluster(opt CreateManagedClusterOption) {
	os.Setenv("OCI_MANAGED_NODE_IMAGE_ID", opt.ImageID)
	os.Setenv("OCI_MANAGED_NODE_SHAPE", opt.Shape)
	os.Setenv("OCI_MANAGED_NODE_MACHINE_TYPE_OCPUS", fmt.Sprintf("%d", opt.MachineTypeOCPU))
	os.Setenv("OCI_SSH_KEY", opt.SSHKey)
	os.Setenv("OCI_REGION", c.Region)
	os.Setenv("OCI_WORKLOAD_REGION", opt.WorkloadRegion)
	os.Setenv("KUBERNETES_VERSION", opt.KubernetesVersion)
	os.Setenv("NAMESPACE", opt.Namespace)
	os.Setenv("NODE_MACHINE_COUNT", fmt.Sprintf("%d", opt.MachineCount))
}
