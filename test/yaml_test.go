package test

import (
	"os"
	"strings"
	"testing"

	"github.com/LyridInc/cluster-api-go-sdk/api"
	"github.com/LyridInc/cluster-api-go-sdk/model"
	"github.com/LyridInc/cluster-api-go-sdk/option"
)

// go test ./test -v -run ^TestReadYamlFromUrl$
func TestReadYaml(t *testing.T) {
	t.Run("read flannel manifest from url", func(t *testing.T) {
		yaml, err := model.ReadYamlFromUrl(option.FLANNEL_MANIFEST_URL)
		if err != nil {
			t.Fatal(error.Error(err))
		}
		t.Log(yaml)
	})
}

// go test ./test -v -run ^TestCloudsYaml$
func TestCloudsYaml(t *testing.T) {
	yamlByte, _ := os.ReadFile("./data/clouds.yaml")
	cloudsYaml := model.CloudsYaml{}
	cloudsYaml.Parse(yamlByte)
	opt := option.OpenstackGenerateClusterOptions{
		ControlPlaneMachineFlavor: "SS2.2",
		NodeMachineFlavor:         "SM8.4",
		ExternalNetworkId:         "79241ddc-c51b-4677-a763-f48c60870923",
		ImageName:                 "ubuntu-2004-kube-v1.24.8",
		SshKeyName:                "kube-key",
		DnsNameServers:            "8.8.8.8",
		FailureDomain:             "az-01", // nova/az-01
		IgnoreVolumeAZ:            true,
	}
	cloudsYaml.SetEnvironment(opt)

	envs := os.Environ()
	filteredEnvs := []string{}
	for _, env := range envs {
		if strings.HasPrefix(env, "OPENSTACK") || strings.HasPrefix(env, "CAPO") {
			filteredEnvs = append(filteredEnvs, env)
		}
	}

	t.Log(filteredEnvs)
}

// go test ./test -v -run ^TestUpdateYaml$
func TestUpdateYaml(t *testing.T) {
	cl := api.OpenstackClient{
		NetworkEndpoint: os.Getenv("OS_NETWORK_ENDPOINT"),
		AuthEndpoint:    os.Getenv("OS_AUTH_ENDPOINT"),
		AuthToken:       os.Getenv("OS_TOKEN"),
		ProjectId:       os.Getenv("OS_PROJECT_ID"),
	}

	url := "https://raw.githubusercontent.com/kubernetes/cloud-provider-openstack/master/manifests/cinder-csi-plugin/cinder-csi-nodeplugin.yaml"

	yaml, err := model.ReadYamlFromUrl(url)
	if err != nil {
		t.Fatal(error.Error(err))
	}

	yamlResult, err := cl.UpdateYamlManifest(yaml, option.ManifestOption{
		DaemonSetKindOption: option.DaemonSetKindOption{
			VolumeSecretName: "test-secret-update",
		},
	})
	if err != nil {
		t.Fatal("Update yaml from url error:", error.Error(err))
	}

	if err := os.WriteFile("./data/test-daemonset.yaml", []byte(yamlResult), 0644); err != nil {
		t.Fatal("Write yaml file error:", url, error.Error(err))
	}
}
