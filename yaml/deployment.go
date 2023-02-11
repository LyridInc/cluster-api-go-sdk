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
		ServiceAccount     string                  `yaml:"serviceAccount" json:"serviceAccount"`
		Containers         []Container             `yaml:"containers" json:"containers"`
		Volumes            []interface{}           `yaml:"volumes" json:"volumes"`
		NodeSelector       *map[string]interface{} `yaml:"nodeSelector" json:"nodeSelector"`
		SecurityContext    *map[string]interface{} `yaml:"securityContext" json:"securityContext"`
		Tolerations        *[]interface{}          `yaml:"tolerations" json:"tolerations"`
		ServiceAccountName *string                 `yaml:"serviceAccountName" json:"serviceAccountName"`
		HostNetwork        *bool                   `yaml:"hostNetwork" json:"hostNetwork"`
	}

	Container struct {
		Name            string                  `yaml:"name" json:"name"`
		Image           string                  `yaml:"image" json:"image"`
		Args            []string                `yaml:"args" json:"args"`
		Env             []ContainerEnv          `yaml:"env" json:"env"`
		ImagePullPolicy string                  `yaml:"imagePullPolicy" json:"imagePullPolicy"`
		VolumeMounts    []VolumeMount           `yaml:"volumeMounts" json:"volumeMounts"`
		SecurityContext *map[string]interface{} `yaml:"securityContext" json:"securityContext"`
		Ports           *[]interface{}          `yaml:"ports" json:"ports"`
		LivenessProbe   *map[string]interface{} `yaml:"livenessProbe" json:"livenessProbe"`
	}

	ContainerEnv struct {
		Name      string                  `yaml:"name" json:"name"`
		Value     interface{}             `yaml:"value" json:"value"`
		ValueFrom *map[string]interface{} `yaml:"valueFrom" json:"valueFrom"`
	}

	VolumeMount struct {
		Name             string                  `yaml:"name" json:"name"`
		MountPath        string                  `yaml:"mountPath" json:"mountPath"`
		HostPath         *map[string]interface{} `yaml:"hostPath" json:"hostPath"`
		Secret           *map[string]interface{} `yaml:"secret" json:"secret"`
		ReadOnly         *bool                   `yaml:"readOnly" json:"readOnly"`
		MountPropagation string                  `yaml:"mountPropagation" json:"mountPropagation"`
	}
)
