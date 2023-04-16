package test

import (
	"os"
	"testing"

	"github.com/LyridInc/cluster-api-go-sdk/api"
	"github.com/LyridInc/cluster-api-go-sdk/model"
)

// export $(< ./test/data/test.env)

// go test ./test -v -run ^TestVirtuozzoAuthentication$
func TestVirtuozzoAuthentication(t *testing.T) {
	authEndpoint := os.Getenv("VT_AUTH_ENDPOINT")
	username := os.Getenv("VT_USERNAME")
	project := os.Getenv("VT_PROJECT_NAME")
	domainId := os.Getenv("VT_DOMAIN_ID")
	password := os.Getenv("VT_PASSWORD")

	vClient := api.NewVirtuozzoClient(authEndpoint)
	response, err := vClient.Authenticate(project, username, password, domainId)
	if err != nil {
		t.Fatal(error.Error(err))
	}

	t.Log(string(response))

	resp, err := vClient.ListPublicEndpoints()
	if err != nil {
		t.Fatal(error.Error(err))
	}

	t.Log(string(resp))
}

// go test ./test -v -run ^TestVirtuozzoListKubernetesClusterTemplates$
func TestVirtuozzoListKubernetesClusterTemplates(t *testing.T) {
	authEndpoint := os.Getenv("VT_AUTH_ENDPOINT")
	username := os.Getenv("VT_USERNAME")
	project := os.Getenv("VT_PROJECT_NAME")
	domainId := os.Getenv("VT_DOMAIN_ID")
	password := os.Getenv("VT_PASSWORD")

	vClient := api.NewVirtuozzoClient(authEndpoint)
	_, err := vClient.Authenticate(project, username, password, domainId)
	if err != nil {
		t.Fatal(error.Error(err))
	}

	resp, err := vClient.ListKubernetesClusterTemplates()
	if err != nil {
		t.Fatal(error.Error(err))
	}

	t.Log(string(resp))
}

// go test ./test -v -run ^TestVirtuozzoCreateKubernetesClusterTemplates$
func TestVirtuozzoCreateKubernetesClusterTemplates(t *testing.T) {
	authEndpoint := os.Getenv("VT_AUTH_ENDPOINT")
	username := os.Getenv("VT_USERNAME")
	project := os.Getenv("VT_PROJECT_NAME")
	domainId := os.Getenv("VT_DOMAIN_ID")
	password := os.Getenv("VT_PASSWORD")

	vClient := api.NewVirtuozzoClient(authEndpoint)
	_, err := vClient.Authenticate(project, username, password, domainId)
	if err != nil {
		t.Fatal(error.Error(err))
	}

	// {"name":"vt-test","version":"v1.23.5","master_node_count":1,"master_flavor":"a2.large-2","containers_volume_storage_policy":"Standard_vDisk-T0","containers_volume_size":10,"key_name":"eranya-ssh","external_network_id":"f30c9e3d-757b-43fb-b4e0-da3ab36708a4","network_id":"f30c9e3d-757b-43fb-b4e0-da3ab36708a4","worker_pools":[{"flavor":"a2.large-2","node_count":1}],"floating_ip_enabled":false,"public_access_enabled":false}

	args := model.VirtuozzoCreateKubernetesClusterTemplateArgs{
		FloatingIPEnabled: true,
		MasterFlavorID:    "a2.large-2",
		FixedNetwork:      "f30c9e3d-757b-43fb-b4e0-da3ab36708a4", // public network
		FixedSubnet:       "103.176.44.0/23",
		HttpsProxy:        "None",
		NoProxy:           "None",
		TLSDisabled:       true,
		KeypairID:         "eranya-ssh",
		Public:            true,
		DockerVolumeSize:  15,
		ServerType:        "vm",
		ExternalNetworkID: "f30c9e3d-757b-43fb-b4e0-da3ab36708a4",
		ImageID:           "b35b2de1-71a1-443c-83cb-1449d358e415",
		// VolumeDriver:        "csi",
		RegistryEnabled:     true,
		DockerStorageDriver: "overlay2",
		Name:                "vt-development",
		NetworkDriver:       "flannel",
		COE:                 "kubernetes", // swarm-mode, swarm, mesos, kubernetes, dcos
		FlavorID:            "a2.large-2",
		MasterLBEnabled:     true,
		DNSNameserver:       "8.8.8.8",
	}

	resp, err := vClient.CreateKubernetesClusterTemplate(args)
	if err != nil {
		t.Fatal(error.Error(err))
	}

	t.Log(string(resp))
}
