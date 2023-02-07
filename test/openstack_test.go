package test

import (
	"os"
	"testing"

	"github.com/LyridInc/cluster-api-go-sdk/api"
	"github.com/LyridInc/cluster-api-go-sdk/model"
	"github.com/LyridInc/cluster-api-go-sdk/option"
)

// export $(< test.env)

// go test ./test -v -run ^TestClientAuthentication$
func TestClientAuthentication(t *testing.T) {
	cl := api.OpenstackClient{
		NetworkEndpoint: os.Getenv("OS_NETWORK_ENDPOINT"),
		AuthEndpoint:    os.Getenv("OS_AUTH_ENDPOINT"),
		AuthToken:       os.Getenv("OS_TOKEN"),
		ProjectId:       os.Getenv("OS_PROJECT_ID"),
	}
	cl.Authenticate(api.OpenstackCredential{
		ApplicationCredentialName:   os.Getenv("OS_APPLICATION_CREDENTIAL_NAME"),
		ApplicationCredentialId:     os.Getenv("OS_APPLICATION_CREDENTIAL_ID"),
		ApplicationCredentialSecret: os.Getenv("OS_APPLICATION_CREDENTIAL_SECRET"),
	})
}

// go test ./test -v -run ^TestCheckAuthToken$
func TestCheckAuthToken(t *testing.T) {
	cl := api.OpenstackClient{
		NetworkEndpoint: os.Getenv("OS_NETWORK_ENDPOINT"),
		AuthEndpoint:    os.Getenv("OS_AUTH_ENDPOINT"),
		AuthToken:       os.Getenv("OS_TOKEN"),
		ProjectId:       os.Getenv("OS_PROJECT_ID"),
	}
	t.Run("check auth token", func(t *testing.T) {
		t.Log("TOKEN:", cl.AuthToken)
		response, err := cl.CheckAuthToken()
		if err != nil {
			t.Fatal(error.Error(err))
		}
		t.Log(response)
	})
}

// go test ./test -v -run ^TestQuota$
func TestQuota(t *testing.T) {
	cl := api.OpenstackClient{
		NetworkEndpoint: os.Getenv("OS_NETWORK_ENDPOINT"),
		AuthEndpoint:    os.Getenv("OS_AUTH_ENDPOINT"),
		AuthToken:       os.Getenv("OS_TOKEN"),
		ProjectId:       os.Getenv("OS_PROJECT_ID"),
	}
	t.Run("show quota", func(t *testing.T) {
		cl.GetProjectQuotas()
	})
}

// go test ./test -v -run ^TestUpdateYamlManifest$
func TestUpdateYamlManifest(t *testing.T) {
	cl := api.OpenstackClient{
		NetworkEndpoint: os.Getenv("OS_NETWORK_ENDPOINT"),
		AuthEndpoint:    os.Getenv("OS_AUTH_ENDPOINT"),
		AuthToken:       os.Getenv("OS_TOKEN"),
		ProjectId:       os.Getenv("OS_PROJECT_ID"),
	}

	t.Run("update manifest for flannel support", func(t *testing.T) {
		yamlByte, err := os.ReadFile("./data/capi-local.yaml") // workload cluster yaml manifest
		if err != nil {
			t.Fatal(error.Error(err))
		}
		yaml := string(yamlByte)

		yamlResult, _ := cl.UpdateYamlManifest(yaml, option.ManifestOption{
			ClusterKindSpecOption: option.ClusterKindSpecOption{
				CidrBlocks: []string{"10.244.0.0/16"},
			},
			InfrastructureKindSpecOption: option.InfrastructureKindSpecOption{
				AllowAllInClusterTraffic: true,
			},
		})
		if err := os.WriteFile("./data/capi-local-flannel.yaml", []byte(yamlResult), 0644); err != nil {
			t.Fatal("Write yaml error:", error.Error(err))
		}
	})

	t.Run("update manifest for csi cinder support", func(t *testing.T) {
		yaml, err := model.ReadYamlFromUrl(option.OPENSTACK_CINDER_MANIFEST_URLS["block"])
		if err != nil {
			t.Fatal("Read yaml from url error:", error.Error(err))
		}

		yamlResult, err := cl.UpdateYamlManifest(yaml, option.ManifestOption{
			StorageClassKindOption: option.StorageClassKindOption{
				Parameters: map[string]interface{}{
					"availability": "nova",
				},
			},
		})
		if err != nil {
			t.Fatal("Update yaml from url error:", error.Error(err))
		}

		if err := os.WriteFile("./data/block-storage.yaml", []byte(yamlResult), 0644); err != nil {
			t.Fatal("Write yaml error:", error.Error(err))
		}
	})

	t.Run("update manifest for csi secret cinderplugin", func(t *testing.T) {
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
			FailureDomain:             "az-01",
		}
		cloudsYaml.SetEnvironment(opt)

		cloudConf := os.Getenv("OPENSTACK_CLOUD_PROVIDER_CONF_B64")
		if cloudConf == "" {
			t.Fatal("Error reading cloud conf: OPENSTACK_CLOUD_PROVIDER_CONF_B64 is not set")
		}

		yaml, err := model.ReadYamlFromUrl(option.OPENSTACK_CINDER_MANIFEST_URLS["secret"])
		if err != nil {
			t.Fatal("Read yaml from url error:", error.Error(err))
		}

		yamlResult, err := cl.UpdateYamlManifest(yaml, option.ManifestOption{
			SecretKindOption: option.SecretKindOption{
				Data: map[string]interface{}{
					"cloud.conf": cloudConf,
				},
			},
		})
		if err != nil {
			t.Fatal("Update yaml from url error:", error.Error(err))
		}

		if err := os.WriteFile("./data/csi-secret-cinderplugin.yaml", []byte(yamlResult), 0644); err != nil {
			t.Fatal("Write yaml error:", error.Error(err))
		}
	})
}
