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
