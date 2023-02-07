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

	ClusterKindSpecOption struct {
		CidrBlocks []string
	}

	InfrastructureKindSpecOption struct {
		AllowAllInClusterTraffic bool
	}

	ManifestSpecOption struct {
		ClusterKindSpecOption        ClusterKindSpecOption
		InfrastructureKindSpecOption InfrastructureKindSpecOption
	}
)

var Namespaces map[string]string = map[string]string{
	"openstack": "capo-system",
	"oci":       "capoci-system",
}

const FLANNEL_MANIFEST_URL = "https://raw.githubusercontent.com/flannel-io/flannel/master/Documentation/kube-flannel.yml"

var OPENSTACK_CLOUD_CONTROLLER_MANIFEST_URLS = []string{
	"https://raw.githubusercontent.com/kubernetes/cloud-provider-openstack/master/manifests/controller-manager/cloud-controller-manager-roles.yaml",
	"https://raw.githubusercontent.com/kubernetes/cloud-provider-openstack/master/manifests/controller-manager/cloud-controller-manager-role-bindings.yaml",
	"https://raw.githubusercontent.com/kubernetes/cloud-provider-openstack/master/manifests/controller-manager/openstack-cloud-controller-manager-ds.yaml",
}

var OPENSTACK_CINDER_MANIFEST_URLS = []string{}
