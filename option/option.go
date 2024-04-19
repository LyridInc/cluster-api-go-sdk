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
		IgnoreVolumeAZ            bool
	}

	GenerateWorkloadClusterOptions struct {
		ClusterName              string
		TargetNamespace          string
		KubernetesVersion        string
		WorkerMachineCount       int64
		ControlPlaneMachineCount int64
		InfrastructureProvider   string
		Flavor                   string
		URL                      string
	}

	GenerateOciWorkloadClusterOption struct {
		CompartmentID     string
		ClusterName       string
		ImageID           string
		Shape             string
		MachineTypeOCPU   int
		SSHKey            string
		WorkloadRegion    string
		Namespace         string
		KubernetesVersion string
		MachineCount      int64
		Region            string
		URL               string
	}

	GenerateCloudStackWorkloadClusterOption struct {
		ZoneName                    string
		NetworkName                 string
		ClusterEndpointIP           string
		ClusterEndpointPort         string
		ControlPlaneMachineOffering string
		WorkerMachineOffering       string
		TemplateName                string
		SshKeyName                  string
		ClusterName                 string
		Namespace                   string
		KubernetesVersion           string
		WorkerMachineCount          int64
		ControlPlaneMachineCount    int64
		URL                         string
	}

	ClusterKindSpecOption struct {
		CidrBlocks []string
	}

	InfrastructureKindSpecOption struct {
		AllowAllInClusterTraffic bool
		NodeCidr                 string
	}

	StorageClassKindOption struct {
		Parameters map[string]interface{}
	}

	SecretKindOption struct {
		Data     map[string]interface{}
		Metadata map[string]interface{}
	}

	DaemonSetKindOption struct {
		VolumeSecretName string
	}

	DeploymentKindOption struct {
		VolumeSecretName string
	}

	PersistentVolumeClaimKindOption struct {
		Metadata         map[string]interface{}
		Storage          string
		VolumeMode       string
		StorageClassName string
	}

	ManifestOption struct {
		ClusterKindSpecOption           ClusterKindSpecOption
		InfrastructureKindSpecOption    InfrastructureKindSpecOption
		StorageClassKindOption          StorageClassKindOption
		SecretKindOption                SecretKindOption
		DeploymentKindOption            DeploymentKindOption
		DaemonSetKindOption             DaemonSetKindOption
		PersistentVolumeClaimKindOption PersistentVolumeClaimKindOption
	}
)

var Namespaces map[string]string = map[string]string{
	"openstack": "capo-system",
	"oci":       "capoci-system",
}

const FLANNEL_MANIFEST_URL = "https://raw.githubusercontent.com/flannel-io/flannel/master/Documentation/kube-flannel.yml"
const WEAVE_MANIFEST_URL = "https://github.com/weaveworks/weave/releases/download/v2.8.1/weave-daemonset-k8s.yaml"

var OPENSTACK_CLOUD_CONTROLLER_MANIFEST_URLS = []string{
	"https://raw.githubusercontent.com/kubernetes/cloud-provider-openstack/master/manifests/controller-manager/cloud-controller-manager-roles.yaml",
	"https://raw.githubusercontent.com/kubernetes/cloud-provider-openstack/master/manifests/controller-manager/cloud-controller-manager-role-bindings.yaml",
	"https://raw.githubusercontent.com/kubernetes/cloud-provider-openstack/master/manifests/controller-manager/openstack-cloud-controller-manager-ds.yaml",
}

var OPENSTACK_CINDER_DRIVER_MANIFEST_URLS = map[string]interface{}{
	"secret": "https://raw.githubusercontent.com/kubernetes/cloud-provider-openstack/master/manifests/cinder-csi-plugin/csi-secret-cinderplugin.yaml",
	"plugins": []string{
		"https://raw.githubusercontent.com/kubernetes/cloud-provider-openstack/master/manifests/cinder-csi-plugin/csi-cinder-driver.yaml",
		"https://raw.githubusercontent.com/kubernetes/cloud-provider-openstack/master/manifests/cinder-csi-plugin/cinder-csi-controllerplugin-rbac.yaml",
		"https://raw.githubusercontent.com/kubernetes/cloud-provider-openstack/master/manifests/cinder-csi-plugin/cinder-csi-controllerplugin.yaml",
		"https://raw.githubusercontent.com/kubernetes/cloud-provider-openstack/master/manifests/cinder-csi-plugin/cinder-csi-nodeplugin-rbac.yaml",
		"https://raw.githubusercontent.com/kubernetes/cloud-provider-openstack/master/manifests/cinder-csi-plugin/cinder-csi-nodeplugin.yaml",
	},
}

var OPENSTACK_CINDER_STORAGE_URLS = map[string]string{
	"block": "https://raw.githubusercontent.com/kubernetes/cloud-provider-openstack/master/examples/cinder-csi-plugin/block/block.yaml",
}

var OCI_CLUSTER_TEMPLATE_URL map[string]string = map[string]string{
	"managed_flannel":    "https://github.com/oracle/cluster-api-provider-oci/releases/download/v0.11.5/cluster-template-managed-flannel.yaml",
	"alternative_region": "https://github.com/oracle/cluster-api-provider-oci/releases/download/v0.11.5/cluster-template-alternative-region.yaml",
}
