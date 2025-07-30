package servercore

import "time"

type ClusterListResponse struct {
	Clusters []Cluster `json:"clusters"`
}

type Cluster struct {
	ID                            string            `json:"id"`
	CreatedAt                     time.Time         `json:"created_at"`
	UpdatedAt                     time.Time         `json:"updated_at"`
	Name                          string            `json:"name"`
	Status                        string            `json:"status"`
	ProjectID                     string            `json:"project_id"`
	NetworkID                     string            `json:"network_id"`
	SubnetID                      string            `json:"subnet_id"`
	KubeAPIIP                     string            `json:"kube_api_ip"`
	KubeVersion                   string            `json:"kube_version"`
	Region                        string            `json:"region"`
	AdditionalSoftware            interface{}       `json:"additional_software"` // use a concrete type if structure known
	EnableAutorepair              bool              `json:"enable_autorepair"`
	EnablePatchVersionAutoUpgrade bool              `json:"enable_patch_version_auto_upgrade"`
	KubernetesOptions             KubernetesOptions `json:"kubernetes_options"`
	Zonal                         bool              `json:"zonal"`
	PrivateKubeAPI                bool              `json:"private_kube_api"`
}

type KubernetesOptions struct {
	FeatureGates         []string    `json:"feature_gates"`
	AdmissionControllers []string    `json:"admission_controllers"`
	AuditLogs            AuditLogs   `json:"audit_logs"`
	OIDC                 OIDCOptions `json:"oidc"`
}

type AuditLogs struct {
	Enabled    bool   `json:"enabled"`
	SecretName string `json:"secret_name"`
}

type OIDCOptions struct {
	Enabled       bool   `json:"enabled"`
	IssuerURL     string `json:"issuer_url"`
	ClientID      string `json:"client_id"`
	UsernameClaim string `json:"username_claim"`
	GroupsClaim   string `json:"groups_claim"`
	ProviderName  string `json:"provider_name"`
}

type CreateClusterResponse struct {
	Cluster Cluster `json:"cluster"`
}

type CreateClusterRequest struct {
	Cluster ClusterRequest `json:"cluster"`
}

type ClusterRequest struct {
	EnableAutorepair              bool                     `json:"enable_autorepair"`
	EnablePatchVersionAutoUpgrade bool                     `json:"enable_patch_version_auto_upgrade"`
	KubeVersion                   string                   `json:"kube_version"`
	KubernetesOptions             KubernetesOptionsRequest `json:"kubernetes_options"`
	MaintenanceWindowStart        string                   `json:"maintenance_window_start"`
	Name                          string                   `json:"name"`
	NetworkID                     string                   `json:"network_id"`
	NodeGroups                    []NodeGroupRequest       `json:"nodegroups"`
	PrivateKubeAPI                bool                     `json:"private_kube_api"`
	Region                        string                   `json:"region"`
	SubnetID                      string                   `json:"subnet_id"`
	Zonal                         bool                     `json:"zonal"`
}

type KubernetesOptionsRequest struct {
	AdmissionControllers    []string         `json:"admission_controllers"`
	AuditLogs               AuditLogsRequest `json:"audit_logs"`
	EnablePodSecurityPolicy bool             `json:"enable_pod_security_policy"`
	FeatureGates            []string         `json:"feature_gates"`
	OIDC                    OIDCRequest      `json:"oidc"`
	X509CACertificates      string           `json:"x509_ca_certificates"`
}

type AuditLogsRequest struct {
	Enabled    bool   `json:"enabled"`
	SecretName string `json:"secret_name"`
}

type OIDCRequest struct {
	CACerts       string `json:"ca_certs"`
	ClientID      string `json:"client_id"`
	Enabled       bool   `json:"enabled"`
	GroupsClaim   string `json:"groups_claim"`
	IssuerURL     string `json:"issuer_url"`
	ProviderName  string `json:"provider_name"`
	UsernameClaim string `json:"username_claim"`
}

type NodeGroupRequest struct {
	AffinityPolicy            string            `json:"affinity_policy"`
	AutoscaleMaxNodes         int               `json:"autoscale_max_nodes"`
	AutoscaleMinNodes         int               `json:"autoscale_min_nodes"`
	AvailabilityZone          string            `json:"availability_zone"`
	Count                     int               `json:"count"`
	CPUs                      int               `json:"cpus"`
	EnableAutoscale           bool              `json:"enable_autoscale"`
	FlavorID                  string            `json:"flavor_id"`
	InstallNvidiaDevicePlugin bool              `json:"install_nvidia_device_plugin"`
	KeypairName               string            `json:"keypair_name"`
	Labels                    map[string]string `json:"labels"`
	LocalVolume               bool              `json:"local_volume"`
	Preemptible               bool              `json:"preemptible"`
	RAMMB                     int               `json:"ram_mb"`
	Taints                    []string          `json:"taints"` // Can be a struct if needed
	UserData                  string            `json:"user_data"`
	VolumeGB                  int               `json:"volume_gb"`
	VolumeType                string            `json:"volume_type"`
}
