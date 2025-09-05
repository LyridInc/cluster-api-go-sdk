package test

import (
	"os"
	"testing"

	svcapi "github.com/LyridInc/cluster-api-go-sdk/api/servercore"
	svcmodel "github.com/LyridInc/cluster-api-go-sdk/model/servercore"
)

// export $(< ./test/data/test.env)

func setupClient(t *testing.T) (*svcapi.ServercoreClient, string) {
	t.Helper() // marks function as test helper

	cl := svcapi.NewServercoreClient(svcmodel.Config{
		ApiKey:                  os.Getenv("SERVERCORE_API_KEY"),
		ApiUrl:                  os.Getenv("SERVERCORE_API_URL"),
		CloudApiUrl:             os.Getenv("SERVERCORE_CLOUD_API_URL"),
		ManagedKubernetesApiUrl: os.Getenv("SERVERCORE_MANAGED_KUBERNETES_API_URL"),
		ProjectID:               os.Getenv("SERVERCORE_PROJECT_ID"),
	})

	_, token, err := cl.Authenticate(svcmodel.AuthConfig{
		Username:    os.Getenv("SERVERCORE_USER_NAME"),
		Password:    os.Getenv("SERVERCORE_USER_PASSWORD"),
		Domain:      os.Getenv("SERVERCORE_DOMAIN_NAME"),
		ProjectName: "Lyrid Development",
	})
	if err != nil {
		t.Fatal(err)
	}

	return cl, token
}

// go test ./test -v -run ^TestAuthenticateServercore$
func TestAuthenticateServercore(t *testing.T) {
	_, token := setupClient(t)
	t.Log(token)
}

// go test ./test -v -run ^TestListServercoreClusters$
func TestListServercoreClusters(t *testing.T) {
	cl, token := setupClient(t)

	t.Log(token)

	clusterListResp, err := cl.GetListClusters()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(clusterListResp)
}

// go test ./test -v -run ^TestGetServercoreClusterByID$
func TestGetServercoreClusterByID(t *testing.T) {
	cl, token := setupClient(t)

	t.Log(token)

	clusterResp, err := cl.GetClusterByID("b78ed91d-9822-4d02-8896-06aecc880f42")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(clusterResp)
}

// go test ./test -v -run ^TestCreateServercoreCluster$
func TestCreateServercoreCluster(t *testing.T) {
	cl, token := setupClient(t)

	t.Log(token)

	createClusterResp, err := cl.CreateCluster(svcmodel.CreateClusterRequest{
		Cluster: svcmodel.ClusterRequest{
			KubeVersion:    "1.32.2",
			Name:           "test-from-go",
			PrivateKubeAPI: false,
			Region:         "ke-1",
			Zonal:          true,
			NodeGroups: []svcmodel.NodeGroupRequest{
				{
					Count:            1,
					CPUs:             2,
					RAMMB:            4096,
					AvailabilityZone: "ke-1a",
					VolumeGB:         50,
					VolumeType:       "universal2.ke-1a",
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(createClusterResp)
}

// go test ./test -v -run ^TestDeleteServercoreClusterByID$
func TestDeleteServercoreClusterByID(t *testing.T) {
	cl, token := setupClient(t)

	t.Log(token)

	clusterResp, err := cl.DeleteClusterByID("687a4e55-4b69-49e4-a365-688a3cfb3651")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(clusterResp)
}

// go test ./test -v -run ^TestGetServercoreClusterKubeconfig$
func TestGetServercoreClusterKubeconfig(t *testing.T) {
	cl, token := setupClient(t)

	t.Log(token)

	resp, err := cl.GetClusterKubeconfig("687a4e55-4b69-49e4-a365-688a3cfb3651")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(*resp)
}

// go test ./test -v -run ^TestListServercoreIPs$
func TestListServercoreIPs(t *testing.T) {
	cl, token := setupClient(t)

	t.Log(token)

	ipListResp, err := cl.GetListIPs()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(ipListResp)
}

// go test ./test -v -run ^TestListServercoreFloatingIPs$
func TestListServercoreFloatingIPs(t *testing.T) {
	cl, token := setupClient(t)

	t.Log(token)

	ipListResp, err := cl.GetListFloatingIPs()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(ipListResp)
}

// go test ./test -v -run ^TestListServercoreFloatingIPDetail$
func TestListServercoreFloatingIPDetail(t *testing.T) {
	cl, token := setupClient(t)

	t.Log(token)

	ipResp, err := cl.GetListFloatingIPDetail("eaef37c1-7c63-4079-837d-a0dea9a0e2d8")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(ipResp)
}

// go test ./test -v -run ^TestCreateServercoreFloatingIP$
func TestCreateServercoreFloatingIP(t *testing.T) {
	cl, token := setupClient(t)

	t.Log(token)

	ipResp, err := cl.CreateFloatingIP(svcmodel.CreateFloatingIPRequest{
		FloatingIPs: []svcmodel.FloatingIPRequest{
			{
				Quantity: 1,
				Region:   "ke-1",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(ipResp)
}

// go test ./test -v -run ^TestDeleteServercoreFloatingIP$
func TestDeleteServercoreFloatingIP(t *testing.T) {
	cl, token := setupClient(t)

	t.Log(token)

	ipResp, err := cl.DeleteFloatingIPByID("ef20e389-5df3-4827-8075-778fb014c738")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(ipResp)
}
