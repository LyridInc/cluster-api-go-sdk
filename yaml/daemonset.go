package yaml

type (
	DaemonSetSpec struct {
		Selector Selector `yaml:"selector" json:"selector"`
		Template Template `yaml:"template" json:"template"`
	}
)
