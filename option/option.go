package option

type (
	OpenstackGenerateClusterOptions struct {
		ControlPlaneMachineFlavor string
		NodeMachineFlavor         string
		ExternalNetworkId         string
		ImageName                 string
		SshKeyName                string
		DnsNameServers            string
		FailureDomain             string
	}

	GenerateWorkloadClusterOptions struct {
		ClusterName              string
		KubernetesVersion        string
		WorkerMachineCount       int64
		ControlPlaneMachineCount int64
		InfrastructureProvider   string
		Flavor                   string
	}
)

var Namespaces map[string]string = map[string]string{
	"openstack": "capi-system",
	"oci":       "capoci-system",
}
