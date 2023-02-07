package yaml

type (
	StorageClass struct {
		ApiVersion  string                 `yaml:"apiVersion" json:"apiVersion"`
		Kind        string                 `yaml:"kind" json:"kind"`
		Metadata    map[string]interface{} `yaml:"metadata" json:"metadata"`
		Provisioner string                 `yaml:"provisioner" json:"provisioner"`
		Parameters  map[string]interface{} `yaml:"parameters" json:"parameters"`
	}

	PersistentVolumeClaimSpec struct {
		AccessModes      []string  `yaml:"accessModes" json:"accessModes"`
		VolumeMode       string    `yaml:"volumeMode" json:"volumeMode"`
		Resources        Resources `yaml:"resources" json:"resources"`
		StorageClassName string    `yaml:"storageClassName" json:"storageClassName"`
	}

	Resources struct {
		Requests Requests `yaml:"requests" json:"requests"`
	}

	Requests struct {
		Storage string `yaml:"storage" json:"storage"`
	}
)
