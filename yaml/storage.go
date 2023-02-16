package yaml

type (
	StorageClass struct {
		ApiVersion  string                 `yaml:"apiVersion" json:"apiVersion"`
		Kind        string                 `yaml:"kind" json:"kind"`
		Metadata    map[string]interface{} `yaml:"metadata" json:"metadata"`
		Provisioner string                 `yaml:"provisioner" json:"provisioner"`
		Parameters  map[string]interface{} `yaml:"parameters" json:"parameters"`
	}
)
