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
)
