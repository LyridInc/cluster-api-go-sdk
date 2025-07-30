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
