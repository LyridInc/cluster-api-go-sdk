package test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/LyridInc/cluster-api-go-sdk/api"
	"github.com/LyridInc/cluster-api-go-sdk/model"
)

// go test ./test -v -run ^TestCreateCivoCluster$
func TestCreateCivoCluster(t *testing.T) {
	token := os.Getenv("CIVO_TOKEN")
	endpoint := os.Getenv("CIVO_API_ENDPOINT")
	client := api.NewCivoClient(token, endpoint)

	args := api.CivoCreateClusterArgs{
		Name:      "lyrid-sdk",
		NetworkID: "8674acc1-2fcd-4880-b62c-4605f5fe578d",
		Region:    "NYC1",
		CNIPlugin: "flannel",
		Pools: []model.CivoPool{
			{
				ID:    "sdk-pool-1",
				Count: 1,
				Size:  "g4s.kube.medium",
			},
		},
		KubernetesVersion: "1.29.8-k3s1 ",
		InstanceFirewall:  "6cbf5e4c-6256-4f37-80b2-7cdab7d0ac1c",
	}

	res, err := client.CreateCluster(args)
	if err != nil {
		t.Fatal(err)
	}

	b, _ := json.Marshal(res)
	// eyJjb2RlIjoia3ViZXJuZXRlc19pbnZhbGlkX3NpemUiLCJyZWFzb24iOiJGYWlsZWQgdG8gY3JlYXRlIGEga3ViZXJuZXRlcyBjbHVzdGVyIHdpdGggdGhlIGdpdmVuIHNpemUuIFBsZWFzZSBzZWxlY3QgYSB2YWxpZCBzaXplIG9mIFR5cGU6IEt1YmVybmV0ZXMsIHlvdSBjYW4gbGlzdCBhbGwgc2l6ZXMgd2l0aCBDTEkgOiBgY2l2byBzaXplcyBscyJ9Cg==

	// {
	// 	"id": "33030763-1b6e-44bb-9399-6c14932c5a44",
	// 	"name": "lyrid-sdk",
	// 	"version": "1.28.7-k3s1",
	// 	"cluster_type": "k3s",
	// 	"status": "BUILDING",
	// 	"ready": false,
	// 	"num_target_nodes": 1,
	// 	"target_nodes_size": "",
	// 	"built_at": "0001-01-01T00:00:00Z",
	// 	"kubernetes_version": "1.28.7-k3s1",
	// 	"api_endpoint": "https://:6443",
	// 	"dns_entry": "33030763-1b6e-44bb-9399-6c14932c5a44.k8s.civo.com",
	// 	"created_at": "2025-03-26T21:04:28Z",
	// 	"master_ip": "",
	// 	"pools": null,
	// 	"required_pools": [
	// 		{
	// 			"id": "sdk-pool-1",
	// 			"size": "g4s.kube.medium",
	// 			"count": 1,
	// 			"labels": null,
	// 			"taints": null
	// 		}
	// 	],
	// 	"firewall_id": "6cbf5e4c-6256-4f37-80b2-7cdab7d0ac1c",
	// 	"master_ipv6": "",
	// 	"network_id": "8674acc1-2fcd-4880-b62c-4605f5fe578d",
	// 	"namespace": "cust-default-b51ba92a-75059924730c",
	// 	"size": "",
	// 	"count": 0,
	// 	"kubeconfig": null,
	// 	"instances": null,
	// 	"installed_applications": null,
	// 	"conditions": [
	// 		{
	// 			"type": "ClusterVersionSync",
	// 			"status": "Unknown",
	// 			"synced": true,
	// 			"last_transition_time": null
	// 		},
	// 		...
	// 	],
	// 	"cni_plugin": "flannel",
	// 	"ccm_installed": "true"
	// }

	t.Log(string(b))
}

// go test ./test -v -run ^TestListInstanceSizes$
func TestListInstanceSizes(t *testing.T) {
	token := os.Getenv("CIVO_TOKEN")
	endpoint := os.Getenv("CIVO_API_ENDPOINT")
	client := api.NewCivoClient(token, endpoint)

	res, err := client.ListInstanceSizes(nil)
	if err != nil {
		t.Fatal(err)
	}

	b, _ := json.Marshal(res)
	// {
	//   "name": "g3.medium",
	//   "nice_name": "Medium",
	//   "cpu_cores": 2,
	//   "ram_mb": 4096,
	//   "disk_gb": 50,
	//   "description": "Medium",
	//   "selectable": true
	// },

	t.Log(string(b))
}

// go test ./test -v -run ^TestListNetworks$
func TestListNetworks(t *testing.T) {
	token := os.Getenv("CIVO_TOKEN")
	endpoint := os.Getenv("CIVO_API_ENDPOINT")
	client := api.NewCivoClient(token, endpoint)

	res, err := client.ListNetworks(nil)
	if err != nil {
		t.Fatal(err)
	}

	b, _ := json.Marshal(res)
	// [{"result":"","label":"default","id":"8674acc1-2fcd-4880-b62c-4605f5fe578d"}]

	t.Log(string(b))
}

// go test ./test -v -run ^TestListKubernetesVersions$
func TestListKubernetesVersions(t *testing.T) {
	token := os.Getenv("CIVO_TOKEN")
	endpoint := os.Getenv("CIVO_API_ENDPOINT")
	client := api.NewCivoClient(token, endpoint)

	res, err := client.ListKubernetesVersions(nil)
	if err != nil {
		t.Fatal(err)
	}

	b, _ := json.Marshal(res)
	// {
	//   "version": "1.28.7-k3s1",
	//   "flat_version": "",
	//   "label": "1.28.7-k3s1",
	//   "type": "stable",
	//   "default": true
	// },

	t.Log(string(b))
}

// go test ./test -v -run ^TestListFirewalls$
func TestListFirewalls(t *testing.T) {
	token := os.Getenv("CIVO_TOKEN")
	endpoint := os.Getenv("CIVO_API_ENDPOINT")
	client := api.NewCivoClient(token, endpoint)

	res, err := client.ListFirewalls(nil)
	if err != nil {
		t.Fatal(err)
	}

	b, _ := json.Marshal(res)
	// [{"id":"6cbf5e4c-6256-4f37-80b2-7cdab7d0ac1c","name":"default-default","account_id":"b51ba92a-3537-4257-b7b6-5181684b17cf","rules_count":12,"default":"true","label":"","network_id":"8674acc1-2fcd-4880-b62c-4605f5fe578d"}]

	t.Log(string(b))
}

// go test ./test -v -run ^TestListFirewallRules$
func TestListFirewallRules(t *testing.T) {
	token := os.Getenv("CIVO_TOKEN")
	endpoint := os.Getenv("CIVO_API_ENDPOINT")
	client := api.NewCivoClient(token, endpoint)

	res, err := client.ListFirewallRules("6cbf5e4c-6256-4f37-80b2-7cdab7d0ac1c")
	if err != nil {
		t.Fatal(err)
	}

	b, _ := json.Marshal(res)

	t.Log(string(b))
}

// go test ./test -v -run ^TestGetCivoClusterDetail$
func TestGetCivoClusterDetail(t *testing.T) {
	token := os.Getenv("CIVO_TOKEN")
	endpoint := os.Getenv("CIVO_API_ENDPOINT")
	client := api.NewCivoClient(token, endpoint)

	res, err := client.GetClusterDetail("33030763-1b6e-44bb-9399-6c14932c5a44")
	if err != nil {
		t.Fatal(err)
	}

	b, _ := json.Marshal(res)

	t.Log(string(b))
}
