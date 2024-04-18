package test

import (
	"encoding/base64"
	"fmt"
	"os"
	"testing"

	"github.com/LyridInc/cluster-api-go-sdk/api"
	"github.com/LyridInc/cluster-api-go-sdk/option"
	"github.com/apache/cloudstack-go/cloudstack"
)

// go test ./test -v -run ^TestGetListPublicIPAddresses$
func TestGetListPublicIPAddresses(t *testing.T) {
	apiUrl := os.Getenv("CLOUDSTACK_API_URL")
	apiKey := os.Getenv("CLOUDSTACK_API_KEY")
	secret := os.Getenv("CLOUDSTACK_SECRET")
	cs := cloudstack.NewAsyncClient(apiUrl, apiKey, secret, true)

	params := cloudstack.ListPublicIpAddressesParams{}
	zoneId := os.Getenv("CLOUDSTACK_ZONE_ID")
	params.SetZoneid(zoneId)
	params.SetState("Free")

	resp, err := cs.Address.ListPublicIpAddresses(&params)
	if err != nil {
		t.Fatal(err)
	}

	ipAddresses := []string{}
	for _, ip := range resp.PublicIpAddresses {
		ipAddresses = append(ipAddresses, ip.Ipaddress)
	}
	t.Log(ipAddresses)
	t.Log(len(ipAddresses))
}

// go test ./test -v -run ^TestGenerateCloudStackClusterTemplate$
func TestGenerateCloudStackClusterTemplate(t *testing.T) {
	infrastructure := "cloudstack"
	capi, _ := api.NewClusterApiClient("", "./data/az-vega.kubeconfig")

	ready, err := capi.InfrastructureReadiness(infrastructure)
	if !ready && err == nil {
		t.Log("initialize infrastructure")
		capi.InitInfrastructure(infrastructure)
	}

	// get public IP
	apiUrl := os.Getenv("CLOUDSTACK_API_URL")
	apiKey := os.Getenv("CLOUDSTACK_API_KEY")
	secret := os.Getenv("CLOUDSTACK_SECRET")
	cs := cloudstack.NewAsyncClient(apiUrl, apiKey, secret, true)

	params := cloudstack.ListPublicIpAddressesParams{}
	zoneId := os.Getenv("CLOUDSTACK_ZONE_ID")
	params.SetZoneid(zoneId)
	params.SetState("Free")

	resp, err := cs.Address.ListPublicIpAddresses(&params)
	if err != nil {
		t.Fatal(err)
	}

	if len(resp.PublicIpAddresses) <= 0 {
		t.Fatal("No Public IP is available")
	}

	ipAddress := resp.PublicIpAddresses[0].Ipaddress
	fmt.Println("Cluster Endpoint IP:", ipAddress)

	zone, _, err := cs.Zone.GetZoneByID(zoneId)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Zone:", zone.Name)

	clusterName := "capc-test"
	clusterOpt := option.GenerateCloudStackWorkloadClusterOption{
		ZoneName:                    zone.Name,
		NetworkName:                 fmt.Sprintf("%s-network", clusterName),
		ClusterEndpointIP:           ipAddress, // TODO
		ClusterEndpointPort:         "6443",
		ControlPlaneMachineOffering: "DBaaS Premium Control Node 2/4/100",
		WorkerMachineOffering:       "DBaaS Premium Worker Node 16/32/100",
		TemplateName:                "KUBE Ubuntu 20.04 CAPI",
		SshKeyName:                  "azhary-keypair",
		ClusterName:                 clusterName,
		Namespace:                   "capc",
		KubernetesVersion:           "v1.27.3",
		WorkerMachineCount:          1,
		ControlPlaneMachineCount:    1,
		URL:                         "./data/cloudstack/cluster-template.yaml",
	}
	yaml, err := capi.GenerateCloudStackWorkloadClusterYaml(clusterOpt)
	if err != nil {
		t.Fatal("Generate workload cluster error:", err)
	}

	if err := os.WriteFile(fmt.Sprintf("./data/%s.yaml", clusterName), []byte(yaml), 0644); err != nil {
		t.Fatal("Write yaml error:", err)
	}
}

// go test ./test -v -run ^TestInstallCloudStackCloudControllerManager$
func TestInstallCloudStackCloudControllerManager(t *testing.T) {
	apiUrl := os.Getenv("CLOUDSTACK_API_URL")
	apiKey := os.Getenv("CLOUDSTACK_API_KEY")
	secret := os.Getenv("CLOUDSTACK_SECRET")
	zoneId := os.Getenv("CLOUDSTACK_ZONE_ID")

	cs := cloudstack.NewAsyncClient(apiUrl, apiKey, secret, true)
	zone, _, err := cs.Zone.GetZoneByID(zoneId)
	if err != nil {
		t.Fatal("error GetZoneByID", err)
	}

	cloudConfigValues := "[Global]\n"
	cloudConfigValues = cloudConfigValues + fmt.Sprintf("api-url = %s\n", apiUrl)
	cloudConfigValues = cloudConfigValues + fmt.Sprintf("api-key = %s\n", apiKey)
	cloudConfigValues = cloudConfigValues + fmt.Sprintf("secret-key = %s\n", secret)
	cloudConfigValues = cloudConfigValues + fmt.Sprintf("zone = %s\n", zone.Name)
	cloudConfigValues = cloudConfigValues + fmt.Sprintf("ssl-no-verify = %s\n", "false")

	cloudConfigValuesB64 := base64.StdEncoding.EncodeToString([]byte(cloudConfigValues))
	t.Log(cloudConfigValues)
	t.Log(cloudConfigValuesB64)
}

// go test ./test -v -run ^TestCreateCloudStackSSHKeypair$
func TestCreateCloudStackSSHKeypair(t *testing.T) {
	apiUrl := os.Getenv("CLOUDSTACK_API_URL")
	apiKey := os.Getenv("CLOUDSTACK_API_KEY")
	secret := os.Getenv("CLOUDSTACK_SECRET")

	cs := cloudstack.NewAsyncClient(apiUrl, apiKey, secret, true)
	params := cloudstack.CreateSSHKeyPairParams{}
	params.SetName("test-keypair")
	resp, err := cs.SSH.CreateSSHKeyPair(&params)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp.Privatekey)
}
