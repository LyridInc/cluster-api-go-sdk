package yaml

import "github.com/LyridInc/cluster-api-go-sdk/option"

type (
	InfrastructureSpec struct {
		ExternalNetwork       ExternalNetworkId `yaml:"externalNetwork" json:"externalNetwork"`
		NodeCidr              string            `yaml:"nodeCidr" json:"nodeCidr"`
		ManagedSecurityGroups any               `yaml:"managedSecurityGroups" json:"managedSecurityGroups"`
		// AllowAllInClusterTraffic bool                  `yaml:"allowAllInClusterTraffic" json:"allowAllInClusterTraffic"`

		IdentityRef           IdentityRef            `yaml:"identityRef" json:"identityRef"`
		ApiServerLoadBalancer ApiServerLoadBalancer  `yaml:"apiServerLoadBalancer" json:"apiServerLoadBalancer"`
		ManagedSubnets        []option.ManagedSubnet `yaml:"managedSubnets" json:"managedSubnets"`
	}

	ApiServerLoadBalancer struct {
		Enabled bool `yaml:"enabled" json:"enabled"`
	}

	IdentityRef struct {
		// Kind string `yaml:"kind" json:"kind"`
		CloudName string `yaml:"cloudName" json:"cloudName"`
		Name      string `yaml:"name" json:"name"`
	}

	ExternalNetworkId struct {
		Id string `yaml:"id" json:"id"`
	}
)
