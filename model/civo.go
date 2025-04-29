package model

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
