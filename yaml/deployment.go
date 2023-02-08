package yaml

type (
	DeploymentSpec struct {
		Replicas float64  `yaml:"replicas" json:"replicas"`
		Strategy Strategy `yaml:"strategy" json:"strategy"`
		Selector Selector `yaml:"selector" json:"selector"`
		Template Template `yaml:"template" json:"template"`
	}

	Strategy struct {
		Type          string        `yaml:"type" json:"type"`
		RollingUpdate RollingUpdate `yaml:"rollingUpdate" json:"rollingUpdate"`
	}

	RollingUpdate struct {
		MaxUnavailable float64 `yaml:"maxUnavailable" json:"maxUnavailable"`
		MaxSurge       float64 `yaml:"maxSurge" json:"maxSurge"`
	}

	Selector struct {
		MatchLabels map[string]interface{} `yaml:"matchLabels" json:"matchLabels"`
	}

	Template struct {
		Metadata TemplateMetadata `yaml:"metadata" json:"metadata"`
		Spec     TemplateSpec     `yaml:"spec" json:"spec"`
	}

	TemplateMetadata struct {
		Labels map[string]interface{} `yaml:"labels" json:"labels"`
	}

	TemplateSpec struct {
		ServiceAccount string                   `yaml:"serviceAccount" json:"serviceAccount"`
		Containers     []Container              `yaml:"containers" json:"containers"`
		Volumes        []map[string]interface{} `yaml:"volumes" json:"volumes"`
	}

	Container struct {
		Name            string         `yaml:"name" json:"name"`
		Image           string         `yaml:"image" json:"image"`
		Args            []string       `yaml:"args" json:"args"`
		Env             []ContainerEnv `yaml:"env" json:"env"`
		ImagePullPolicy string         `yaml:"imagePullPolicy" json:"imagePullPolicy"`
		VolumeMounts    []VolumeMount  `yaml:"volumeMounts" json:"volumeMounts"`
	}

	ContainerEnv struct {
		Name  string      `yaml:"name" json:"name"`
		Value interface{} `yaml:"value" json:"value"`
	}

	VolumeMount struct {
		Name      string `yaml:"name" json:"name"`
		MountPath string `yaml:"mountPath" json:"mountPath"`
	}
)
