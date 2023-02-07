package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/LyridInc/cluster-api-go-sdk/api"
	"github.com/LyridInc/cluster-api-go-sdk/model"
	"github.com/LyridInc/cluster-api-go-sdk/option"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// export $(< test.env)

func TestInfrastructureReadiness(t *testing.T) {
	capi := api.NewClusterApiClient("", "../local.kubeconfig")
	_, err := capi.InfrastructureReadiness("openstack")
	if err != nil {
		t.Fatalf(error.Error(err))
	}
}

func TestUnavailableInfrastructureReadiness(t *testing.T) {
	capi := api.NewClusterApiClient("", "../local.kubeconfig")
	_, err := capi.InfrastructureReadiness("oci")
	if err != nil {
		t.Fatalf(error.Error(err))
	}
}

// go test ./test -v -run ^TestGetWorkloadClusterKubeconfig$
func TestGetWorkloadClusterKubeconfig(t *testing.T) {
	capi := api.NewClusterApiClient("", "./data/local.kubeconfig")
	t.Run("cluster exists", func(t *testing.T) {
		clusterName := "capi-local-3"
		conf, err := capi.GetWorkloadClusterKubeconfig(clusterName)
		if err != nil {
			t.Fatalf(error.Error(err))
		}
		if err := os.WriteFile(fmt.Sprintf("./data/%s.kubeconfig", clusterName), []byte(*conf), 0644); err != nil {
			t.Fatal("Write kubeconfig error:", err)
		}
	})
	t.Run("cluster doesn't exist", func(t *testing.T) {
		_, err := capi.GetWorkloadClusterKubeconfig("capi-local-99")
		if err != nil {
			t.Fatalf(error.Error(err))
		}
	})
}

// go test ./test -v -run ^TestKubectlManifest$
func TestKubectlManifest(t *testing.T) {
	capi := api.NewClusterApiClient("", "../local.kubeconfig")
	yamlByte, err := os.ReadFile("../capi-local.yaml") // workload cluster yaml manifest
	if err != nil {
		t.Fatal(error.Error(err))
	}
	yaml := string(yamlByte)

	t.Run("kubectl apply -f", func(t *testing.T) {
		cl := api.OpenstackClient{
			NetworkEndpoint: os.Getenv("OS_NETWORK_ENDPOINT"),
			AuthEndpoint:    os.Getenv("OS_AUTH_ENDPOINT"),
			AuthToken:       os.Getenv("OS_TOKEN"),
			ProjectId:       os.Getenv("OS_PROJECT_ID"),
		}

		if quotaAvailable, err := cl.ValidateQuotas(); !quotaAvailable || err != nil {
			t.Fatal("Quota problem:", error.Error(err))
		}

		if err := capi.ApplyYaml(yaml); err != nil {
			t.Fatal(error.Error(err))
		}
	})

	// go test ./test -run ^\QTestKubectlManifest\E$/^\Qkubectl_delete_-f\E$
	t.Run("kubectl delete -f", func(t *testing.T) {
		if err := capi.DeleteYaml(yaml); err != nil {
			t.Fatal(error.Error(err))
		}
	})
}

// go test ./test -v -run ^TestCreateSecret$
func TestCreateSecret(t *testing.T) {
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

	cloudConf := os.Getenv("OPENSTACK_CLOUD_PROVIDER_CONF")
	if cloudConf == "" {
		t.Fatal("Error reading cloud conf: OPENSTACK_CLOUD_PROVIDER_CONF is not set")
	}

	capi := api.NewClusterApiClient("", "./data/local.kubeconfig")
	if err := capi.SetKubernetesClientset("./data/capi-local-3.kubeconfig"); err != nil {
		t.Fatal("Error set kubeconfig:", error.Error(err))
	}

	// 	NAMESPACE     NAME           TYPE     DATA   AGE
	//  kube-system   cloud-config   Opaque   1      19s

	// {"cloud-2.conf": "base64string"}
	// kubectl get secret -n kube-system --kubeconfig=/mnt/c/Users/Lyrid/Documents/Projects/cluster-api-sdk/test/data/capi-local-2.kubeconfig cloud-config -o jsonpath="{.data.cloud\.conf}" | base64 --decode
	// kubectl get pods -A --kubeconfig=/mnt/c/Users/Lyrid/Documents/Projects/cluster-api-sdk/test/data/capi-local-2.kubeconfig
	// kubectl --kubeconfig=./${CLUSTER_NAME}.kubeconfig create secret -n kube-system generic cloud-config --from-file=/tmp/cloud.conf
	secret := v1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: v1.SchemeGroupVersion.String(),
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cloud-config",
			Namespace: "kube-system",
		},
		Data: map[string][]byte{
			"cloud.conf": []byte(cloudConf),
		},
	}
	secretValue, err := capi.CreateSecret(secret)
	if err != nil {
		t.Fatal("Error create secret:", error.Error(err))
	}
	t.Log((*secretValue).Data)
}
