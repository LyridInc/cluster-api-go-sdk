package model

type (
	VirtuozzoAuth struct {
		Identity VirtuozzoIdentity `json:"identity"`
		Scope    VirtuozzoScope    `json:"scope"`
	}

	VirtuozzoIdentity struct {
		Methods  []string         `json:"methods"`
		Password IdentityPassword `json:"password"`
	}

	VirtuozzoScope struct {
		Project ScopeProject `json:"project"`
	}

	ScopeProject struct {
		Name   string                 `json:"name"`
		Domain map[string]interface{} `json:"domain"`
	}

	IdentityPassword struct {
		User UserIdentity `json:"user"`
	}

	UserIdentity struct {
		Name     string                 `json:"name"`
		Domain   map[string]interface{} `json:"domain"`
		Password string                 `json:"password"`
	}
)

type (
	// https://docs.virtuozzo.com/virtuozzo_hybrid_infrastructure_5_4_compute_api_reference/index.html#creating-kubernetes-cluster-templates.html
	VirtuozzoCreateKubernetesClusterTemplateArgs struct {
		FloatingIPEnabled   bool              `json:"floating_ip_enabled"`
		FixedSubnet         string            `json:"fixed_subnet"`
		MasterFlavorID      string            `json:"master_flavor_id"`
		NoProxy             string            `json:"no_proxy"`
		HttpsProxy          string            `json:"https_proxy"`
		TLSDisabled         bool              `json:"tls_disabled"`
		KeypairID           string            `json:"keypair_id"`
		Public              bool              `json:"public"`
		Labels              map[string]string `json:"labels"`
		DockerVolumeSize    int               `json:"docker_volume_size"`
		ServerType          string            `json:"server_type"`
		ExternalNetworkID   string            `json:"external_network_id"`
		ImageID             string            `json:"image_id"`
		VolumeDriver        string            `json:"volume_driver,omitempty"`
		RegistryEnabled     bool              `json:"registry_enabled"`
		DockerStorageDriver string            `json:"docker_storage_driver"`
		Name                string            `json:"name"`
		NetworkDriver       string            `json:"network_driver,omitempty"`
		FixedNetwork        string            `json:"fixed_network"`
		COE                 string            `json:"coe,omitempty"`
		FlavorID            string            `json:"flavor_id"`
		MasterLBEnabled     bool              `json:"master_lb_enabled"`
		DNSNameserver       string            `json:"dns_nameserver"`
		Hidden              bool              `json:"hidden"`
	}

	VirtuozzoKubernetesClusterResponse struct {
		InsecureRegistry    string            `json:"insecure_registry"`
		Links               []interface{}     `json:"links"`
		UpdatedAt           string            `json:"updated_at"`
		FloatingIPEnabled   bool              `json:"floating_ip_enabled"`
		FixedSubnet         string            `json:"fixed_subnet"`
		MasterFlavorID      string            `json:"master_flavor_id"`
		UUID                string            `json:"uuid"`
		NoProxy             string            `json:"no_proxy"`
		HttpsProxy          string            `json:"https_proxy"`
		TLSDisabled         bool              `json:"tls_disabled"`
		KeypairID           string            `json:"keypair_id"`
		Public              bool              `json:"public"`
		Labels              map[string]string `json:"labels"`
		DockerVolumeSize    int               `json:"docker_volume_size"`
		ExternalNetworkID   string            `json:"external_network_id"`
		ClusterDistro       string            `json:"cluster_distro"`
		ImageID             string            `json:"image_id"`
		RegistryEnabled     bool              `json:"registry_enabled"`
		DockerStorageDriver string            `json:"docker_storage_driver"`
		APIServerPort       int               `json:"apiserver_port"`
		Name                string            `json:"name"`
		CreatedAt           string            `json:"created_at"`
		NetworkDriver       string            `json:"network_driver"`
		FixedNetwork        string            `json:"fixed_network"`
		COE                 string            `json:"coe"`
		FlavorID            string            `json:"flavor_id"`
		MasterLBEnabled     bool              `json:"master_lb_enabled"`
		DNSNameserver       string            `json:"dns_nameserver"`
		Hidden              bool              `json:"hidden"`
	}
)
