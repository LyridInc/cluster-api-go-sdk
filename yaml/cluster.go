package yaml

type (
	ClusterSpec struct {
		ClusterNetwork    ClusterNetwork `yaml:"clusterNetwork" json:"clusterNetwork"`
		ControlPlaneRef   Ref            `yaml:"controlPlaneRef" json:"controlPlaneRef"`
		InfrastructureRef Ref            `yaml:"infrastructureRef" json:"infrastructureRef"`
	}

	ClusterNetwork struct {
		Pods          Pods   `yaml:"pods" json:"pods"`
		ServiceDomain string `yaml:"serviceDomain" json:"serviceDomain"`
	}

	Pods struct {
		CidrBlocks []string `yaml:"cidrBlocks" json:"cidrBlocks"`
	}

	Ref struct {
		ApiVersion string `yaml:"apiVersion" json:"apiVersion"`
		Kind       string `yaml:"kind" json:"kind"`
		Name       string `yaml:"name" json:"name"`
	}
)
