package model

import "time"

type CivoFirewallResponse struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	AccountID      string  `json:"account_id"`
	RulesCount     int     `json:"rules_count"`
	InstancesCount *int    `json:"instances_count,omitempty"`
	Default        string  `json:"default"`
	Label          *string `json:"label,omitempty"`
	NetworkID      string  `json:"network_id"`
}

type CivoFirewallRuleResponse struct {
	ID         string   `json:"id"`
	FirewallID string   `json:"firewall_id"`
	Protocol   string   `json:"protocol"`
	StartPort  string   `json:"start_port"`
	CIDR       []string `json:"cidr"`
	Direction  string   `json:"direction"`
	Label      string   `json:"label,omitempty"`
}

type CivoKubernetesVersionResponse struct {
	Version     string `json:"version"`
	FlatVersion string `json:"flat_version"`
	Label       string `json:"label"`
	Type        string `json:"type"`
	Default     *bool  `json:"default,omitempty"`
}

type CivoInstanceSizeResponse struct {
	Name        string `json:"name"`
	NiceName    string `json:"nice_name"`
	CPUCores    int    `json:"cpu_cores"`
	RAMMB       int    `json:"ram_mb"`
	DiskGB      int    `json:"disk_gb"`
	Description string `json:"description,omitempty"`
	Selectable  *bool  `json:"selectable,omitempty"`
}

type CivoNetworkResponse struct {
	Result string  `json:"result"`
	Label  *string `json:"label,omitempty"`
	ID     *string `json:"id,omitempty"`
}

type CivoMarketplaceItemResponse struct {
	Name         string     `json:"name"`
	Title        *string    `json:"title,omitempty"`
	Version      string     `json:"version"`
	Default      *string    `json:"default,omitempty"`
	Dependencies []string   `json:"dependencies,omitempty"`
	Maintainer   string     `json:"maintainer"`
	Description  string     `json:"description"`
	PostInstall  *string    `json:"post_install,omitempty"`
	URL          string     `json:"url"`
	Category     string     `json:"category"`
	ImageURL     string     `json:"image_url"`
	Plans        []CivoPlan `json:"plans,omitempty"`
}

type CivoPlan struct {
	Label         string                     `json:"label"`
	Configuration map[string]CivoConfigValue `json:"configuration"`
}

type CivoConfigValue struct {
	Value string `json:"value"`
}

type CivoPool struct {
	ID    string `json:"id"`
	Size  string `json:"size"`
	Count int    `json:"count"`
}

type CivoCreateClusterArgs struct {
	Name              string     `json:"name,omitempty"`
	NetworkID         string     `json:"network_id,omitempty"`
	Region            string     `json:"region,omitempty"`
	CNIPlugin         string     `json:"cni_plugin,omitempty"`
	Pools             []CivoPool `json:"pools,omitempty"`
	KubernetesVersion string     `json:"kubernetes_version,omitempty"`
	Tags              string     `json:"tags,omitempty"`
	InstanceFirewall  string     `json:"instance_firewall,omitempty"`
	FirewallRule      string     `json:"firewall_rule,omitempty"`
	Applications      string     `json:"applications,omitempty"`
}

type CivoCreateNetworkArgs struct {
	Label         string `json:"label,omitempty"`
	Region        string `json:"region,omitempty"`
	CidrV4        string `json:"cidr_v4,omitempty"`
	NameserversV4 string `json:"nameservers_v4,omitempty"`
}

type CivoCreateNetworkResponse struct {
	ID     string `json:"id,omitempty"`
	Label  string `json:"label,omitempty"`
	Result string `json:"result,omitempty"`
}

type CivoCreateFirewallArgs struct {
	Name      string `json:"name,omitempty"`
	Region    string `json:"region,omitempty"`
	NetworkID string `json:"network_id,omitempty"`
}

type CivoCreateFirewallResponse struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Result string `json:"result,omitempty"`
}

type CivoClusterDetailResponse struct {
	ID                    string             `json:"id"`
	Name                  string             `json:"name"`
	Version               string             `json:"version"`
	ClusterType           string             `json:"cluster_type"`
	Status                string             `json:"status"`
	Ready                 bool               `json:"ready"`
	NumTargetNodes        int                `json:"num_target_nodes"`
	TargetNodesSize       string             `json:"target_nodes_size"`
	BuiltAt               time.Time          `json:"built_at"`
	KubernetesVersion     string             `json:"kubernetes_version"`
	APIEndpoint           string             `json:"api_endpoint"`
	DNSEntry              string             `json:"dns_entry"`
	CreatedAt             time.Time          `json:"created_at"`
	UpgradeAvailableTo    string             `json:"upgrade_available_to"`
	MasterIP              string             `json:"master_ip"`
	Pools                 []Pool             `json:"pools"`
	RequiredPools         []CivoRequiredPool `json:"required_pools"`
	FirewallID            string             `json:"firewall_id"`
	MasterIPv6            string             `json:"master_ipv6"`
	Applications          string             `json:"applications"`
	NetworkID             string             `json:"network_id"`
	Namespace             string             `json:"namespace"`
	Size                  string             `json:"size"`
	Count                 int                `json:"count"`
	Kubeconfig            interface{}        `json:"kubeconfig"` // or use *string if nullable string
	Instances             []CivoInstance     `json:"instances"`
	InstalledApplications []CivoApplication  `json:"installed_applications"`
}

type Pool struct {
	ID            string         `json:"id"`
	Size          string         `json:"size"`
	Count         int            `json:"count"`
	InstanceNames []string       `json:"instance_names"`
	Instances     []CivoInstance `json:"instances"`
}

type CivoRequiredPool struct {
	ID     string              `json:"id"`
	Size   string              `json:"size"`
	Count  int                 `json:"count"`
	Labels map[string]string   `json:"labels"` // null in JSON, but could be map
	Taints []map[string]string `json:"taints"` // null in JSON, could also define a Taint struct
}

type CivoInstance struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Hostname        string            `json:"hostname"`
	AccountID       string            `json:"account_id"`
	Size            string            `json:"size"`
	FirewallID      string            `json:"firewall_id"`
	SourceType      string            `json:"source_type"`
	SourceID        string            `json:"source_id"`
	NetworkID       string            `json:"network_id"`
	InitialUser     string            `json:"initial_user"`
	InitialPassword string            `json:"initial_password"`
	SSHKey          string            `json:"ssh_key"`
	Tags            []string          `json:"tags"`
	Script          string            `json:"script"`
	Status          string            `json:"status"`
	CivoStatsdToken string            `json:"civostatsd_token"`
	NamespaceID     string            `json:"namespace_id"`
	Notes           string            `json:"notes"`
	ReverseDNS      string            `json:"reverse_dns"`
	CPUCores        int               `json:"cpu_cores"`
	RAMMB           int               `json:"ram_mb"`
	DiskGB          int               `json:"disk_gb"`
	CreatedAt       time.Time         `json:"created_at"`
	AttachedVolumes interface{}       `json:"attached_volumes"`
	PlacementRule   CivoPlacementRule `json:"placement_rule"`
}

type CivoPlacementRule struct {
	AffinityRules interface{} `json:"affinity_rules"`
	NodeSelector  interface{} `json:"node_selector"`
}

type CivoApplication struct {
	Application      string                 `json:"application"`
	Title            string                 `json:"title"`
	Version          string                 `json:"version"`
	Dependencies     interface{}            `json:"dependencies"`
	Maintainer       string                 `json:"maintainer"`
	Description      string                 `json:"description"`
	PostInstall      string                 `json:"post_install"`
	Installed        bool                   `json:"installed"`
	URL              string                 `json:"url"`
	Category         string                 `json:"category"`
	RequiredServices interface{}            `json:"required_services"`
	UpdatedAt        time.Time              `json:"updated_at"`
	ImageURL         string                 `json:"image_url"`
	Plan             string                 `json:"plan"`
	ClusterType      string                 `json:"cluster_type"`
	Configuration    map[string]interface{} `json:"configuration"`
	KubernetesYAML   string                 `json:"kubernetes_yaml"`
	PreInstallScript string                 `json:"pre_install_script"`
	InstallScript    string                 `json:"install_script"`
	ResourceScript   string                 `json:"resource_script"`
	Disabled         bool                   `json:"disabled"`
	Name             string                 `json:"name"`
}
