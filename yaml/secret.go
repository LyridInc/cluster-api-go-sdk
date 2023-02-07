package yaml

type (
	Secret struct {
		ApiVersion string                 `yaml:"apiVersion" json:"apiVersion"`
		Kind       string                 `yaml:"kind" json:"kind"`
		Metadata   map[string]interface{} `yaml:"metadata" json:"metadata"`
		Data       map[string]interface{} `yaml:"data" json:"data"`
	}
)
