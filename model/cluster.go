package model

type DockerConfigEntry struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty" datapolicy:"password"`
	Email    string `json:"email,omitempty"`
	Auth     string `json:"auth,omitempty" datapolicy:"token"`
}

type DockerConfig map[string]DockerConfigEntry

type DockerConfigJSON struct {
	Auths DockerConfig `json:"auths" datapolicy:"token"`
	// +optional
	HttpHeaders map[string]string `json:"HttpHeaders,omitempty" datapolicy:"token"`
}

type CreateDockerRegistrySecretArgs struct {
	Username    string            `json:"username"`
	Password    string            `json:"password"`
	Email       string            `json:"email"`
	Server      string            `json:"server"`
	Annotations map[string]string `json:"annotations"`
}

type KubeconfigCluster struct {
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
	Server                   string `yaml:"server"`
}

type KubeconfigUser struct {
	ClientCertificateData string `yaml:"client-certificate-data"`
	ClientKeyData         string `yaml:"client-key-data"`
	Token                 string `yaml:"token"`
}

type KubeconfigContext struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

type KubeconfigConfig struct {
	ApiVersion     string              `yaml:"apiVersion"`
	Kind           string              `yaml:"kind"`
	Clusters       []KubeconfigCluster `yaml:"clusters"`
	Users          []KubeconfigUser    `yaml:"users"`
	Contexts       []KubeconfigContext `yaml:"contexts"`
	CurrentContext string              `yaml:"current-context"`
}

type CRDResource struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Labels map[string]string `yaml:"labels"`
	} `yaml:"metadata"`
	Status struct {
		Phase string `yaml:"phase"`
	} `yaml:"status"`
	Spec struct {
		ProviderID        string `yaml:"providerID"`
		InfrastructureRef struct {
			Name string `yaml:"name"`
		} `yaml:"infrastructureRef"`
	} `yaml:"spec"`
}
