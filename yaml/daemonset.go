package yaml

type (
	DaemonSetSpec struct {
		Selector       Selector
		UpdateStrategy map[string]interface{}
		Template       Template
	}
)
