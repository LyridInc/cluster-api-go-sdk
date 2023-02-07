package yaml

type (
	InfrastructureSpec struct {
		CloudName                string                `yaml:"cloudName" json:"cloudName"`
		ExternalNetworkId        string                `yaml:"externalNetworkId" json:"externalNetworkId"`
		NodeCidr                 string                `yaml:"nodeCidr" json:"nodeCidr"`
		ManagedSecurityGroups    bool                  `yaml:"managedSecurityGroups" json:"managedSecurityGroups"`
		AllowAllInClusterTraffic bool                  `yaml:"allowAllInClusterTraffic" json:"allowAllInClusterTraffic"`
		DnsNameServers           []string              `yaml:"dnsNameservers" json:"dnsNameservers"`
		IdentityRef              IdentityRef           `yaml:"identityRef" json:"identityRef"`
		ApiServerLoadBalancer    ApiServerLoadBalancer `yaml:"apiServerLoadBalancer" json:"apiServerLoadBalancer"`
	}

	ApiServerLoadBalancer struct {
		Enabled bool `yaml:"enabled" json:"enabled"`
	}

	IdentityRef struct {
		Kind string `yaml:"kind" json:"kind"`
		Name string `yaml:"name" json:"name"`
	}
)
