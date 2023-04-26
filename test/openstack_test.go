package test

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/LyridInc/cluster-api-go-sdk/api"
	"github.com/LyridInc/cluster-api-go-sdk/model"
	"github.com/LyridInc/cluster-api-go-sdk/option"
	"github.com/LyridInc/cluster-api-go-sdk/utils"
)

// export $(< ./test/data/test.env)

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

// go test ./test -v -run ^TestGetLoadBalancer$
func TestGetLoadBalancer(t *testing.T) {
	cl := api.OpenstackClient{
		NetworkEndpoint:      os.Getenv("OS_NETWORK_ENDPOINT"),
		LoadBalancerEndpoint: os.Getenv("OS_LOADBALANCER_ENDPOINT"),
		AuthEndpoint:         os.Getenv("OS_AUTH_ENDPOINT"),
		AuthToken:            os.Getenv("OS_TOKEN"),
		ProjectId:            os.Getenv("OS_PROJECT_ID"),
	}

	os.Setenv("OS_TOKEN", "")

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
	_, err := cl.Authenticate(credential)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(credential)

	t.Run("get existing load balancer", func(t *testing.T) {
		res, err := cl.GetLoadBalancer("59107a73-6f12-4a8d-9812-383b3e66977a")
		if err != nil {
			t.Fatal(error.Error(err))
		}
		t.Log(res)
	})

	t.Run("get non-existing load balancer", func(t *testing.T) {
		res, err := cl.GetLoadBalancer("xxxyyz")
		if err != nil {
			t.Fatal(error.Error(err))
		}
		t.Log(res)
	})
}

// go test ./test -v -run ^TestMagnumClientPasswordAuthentication$
func TestMagnumClientPasswordAuthentication(t *testing.T) {
	cl := api.OpenstackClient{
		MagnumEndpoint:  os.Getenv("VT_MAGNUM_ENDPOINT"),
		NetworkEndpoint: os.Getenv("VT_NETWORK_ENDPOINT"),
		AuthEndpoint:    os.Getenv("VT_AUTH_ENDPOINT"),
		AuthToken:       os.Getenv("OS_TOKEN"),
		ProjectId:       os.Getenv("VT_PROJECT_ID"),
	}

	credential := api.OpenstackAuth{
		Identity: api.OpenstackIdentity{
			Methods: []string{"password"},
			Password: api.OpenstackPassword{
				User: map[string]interface{}{
					"name":     os.Getenv("VT_USERNAME"),
					"password": os.Getenv("VT_PASSWORD"),
					"domain": map[string]string{
						"id": os.Getenv("VT_DOMAIN_ID"),
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

// go test ./test -v -run ^TestMagnumListClusters$
func TestMagnumListClusters(t *testing.T) {
	cl := api.OpenstackClient{
		MagnumEndpoint:  os.Getenv("VT_MAGNUM_ENDPOINT"),
		NetworkEndpoint: os.Getenv("VT_NETWORK_ENDPOINT"),
		AuthEndpoint:    os.Getenv("VT_AUTH_ENDPOINT"),
		AuthToken:       os.Getenv("OS_TOKEN"),
		ProjectId:       os.Getenv("VT_PROJECT_ID"),
	}

	credential := api.OpenstackAuth{
		Identity: api.OpenstackIdentity{
			Methods: []string{"password"},
			Password: api.OpenstackPassword{
				User: map[string]interface{}{
					"name":     os.Getenv("VT_USERNAME"),
					"password": os.Getenv("VT_PASSWORD"),
					"domain": map[string]string{
						"id": os.Getenv("VT_DOMAIN_ID"),
					},
				},
			},
		},
	}

	_, err := cl.Authenticate(credential)
	if err != nil {
		t.Fatal(err)
	}

	b, err := cl.MagnumListClusters()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(b))
}

// go test ./test -v -run ^TestMagnumCreateCluster$
func TestMagnumCreateCluster(t *testing.T) {
	cl := api.OpenstackClient{
		MagnumEndpoint:  os.Getenv("VT_MAGNUM_ENDPOINT"),
		NetworkEndpoint: os.Getenv("VT_NETWORK_ENDPOINT"),
		AuthEndpoint:    os.Getenv("VT_AUTH_ENDPOINT"),
		AuthToken:       os.Getenv("OS_TOKEN"),
		ProjectId:       os.Getenv("VT_PROJECT_ID"),
	}

	projectName := os.Getenv("VT_PROJECT_NAME")
	credential := api.OpenstackAuth{
		Identity: api.OpenstackIdentity{
			Methods: []string{"password"},
			Password: api.OpenstackPassword{
				User: map[string]interface{}{
					"name":     os.Getenv("VT_USERNAME"),
					"password": os.Getenv("VT_PASSWORD"),
					"domain": map[string]string{
						"id": os.Getenv("VT_DOMAIN_ID"),
					},
				},
			},
		},
		Scope: &api.OpenstackScope{
			Project: &api.OpenstackProject{
				Name: &projectName,
				Domain: &map[string]interface{}{
					"id": os.Getenv("VT_DOMAIN_ID"),
				},
			},
		},
	}

	_, err := cl.Authenticate(credential)
	if err != nil {
		t.Fatal(err)
	}

	labels := map[string]interface{}{
		"etcd_lb_disabled":         "true",
		"csi_attacher_tag":         "v3.1.0",
		"hyperkube_image":          "docker.io/virtuozzo/hci-binary-hyperkube",
		"availability_zone":        "nova",
		"csi_snapshotter_tag":      "v3.0.3",
		"cloud_provider_tag":       "v1.22.0",
		"etcd_tag":                 "v3.4.6",
		"docker_volume_type":       "Standard_vDisk-T0",                    // <-
		"octavia_api_lb_flavor":    "13eda90d-a202-46d9-a2f4-5b03c1d440d0", // ?
		"cgroup_driver":            "systemd",
		"cinder_csi_enabled":       "true",
		"kube_version":             "v1.23.5", // <- kubernetes version
		"kube_tag":                 "v1.23.5", // <- kubernetes version
		"use_podman":               "true",
		"boot_volume_type":         "Standard_vDisk-T0", // <-
		"cinder_csi_plugin_tag":    "v1.22.0",
		"flannel_tag":              "v0.11.0-amd64",
		"boot_volume_size":         "10", // <-
		"heat_container_agent_tag": "5.3.11",
		"octavia_default_flavor":   "13eda90d-a202-46d9-a2f4-5b03c1d440d0", // ?
		"cloud_provider_enabled":   "true",
	}

	args := model.MagnumCreateClusterRequest{
		Name:              "lyrid-test-go",
		MasterFlavorID:    "a2.medium-2",
		MasterCount:       1,
		FlavorID:          "a2.medium-2",
		NodeCount:         1,
		Keypair:           "eranya-ssh",
		DockerVolumeSize:  10,
		ClusterTemplateID: "8764c702-dd08-41fd-9b55-d1147f0144dd",
		CreateTimeout:     120,
		Labels:            labels,
	}

	x, _ := json.Marshal(args)
	t.Log(string(x))

	b, err := cl.MagnumCreateCluster(args)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(b))
}

// go test ./test -v -run ^TestMagnumListClusterTemplates$
func TestMagnumListClusterTemplates(t *testing.T) {
	cl := api.OpenstackClient{
		MagnumEndpoint:  os.Getenv("VT_MAGNUM_ENDPOINT"),
		NetworkEndpoint: os.Getenv("VT_NETWORK_ENDPOINT"),
		AuthEndpoint:    os.Getenv("VT_AUTH_ENDPOINT"),
		AuthToken:       os.Getenv("OS_TOKEN"),
		ProjectId:       os.Getenv("VT_PROJECT_ID"),
	}

	projectName := os.Getenv("VT_PROJECT_NAME")
	credential := api.OpenstackAuth{
		Identity: api.OpenstackIdentity{
			Methods: []string{"password"},
			Password: api.OpenstackPassword{
				User: map[string]interface{}{
					"name":     os.Getenv("VT_USERNAME"),
					"password": os.Getenv("VT_PASSWORD"),
					"domain": map[string]string{
						"id": os.Getenv("VT_DOMAIN_ID"),
					},
				},
			},
		},
		Scope: &api.OpenstackScope{
			Project: &api.OpenstackProject{
				Name: &projectName,
				Domain: &map[string]interface{}{
					"id": os.Getenv("VT_DOMAIN_ID"),
				},
			},
		},
	}

	_, err := cl.Authenticate(credential)
	if err != nil {
		t.Fatal(err)
	}

	b, err := cl.MagnumListClusterTemplates()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(b))
}

// go test ./test -v -run ^TestMagnumCreateClusterTemplate$
func TestMagnumCreateClusterTemplate(t *testing.T) {
	os.Setenv("OS_TOKEN", "")
	cl := api.OpenstackClient{
		MagnumEndpoint:  os.Getenv("VT_MAGNUM_ENDPOINT"),
		NetworkEndpoint: os.Getenv("VT_NETWORK_ENDPOINT"),
		AuthEndpoint:    os.Getenv("VT_AUTH_ENDPOINT"),
		AuthToken:       os.Getenv("OS_TOKEN"),
		ProjectId:       os.Getenv("VT_PROJECT_ID"),
	}

	projectName := os.Getenv("VT_PROJECT_NAME")
	credential := api.OpenstackAuth{
		Identity: api.OpenstackIdentity{
			Methods: []string{"password"},
			Password: api.OpenstackPassword{
				User: map[string]interface{}{
					"name":     os.Getenv("VT_USERNAME"),
					"password": os.Getenv("VT_PASSWORD"),
					"domain": map[string]string{
						"id": os.Getenv("VT_DOMAIN_ID"),
					},
				},
			},
		},
		Scope: &api.OpenstackScope{
			Project: &api.OpenstackProject{
				Name: &projectName,
				Domain: &map[string]interface{}{
					"id": os.Getenv("VT_DOMAIN_ID"),
				},
			},
		},
	}

	_, err := cl.Authenticate(credential)
	if err != nil {
		t.Fatal(err)
	}

	b, err := cl.MagnumListClusterTemplates()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(b))

	args := model.MagnumCreateClusterTemplateRequest{
		FloatingIPEnabled:   true,
		DockerVolumeSize:    10,
		ServerType:          "vm",
		ExternalNetworkID:   "f30c9e3d-757b-43fb-b4e0-da3ab36708a4",
		ImageID:             "f4c9aa27-fe3a-4c41-819a-67bdf75baa9f", // get from "openstack image list"
		ClusterDistro:       "fedora-coreos",
		VolumeDriver:        "cinder",
		DockerStorageDriver: "overlay2",
		Name:                "test-lyrid-wqq_template",
		NetworkDriver:       "flannel",
		FixedNetwork:        "a7e5dec7-44e3-455e-9ef8-e30af85328e8", //private - get from "openstack network list"
		FixedSubnet:         "a28ce377-5812-4f09-91a4-dda327102543", // get from "openstack subnet list"
		COE:                 "kubernetes",
		MasterLBEnabled:     false,
	}

	x, _ := json.Marshal(args)
	t.Log(string(x))

	bb, err := cl.MagnumCreateClusterTemplate(args)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(bb))
}

// go test ./test -v -run ^TestImageList$
func TestImageList(t *testing.T) {
	cl := api.OpenstackClient{
		MagnumEndpoint:  os.Getenv("VT_MAGNUM_ENDPOINT"),
		NetworkEndpoint: os.Getenv("VT_NETWORK_ENDPOINT"),
		AuthEndpoint:    os.Getenv("VT_AUTH_ENDPOINT"),
		ImageEndpoint:   os.Getenv("VT_IMAGE_ENDPOINT"),
		AuthToken:       os.Getenv("OS_TOKEN"),
		ProjectId:       os.Getenv("VT_PROJECT_ID"),
	}

	projectName := os.Getenv("VT_PROJECT_NAME")
	credential := api.OpenstackAuth{
		Identity: api.OpenstackIdentity{
			Methods: []string{"password"},
			Password: api.OpenstackPassword{
				User: map[string]interface{}{
					"name":     os.Getenv("VT_USERNAME"),
					"password": os.Getenv("VT_PASSWORD"),
					"domain": map[string]string{
						"id": os.Getenv("VT_DOMAIN_ID"),
					},
				},
			},
		},
		Scope: &api.OpenstackScope{
			Project: &api.OpenstackProject{
				Name: &projectName,
				Domain: &map[string]interface{}{
					"id": os.Getenv("VT_DOMAIN_ID"),
				},
			},
		},
	}

	_, err := cl.Authenticate(credential)
	if err != nil {
		t.Fatal(err)
	}

	response, err := cl.GetImageList(nil)
	if err != nil {
		t.Fatal(err)
	}

	b, _ := json.Marshal(response)
	t.Log(string(b))

	filter, _ := utils.GetQueryStringFromURL(response.Next)
	t.Log(filter)

	response, err = cl.GetImageList(filter) // get result from next page
	if err != nil {
		t.Fatal(err)
	}

	b, _ = json.Marshal(response)
	t.Log(string(b))

	response, err = cl.GetImageList(map[string]string{"name": "fedora-coreos-x64-k8saas"}) // get image by name filter
	if err != nil {
		t.Fatal(err)
	}

	b, _ = json.Marshal(response)
	t.Log(string(b))

	resp, err := cl.GetImageByID("f4c9aa27-fe3a-4c41-819a-67bdf75baa9f")
	if err != nil {
		t.Fatal(err)
	}

	b, _ = json.Marshal(resp)
	t.Log(string(b))
}

// go test ./test -v -run ^TestNetworkList$
func TestNetworkList(t *testing.T) {
	cl := api.OpenstackClient{
		MagnumEndpoint:  os.Getenv("VT_MAGNUM_ENDPOINT"),
		NetworkEndpoint: os.Getenv("VT_NETWORK_ENDPOINT"),
		AuthEndpoint:    os.Getenv("VT_AUTH_ENDPOINT"),
		ImageEndpoint:   os.Getenv("VT_IMAGE_ENDPOINT"),
		AuthToken:       os.Getenv("OS_TOKEN"),
		ProjectId:       os.Getenv("VT_PROJECT_ID"),
	}

	projectName := os.Getenv("VT_PROJECT_NAME")
	credential := api.OpenstackAuth{
		Identity: api.OpenstackIdentity{
			Methods: []string{"password"},
			Password: api.OpenstackPassword{
				User: map[string]interface{}{
					"name":     os.Getenv("VT_USERNAME"),
					"password": os.Getenv("VT_PASSWORD"),
					"domain": map[string]string{
						"id": os.Getenv("VT_DOMAIN_ID"),
					},
				},
			},
		},
		Scope: &api.OpenstackScope{
			Project: &api.OpenstackProject{
				Name: &projectName,
				Domain: &map[string]interface{}{
					"id": os.Getenv("VT_DOMAIN_ID"),
				},
			},
		},
	}

	_, err := cl.Authenticate(credential)
	if err != nil {
		t.Fatal(err)
	}

	response, err := cl.GetNetworkList(nil)
	if err != nil {
		t.Fatal(err)
	}

	b, _ := json.Marshal(response)
	t.Log(string(b))

	resp, err := cl.GetNetworkByID("a7e5dec7-44e3-455e-9ef8-e30af85328e8")
	if err != nil {
		t.Fatal(err)
	}

	b, _ = json.Marshal(resp)
	t.Log(string(b))
}

// go test ./test -v -run ^TestSubnetList$
func TestSubnetList(t *testing.T) {
	cl := api.OpenstackClient{
		MagnumEndpoint:  os.Getenv("VT_MAGNUM_ENDPOINT"),
		NetworkEndpoint: os.Getenv("VT_NETWORK_ENDPOINT"),
		AuthEndpoint:    os.Getenv("VT_AUTH_ENDPOINT"),
		ImageEndpoint:   os.Getenv("VT_IMAGE_ENDPOINT"),
		AuthToken:       os.Getenv("OS_TOKEN"),
		ProjectId:       os.Getenv("VT_PROJECT_ID"),
	}

	projectName := os.Getenv("VT_PROJECT_NAME")
	credential := api.OpenstackAuth{
		Identity: api.OpenstackIdentity{
			Methods: []string{"password"},
			Password: api.OpenstackPassword{
				User: map[string]interface{}{
					"name":     os.Getenv("VT_USERNAME"),
					"password": os.Getenv("VT_PASSWORD"),
					"domain": map[string]string{
						"id": os.Getenv("VT_DOMAIN_ID"),
					},
				},
			},
		},
		Scope: &api.OpenstackScope{
			Project: &api.OpenstackProject{
				Name: &projectName,
				Domain: &map[string]interface{}{
					"id": os.Getenv("VT_DOMAIN_ID"),
				},
			},
		},
	}

	_, err := cl.Authenticate(credential)
	if err != nil {
		t.Fatal(err)
	}

	response, err := cl.GetSubnetList(nil)
	if err != nil {
		t.Fatal(err)
	}

	b, _ := json.Marshal(response)
	t.Log(string(b))

	response, err = cl.GetSubnetList(map[string]string{"network_id": "a7e5dec7-44e3-455e-9ef8-e30af85328e8"}) // get subnet by network id
	if err != nil {
		t.Fatal(err)
	}

	b, _ = json.Marshal(response)
	t.Log(string(b))

	resp, err := cl.GetSubnetByID("8f2acb4a-13bf-4e34-a0fc-fc22980e918d")
	if err != nil {
		t.Fatal(err)
	}

	b, _ = json.Marshal(resp)
	t.Log(string(b))
}
