package test

import (
	"io"
	"os"
	"testing"

	"github.com/LyridInc/cluster-api-go-sdk/api"
	"github.com/LyridInc/cluster-api-go-sdk/model"
	"github.com/LyridInc/cluster-api-go-sdk/option"
)

// export $(< test.env)

// go test ./test -v -run ^TestClientCredentialAuthentication$
func TestClientCredentialAuthentication(t *testing.T) {
	cl := api.OpenstackClient{
		NetworkEndpoint: os.Getenv("OS_NETWORK_ENDPOINT"),
		AuthEndpoint:    os.Getenv("OS_AUTH_ENDPOINT"),
		AuthToken:       os.Getenv("OS_TOKEN"),
		ProjectId:       os.Getenv("OS_PROJECT_ID"),
	}

	credential := api.OpenstackAuth{
		Identity: api.OpenstackIdentity{
			Methods: []string{"application_credential"},
			ApplicationCredential: api.OpenstackCredential{
				ApplicationCredentialName:   os.Getenv("OS_APPLICATION_CREDENTIAL_NAME"),
				ApplicationCredentialId:     os.Getenv("OS_APPLICATION_CREDENTIAL_ID"),
				ApplicationCredentialSecret: os.Getenv("OS_APPLICATION_CREDENTIAL_SECRET"),
			},
		},
	}

	if _, err := cl.Authenticate(credential); err != nil {
		t.Fatal(err)
	}
	if _, err := cl.Authenticate(credential); err != nil {
		t.Fatal("#2: ", err)
	}

	os.Setenv("OS_TOKEN", "")
}

// go test ./test -v -run ^TestClientPasswordAuthentication$
func TestClientPasswordAuthentication(t *testing.T) {
	cl := api.OpenstackClient{
		NetworkEndpoint: os.Getenv("OS_NETWORK_ENDPOINT"),
		AuthEndpoint:    os.Getenv("OS_AUTH_ENDPOINT"),
		AuthToken:       os.Getenv("OS_TOKEN"),
		ProjectId:       os.Getenv("OS_PROJECT_ID"),
	}

	credential := api.OpenstackAuth{
		Identity: api.OpenstackIdentity{
			Methods: []string{"password"},
			Password: api.OpenstackPassword{
				User: map[string]interface{}{
					"name":     os.Getenv("OS_USERNAME"),
					"password": os.Getenv("OS_PASSWORD"),
					"domain": map[string]string{
						"id": os.Getenv("OS_USER_DOMAIN_NAME"),
					},
				},
			},
		},
	}

	res, err := cl.Authenticate(credential)
	if err != nil {
		t.Fatal(err)
	}

	b, _ := io.ReadAll(res.Body)
	t.Log(string(b))

	if _, err := cl.Authenticate(credential); err != nil {
		t.Fatal("#2: ", err)
	}
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
		res, err := cl.GetProjectQuotas()
		if err != nil {
			t.Fatal(error.Error(err))
		}
		t.Log(res)
	})
}

// go test ./test -v -run ^TestUpdateDeploymentKindManifest$
func TestUpdateDeploymentKindManifest(t *testing.T) {
	cl := api.OpenstackClient{
		NetworkEndpoint: os.Getenv("OS_NETWORK_ENDPOINT"),
		AuthEndpoint:    os.Getenv("OS_AUTH_ENDPOINT"),
		AuthToken:       os.Getenv("OS_TOKEN"),
		ProjectId:       os.Getenv("OS_PROJECT_ID"),
	}

	url := "https://raw.githubusercontent.com/kubernetes/cloud-provider-openstack/master/manifests/cinder-csi-plugin/cinder-csi-controllerplugin.yaml"

	yaml, err := model.ReadYamlFromUrl(url)
	if err != nil {
		t.Fatal(error.Error(err))
	}

	yamlResult, err := cl.UpdateYamlManifest(yaml, option.ManifestOption{
		DeploymentKindOption: option.DeploymentKindOption{
			VolumeSecretName: "capi-local-2-csi-secret",
		},
	})
	if err != nil {
		t.Fatal("Update yaml from url error:", error.Error(err))
	}

	if err := os.WriteFile("./data/test-controller-plugin-deployment.yaml", []byte(yamlResult), 0644); err != nil {
		t.Fatal("Write yaml file error:", url, error.Error(err))
	}
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
		yaml, err := model.ReadYamlFromUrl(option.OPENSTACK_CINDER_STORAGE_URLS["block"])
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
			FailureDomain:             "az-01", // nova/az-01
			IgnoreVolumeAZ:            true,
		}
		cloudsYaml.SetEnvironment(opt)

		cloudConf := os.Getenv("OPENSTACK_CLOUD_PROVIDER_CONF_B64")
		if cloudConf == "" {
			t.Fatal("Error reading cloud conf: OPENSTACK_CLOUD_PROVIDER_CONF_B64 is not set")
		}

		secretCinderManifestUrl := option.OPENSTACK_CINDER_DRIVER_MANIFEST_URLS["secret"].(string)
		yaml, err := model.ReadYamlFromUrl(secretCinderManifestUrl)
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

// go test ./test -v -run ^TestDeployCniFlannel$
func TestDeployCniFlannel(t *testing.T) {
	capi, _ := api.NewClusterApiClient("", "./data/local.kubeconfig")
	if err := capi.SetKubernetesClientset("./data/capi-local-2.kubeconfig"); err != nil {
		t.Fatal("Error set kubeconfig:", error.Error(err))
	}

	// kubectl --kubeconfig=C:/Users/Lyrid/Documents/Projects/cluster-api-sdk/test/data/capi-local-3.kubeconfig delete -f ./test/data/cloud-controller-manager-roles.yaml

	t.Run("apply flannel cni", func(t *testing.T) {
		yaml, err := model.ReadYamlFromUrl(option.FLANNEL_MANIFEST_URL)
		if err != nil {
			t.Fatal(error.Error(err))
		}

		if err := capi.ApplyYaml(yaml); err != nil {
			t.Fatal("Error apply flannel cni yaml:", error.Error(err))
		}
	})

	t.Run("apply controller manager", func(t *testing.T) {
		for _, url := range option.OPENSTACK_CLOUD_CONTROLLER_MANIFEST_URLS {
			yaml, err := model.ReadYamlFromUrl(url)
			if err != nil {
				t.Fatal(error.Error(err))
			}

			if err := capi.ApplyYaml(yaml); err != nil {
				t.Fatal("Error apply yaml:", url, " - ", error.Error(err))
			}
		}
	})
}

// go test ./test -v -run ^TestInstallCinderCsiDriver$
func TestInstallCinderCsiDriver(t *testing.T) {
	capi, _ := api.NewClusterApiClient("", "./data/local.kubeconfig")
	if err := capi.SetKubernetesClientset("./data/capi-local-2.kubeconfig"); err != nil {
		t.Fatal("Error set kubeconfig:", error.Error(err))
	}

	t.Run("create cinder csi secret", func(t *testing.T) {
		yamlByte, _ := os.ReadFile("./data/csi-secret-cinderplugin.yaml")
		yaml := string(yamlByte)

		if err := capi.ApplyYaml(yaml); err != nil {
			t.Fatal("Error create cinder csi secret:", error.Error(err))
		}
	})

	t.Run("install cinder csi driver", func(t *testing.T) {
		pluginUrls := option.OPENSTACK_CINDER_DRIVER_MANIFEST_URLS["plugins"].([]string)
		for _, url := range pluginUrls {
			yaml, err := model.ReadYamlFromUrl(url)
			if err != nil {
				t.Fatal(error.Error(err))
			}

			if err := capi.ApplyYaml(yaml); err != nil {
				t.Fatal("Error apply yaml:", url, " - ", error.Error(err))
			}
		}
	})

	t.Run("provision block storage", func(t *testing.T) {
		yamlByte, _ := os.ReadFile("./data/block-storage.yaml")
		yaml := string(yamlByte)

		if err := capi.ApplyYaml(yaml); err != nil {
			t.Fatal("Error provisioning block storage:", error.Error(err))
		}
	})
}

// go test ./test -v -run ^TestCreateStorageClass$
func TestCreateStorageClass(t *testing.T) {
	capi, _ := api.NewClusterApiClient("", "./data/capi-testing.kubeconfig")

	yamlByte, _ := os.ReadFile("./data/sc.yaml")
	yaml := string(yamlByte)

	if err := capi.ApplyYaml(yaml); err != nil {
		// errMessage := error.Error(err)
		// strings.Contains()
		t.Fatal("Error create storage class:", error.Error(err))
	}
}

// go test ./test -v -run ^TestCreatePersistentVolumeClaim$
func TestCreatePersistentVolumeClaim(t *testing.T) {
	cl := api.OpenstackClient{
		NetworkEndpoint: os.Getenv("OS_NETWORK_ENDPOINT"),
		AuthEndpoint:    os.Getenv("OS_AUTH_ENDPOINT"),
		AuthToken:       os.Getenv("OS_TOKEN"),
		ProjectId:       os.Getenv("OS_PROJECT_ID"),
	}
	capi, _ := api.NewClusterApiClient("", "./data/capi-testing.kubeconfig")

	yamlByte, _ := os.ReadFile("./data/custom-pvc.yaml")
	yaml := string(yamlByte)

	yamlResult, err := cl.UpdateYamlManifest(yaml, option.ManifestOption{
		PersistentVolumeClaimKindOption: option.PersistentVolumeClaimKindOption{
			Metadata: map[string]interface{}{
				"name": "csi-pvc-cinderplugin-custom",
			},
			Storage: "10Gi",
		},
	})
	if err != nil {
		t.Fatal("Update yaml from url error:", error.Error(err))
	}

	if err := os.WriteFile("./data/custom-pvc.yaml", []byte(yamlResult), 0644); err != nil {
		t.Fatal("Write yaml error:", error.Error(err))
	}

	if err := capi.ApplyYaml(yaml); err != nil {
		t.Fatal("Error create storage class:", error.Error(err))
	}
}
