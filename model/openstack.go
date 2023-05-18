package model

import "time"

type (
	QuotasResponse struct {
		Quota Quota `json:"quota"`
	}

	Quota struct {
		Network           QuotaValues `json:"network"`
		Subnet            QuotaValues `json:"subnet"`
		SubnetPool        QuotaValues `json:"subnetpool"`
		Port              QuotaValues `json:"port"`
		Router            QuotaValues `json:"router"`
		FloatingIp        QuotaValues `json:"floatingip"`
		RbacPolicy        QuotaValues `json:"rbac_policy"`
		SecurityGroup     QuotaValues `json:"security_group"`
		SecurityGroupRule QuotaValues `json:"security_group_rule"`
		Trunk             QuotaValues `json:"trunk"`
	}

	QuotaValues struct {
		Limit    int `json:"limit"`
		Used     int `json:"used"`
		Reserved int `json:"reserved"`
	}

	LoadBalancerResponse struct {
		LoadBalancer LoadBalancer `json:"loadbalancer"`
	}

	LoadBalancer struct {
		ID                 string        `json:"id"`
		AdminStateUp       bool          `json:"admin_state_up"`
		Description        string        `json:"description"`
		ProjectID          string        `json:"project_id"`
		ProvisioningStatus string        `json:"provisioning_status"`
		FlavorID           string        `json:"flavor_id"`
		VipSubnetID        string        `json:"vip_subnet_id"`
		VipAddress         string        `json:"vip_address"`
		VipNetworkID       string        `json:"vip_network_id"`
		VipPortID          string        `json:"vip_port_id"`
		AdditionalVips     []interface{} `json:"additional_vips"`
		Provider           string        `json:"provider"`
		CreatedAt          *time.Time    `json:"created_at"`
		UpdatedAt          *time.Time    `json:"updated_at"`
		OperatingStatus    string        `json:"operating_status"`
		Name               string        `json:"name"`
		VipQosPolicyID     string        `json:"vip_qos_policy_id"`
		AvailabilityZone   string        `json:"availability_zone"`
		Tags               []interface{} `json:"tags"`
	}

	MagnumCreateClusterRequest struct {
		Name              string                 `json:"name"`
		MasterCount       int                    `json:"master_count"`
		MasterFlavorID    string                 `json:"master_flavor_id,omitempty"`
		NodeCount         int                    `json:"node_count"`
		FlavorID          string                 `json:"flavor_id"`
		Keypair           string                 `json:"keypair"`
		DockerVolumeSize  int                    `json:"docker_volume_size"`
		ClusterTemplateID string                 `json:"cluster_template_id"`
		CreateTimeout     int                    `json:"create_timeout"`
		Labels            map[string]interface{} `json:"labels,omitempty"`
	}

	MagnumCreateClusterTemplateRequest struct {
		Labels              map[string]string `json:"labels,omitempty"`
		FixedSubnet         string            `json:"fixed_subnet,omitempty"`
		MasterFlavorID      string            `json:"master_flavor_id,omitempty"`
		NoProxy             string            `json:"no_proxy,omitempty"`
		HttpsProxy          string            `json:"https_proxy,omitempty"`
		HttpProxy           string            `json:"http_proxy,omitempty"`
		TLSDisabled         bool              `json:"tls_disabled"`
		KeypairID           string            `json:"keypair_id,omitempty"`
		Public              bool              `json:"public"`
		DockerVolumeSize    int               `json:"docker_volume_size"`
		ServerType          string            `json:"server_type"`
		ExternalNetworkID   string            `json:"external_network_id"`
		ImageID             string            `json:"image_id"`
		VolumeDriver        string            `json:"volume_driver,omitempty"`
		RegistryEnabled     bool              `json:"registry_enabled,omitempty"`
		DockerStorageDriver string            `json:"docker_storage_driver"`
		Name                string            `json:"name"`
		NetworkDriver       string            `json:"network_driver"`
		FixedNetwork        string            `json:"fixed_network,omitempty"`
		COE                 string            `json:"coe"`
		FlavorID            string            `json:"flavor_id,omitempty"`
		MasterLBEnabled     bool              `json:"master_lb_enabled"`
		DNSNameserver       string            `json:"dns_nameserver,omitempty"`
		FloatingIPEnabled   bool              `json:"floating_ip_enabled"`
		Hidden              bool              `json:"hidden,omitempty"`
		Tags                []string          `json:"tags,omitempty"`
		ClusterDistro       string            `json:"cluster_distro,omitempty"`
	}

	MagnumClusterTemplateResponse struct {
		InsecureRegistry    string                 `json:"insecure_registry"`
		Links               []interface{}          `json:"links"`
		HTTPProxy           string                 `json:"http_proxy"`
		UpdatedAt           string                 `json:"updated_at"`
		FloatingIPEnabled   bool                   `json:"floating_ip_enabled"`
		FixedSubnet         string                 `json:"fixed_subnet"`
		MasterFlavorID      string                 `json:"master_flavor_id"`
		UserID              string                 `json:"user_id"`
		UUID                string                 `json:"uuid"`
		NoProxy             string                 `json:"no_proxy"`
		HTTPSProxy          string                 `json:"https_proxy"`
		TLSDisabled         bool                   `json:"tls_disabled"`
		KeypairID           string                 `json:"keypair_id"`
		Hidden              bool                   `json:"hidden"`
		ProjectID           string                 `json:"project_id"`
		Public              bool                   `json:"public"`
		Labels              map[string]interface{} `json:"labels"`
		DockerVolumeSize    int                    `json:"docker_volume_size"`
		ServerType          string                 `json:"server_type"`
		ExternalNetworkID   string                 `json:"external_network_id"`
		ClusterDistro       string                 `json:"cluster_distro"`
		ImageID             string                 `json:"image_id"`
		VolumeDriver        string                 `json:"volume_driver"`
		RegistryEnabled     bool                   `json:"registry_enabled"`
		DockerStorageDriver string                 `json:"docker_storage_driver"`
		APIServerPort       int                    `json:"apiserver_port"`
		Name                string                 `json:"name"`
		CreatedAt           string                 `json:"created_at"`
		NetworkDriver       string                 `json:"network_driver"`
		FixedNetwork        string                 `json:"fixed_network"`
		COE                 string                 `json:"coe"`
		FlavorID            string                 `json:"flavor_id"`
		MasterLBEnabled     bool                   `json:"master_lb_enabled"`
		DNSNameserver       string                 `json:"dns_nameserver"`
	}

	Image struct {
		Status          string      `json:"status"`
		Name            string      `json:"name"`
		Tags            []string    `json:"tags"`
		ContainerFormat string      `json:"container_format"`
		CreatedAt       time.Time   `json:"created_at"`
		DiskFormat      string      `json:"disk_format"`
		UpdatedAt       time.Time   `json:"updated_at"`
		Visibility      string      `json:"visibility"`
		Self            string      `json:"self"`
		MinDisk         int         `json:"min_disk"`
		Protected       bool        `json:"protected"`
		ID              string      `json:"id"`
		File            string      `json:"file"`
		Checksum        string      `json:"checksum"`
		OSHashAlgo      string      `json:"os_hash_algo"`
		OSHashValue     string      `json:"os_hash_value"`
		OSHidden        bool        `json:"os_hidden"`
		Owner           string      `json:"owner"`
		Size            int64       `json:"size"`
		MinRAM          int         `json:"min_ram"`
		Schema          string      `json:"schema"`
		VirtualSize     interface{} `json:"virtual_size"` // can be null or int
	}

	ImageListResponse struct {
		Images []Image `json:"images"`
		First  string  `json:"first"`
		Next   string  `json:"next"`
		Schema string  `json:"schema"`
	}

	Network struct {
		AdminStateUp          bool     `json:"admin_state_up"`
		AvailabilityZoneHints []string `json:"availability_zone_hints"`
		AvailabilityZones     []string `json:"availability_zones"`
		CreatedAt             string   `json:"created_at"`
		DNSDomain             string   `json:"dns_domain"`
		ID                    string   `json:"id"`
		IPv4AddressScope      *string  `json:"ipv4_address_scope"`
		IPv6AddressScope      *string  `json:"ipv6_address_scope"`
		L2Adjacency           bool     `json:"l2_adjacency"`
		MTU                   int      `json:"mtu"`
		Name                  string   `json:"name"`
		PortSecurityEnabled   bool     `json:"port_security_enabled"`
		ProjectID             string   `json:"project_id"`
		QoSPolicyID           string   `json:"qos_policy_id"`
		RevisionNumber        int      `json:"revision_number"`
		RouterExternal        bool     `json:"router:external"`
		Shared                bool     `json:"shared"`
		Status                string   `json:"status"`
		Subnets               []string `json:"subnets"`
		TenantID              string   `json:"tenant_id"`
		UpdatedAt             string   `json:"updated_at"`
		VLanTransparent       bool     `json:"vlan_transparent"`
		Description           string   `json:"description"`
		IsDefault             bool     `json:"is_default"`
	}

	NetworkListResponse struct {
		Networks []Network `json:"networks"`
	}

	NetworkResponse struct {
		Network *Network `json:"network"`
	}

	Subnet struct {
		Name              string        `json:"name"`
		EnableDHCP        bool          `json:"enable_dhcp"`
		NetworkID         string        `json:"network_id"`
		SegmentID         *string       `json:"segment_id"`
		ProjectID         string        `json:"project_id"`
		TenantID          string        `json:"tenant_id"`
		DNSNameServers    []string      `json:"dns_nameservers"`
		DNSPublishFixedIP bool          `json:"dns_publish_fixed_ip"`
		AllocationPools   []Allocation  `json:"allocation_pools"`
		HostRoutes        []interface{} `json:"host_routes"`
		IPVersion         int           `json:"ip_version"`
		GatewayIP         string        `json:"gateway_ip"`
		CIDR              string        `json:"cidr"`
		ID                string        `json:"id"`
		CreatedAt         string        `json:"created_at"`
		Description       string        `json:"description"`
		IPv6AddressMode   *string       `json:"ipv6_address_mode"`
		IPv6RAMode        *string       `json:"ipv6_ra_mode"`
		RevisionNumber    int           `json:"revision_number"`
		ServiceTypes      []interface{} `json:"service_types"`
		SubnetPoolID      *string       `json:"subnetpool_id"`
		Tags              []string      `json:"tags"`
		UpdatedAt         string        `json:"updated_at"`
	}

	Allocation struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}

	SubnetListResponse struct {
		Subnets []Subnet `json:"subnets"`
	}

	SubnetResponse struct {
		Subnet *Subnet `json:"subnet"`
	}

	Flavor struct {
		Disabled  bool   `json:"OS-FLV-DISABLED:disabled"`
		Disk      int    `json:"disk"`
		Ephemeral int    `json:"OS-FLV-EXT-DATA:ephemeral"`
		IsPublic  bool   `json:"os-flavor-access:is_public"`
		ID        string `json:"id"`
		Links     []struct {
			Href string `json:"href"`
			Rel  string `json:"rel"`
		} `json:"links"`
		Name        string                 `json:"name"`
		RAM         int                    `json:"ram"`
		Swap        int                    `json:"swap"`
		VCPUs       int                    `json:"vcpus"`
		RxtxFactor  float64                `json:"rxtx_factor"`
		Description string                 `json:"description"`
		ExtraSpecs  map[string]interface{} `json:"extra_specs"`
	}

	FlavorListResponse struct {
		Flavors []Flavor `json:"flavors"`
	}

	FlavorResponse struct {
		Flavor *Flavor `json:"flavor"`
	}
)
