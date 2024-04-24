package test

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/LyridInc/cluster-api-go-sdk/api"
	"github.com/LyridInc/cluster-api-go-sdk/model"
	"github.com/LyridInc/cluster-api-go-sdk/option"
	"github.com/LyridInc/cluster-api-go-sdk/utils"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// export $(< test.env)

// go test ./test -v -run ^TestGenerateClusterTemplate$
func TestGenerateClusterTemplate(t *testing.T) {
	yamlByte, _ := os.ReadFile("./data/elitery-clouds.yaml")
	cloudsYaml := model.CloudsYaml{}
	cloudsYaml.Parse(yamlByte)
	opt := option.OpenstackGenerateClusterOptions{
		ControlPlaneMachineFlavor: "a2.medium-1",
		NodeMachineFlavor:         "a2.large-2",
		ExternalNetworkId:         "f30c9e3d-757b-43fb-b4e0-da3ab36708a4",
		ImageName:                 "Ubuntu-22.04-eranyaImage-v1.0",
		SshKeyName:                "eranya-ssh",
		DnsNameServers:            "8.8.8.8",
		FailureDomain:             "az-01", // nova/az-01
		IgnoreVolumeAZ:            true,
	}
	cloudsYaml.SetEnvironment(opt)

	infrastructure := "openstack"
	capi, _ := api.NewClusterApiClient("", "./data/az-vega.kubeconfig")

	clusterName := "capi-elitery"
	ready, err := capi.InfrastructureReadiness(infrastructure)
	if !ready && err == nil {
		t.Log("initialize infrastructure")
		capi.InitInfrastructure(infrastructure)
	}

	t.Log("Generate workload cluster YAML")
	clusterOpt := option.GenerateWorkloadClusterOptions{
		ClusterName:              clusterName,
		KubernetesVersion:        "v1.24.8",
		WorkerMachineCount:       1,
		ControlPlaneMachineCount: 1,
		InfrastructureProvider:   infrastructure,
		Flavor:                   "",
		URL:                      "./data/template-external-cloud-provider.yaml",
	}
	yaml, err := capi.GenerateWorkloadClusterYaml(clusterOpt)
	if err != nil {
		t.Fatal("Generate workload cluster error:", err)
	}

	if err := os.WriteFile(fmt.Sprintf("./data/%s.yaml", clusterName), []byte(yaml), 0644); err != nil {
		t.Fatal("Write yaml error:", err)
	}
}

// go test ./test -v -run ^TestGetWorkloadClusterKubeconfig$
func TestGetWorkloadClusterKubeconfig(t *testing.T) {
	capi, _ := api.NewClusterApiClient("", "./data/local.kubeconfig")
	t.Run("cluster exists", func(t *testing.T) {
		clusterName := "capi-testing"
		namespace := "default"
		conf, err := capi.GetWorkloadClusterKubeconfig(clusterName, namespace)
		if err != nil {
			t.Fatalf(error.Error(err))
		}
		if err := os.WriteFile(fmt.Sprintf("./data/%s.kubeconfig", clusterName), []byte(*conf), 0644); err != nil {
			t.Fatal("Write kubeconfig error:", err)
		}
	})
	t.Run("cluster doesn't exist", func(t *testing.T) {
		_, err := capi.GetWorkloadClusterKubeconfig("capi-local", "default")
		if err != nil {
			t.Fatalf(error.Error(err))
		}
	})
}

// go test ./test -v -run ^TestSetClientsetFromConfigBytes$
func TestSetClientsetFromConfigBytes(t *testing.T) {
	capi, _ := api.NewClusterApiClient("", "./data/local.kubeconfig")
	clusterName := "capi-local-2"
	namespace := "default"
	conf, err := capi.GetWorkloadClusterKubeconfig(clusterName, namespace)
	if err != nil {
		t.Fatalf(error.Error(err))
	}
	if err := os.WriteFile(fmt.Sprintf("./data/%s.kubeconfig", clusterName), []byte(*conf), 0644); err != nil {
		t.Fatal("Write kubeconfig error:", err)
	}

	capi.SetKubernetesClientsetFromConfigBytes([]byte(*conf))

	secret := v1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: v1.SchemeGroupVersion.String(),
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: "kube-system",
		},
		Data: map[string][]byte{
			"key": []byte("test secret"),
		},
	}

	if _, err := capi.CreateSecret(secret); err != nil {
		t.Fatal("Error create secret:", error.Error(err))
	}
}

// go test ./test -v -run ^TestKubectlManifest$
func TestKubectlManifest(t *testing.T) {
	capi, _ := api.NewClusterApiClient("", "../local.kubeconfig")
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
		FailureDomain:             "az-01", // nova/az-01
	}
	cloudsYaml.SetEnvironment(opt)

	cloudConf := os.Getenv("OPENSTACK_CLOUD_PROVIDER_CONF")
	if cloudConf == "" {
		t.Fatal("Error reading cloud conf: OPENSTACK_CLOUD_PROVIDER_CONF is not set")
	}

	t.Log(cloudConf)

	capi, _ := api.NewClusterApiClient("", "./data/local.kubeconfig")
	if err := capi.SetKubernetesClientset("./data/capi-local-2.kubeconfig"); err != nil {
		t.Fatal("Error set kubeconfig:", error.Error(err))
	}

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
	x := (*secretValue).Data
	t.Log(string(x["cloud.conf"]))
}

// go test ./test -v -run ^TestInitializeInfrastructure$
func TestInitializeInfrastructure(t *testing.T) {
	capi, err := api.NewClusterApiClient("", "C:/Users/Lyrid/.kube/beta.config")
	if err != nil {
		t.Fatal(err)
	}

	ready, err := capi.InfrastructureReadiness("openstack")
	if !ready || err != nil {
		clComponents, err := capi.InitInfrastructure("openstack")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(clComponents)
	}
}

// go test ./test -v -run ^TestCreateNamespace$
func TestCreateNamespace(t *testing.T) {
	capi, _ := api.NewClusterApiClient("", "./data/capi-helm-testing.kubeconfig")
	namespace := "test-ns"

	t.Run("create namespace", func(t *testing.T) {
		ns, err := capi.CreateNamespace(namespace)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(ns)
	})

	t.Run("add label namespace", func(t *testing.T) {
		ns, err := capi.AddLabelNamespace(namespace, "istio-injection", "enabled")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(ns)
	})
}

// go test ./test -v -run ^TestKubectlManifestWithLabelSelector$
func TestKubectlManifestWithLabelSelector(t *testing.T) {
	capi, _ := api.NewClusterApiClient("", "./data/capi-helm-testing.kubeconfig")

	// kubectl apply -l knative.dev/crd-install=true -f https://github.com/knative/net-istio/releases/download/knative-v1.8.1/istio.yaml
	t.Run("kubectl apply -l knative.dev/crd-install=true -f", func(t *testing.T) {
		capi.LabelSelector = &metav1.LabelSelector{MatchLabels: map[string]string{"knative.dev/crd-install": "true"}}
		yaml, err := model.ReadYamlFromUrl("https://github.com/knative/net-istio/releases/download/knative-v1.8.1/istio.yaml")
		if err != nil {
			t.Fatal(error.Error(err))
		}

		if err := capi.ApplyYaml(yaml); err != nil {
			t.Fatal(error.Error(err))
		}
	})
}

// go test ./test -v -run ^TestGetKubeconfigValues$
func TestGetKubeconfigValues(t *testing.T) {
	t.Run("has ca data", func(t *testing.T) {
		capi, _ := api.NewClusterApiClient("", "./data/beta.config")
		b, _ := os.ReadFile("./data/experiment.kubeconfig")
		values, _ := capi.GetConfigValues(b)
		t.Log(values["cert_data"])
		t.Log(values["certificate_authority_data"])
		t.Log(values["server"])
		t.Log(values["bearer_token"])
	})
}

// go test ./test -v -run ^TestGetService$
func TestGetService(t *testing.T) {
	t.Run("check service", func(t *testing.T) {
		capi, _ := api.NewClusterApiClient("", "./data/certificatetest-yahv.kubeconfig")
		s, err := capi.GetService("vega", "lyrid-9cc8b789-e6df-434a-afbb-371e8280ec1a")
		if err != nil {
			t.Fatal(err)
		}
		b, _ := json.Marshal(s)
		t.Log(string(b))

		s, err = capi.GetService("lyrid-f2d95330-796c-4735-b4dd-5ae021779141", "lyrid-faa37fc4-f5b8-441c-89c9-1f56fd41356f")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(s.Status.Conditions)
	})

	t.Run("check LB", func(t *testing.T) {
		capi, _ := api.NewClusterApiClient("", "./data/capi-az-local.kubeconfig")
		s, err := capi.GetService("ingress-dftw8qey-ingress-nginx-controller", "lyrid-9cc8b789-e6df-434a-afbb-371e8280ec1a")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(s.Spec.LoadBalancerIP)
		t.Log(s.Spec.ExternalIPs)
		t.Log(s.Spec.ClusterIP)
		t.Log(s.Status.LoadBalancer.Ingress)
		t.Log(s.Annotations["loadbalancer.openstack.org/load-balancer-id"])
	})

}

// go test ./test -v -run ^TestGetSecret$
func TestGetSecret(t *testing.T) {
	capi, _ := api.NewClusterApiClient("", "./data/experiment.kubeconfig")
	s, err := capi.GetSecret("lyrid-admin", "kube-system")
	if err != nil {
		t.Fatal(err)
	}
	jsonData, err := json.Marshal(s.Data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jsonData))
}

// go test ./test -v -run ^TestCreateRegistrySecret$
func TestCreateRegistrySecret(t *testing.T) {
	capi, _ := api.NewClusterApiClient("", "./data/local.kubeconfig")

	secretName := "lyridlocal.key"
	secretNs := "lyrid-9cc8b789-e6df-434a-afbb-371e8280ec1a"
	dockerUsername := "<docker-username>"
	dockerPassword := "<docker-password>"
	dockerServer := "<docker-server>"

	secretValue, err := capi.CreateDockerRegistrySecret(secretName, secretNs, model.CreateDockerRegistrySecretArgs{
		Username: dockerUsername,
		Password: dockerPassword,
		Email:    "admin@mail.io",
		Server:   dockerServer,
	})
	if err != nil {
		t.Fatal("Error create secret:", error.Error(err))
	}
	x := (*secretValue).Data
	t.Log(x)

}

// go test ./test -v -run ^TestPatchServiceAccount$
func TestPatchServiceAccount(t *testing.T) {
	capi, _ := api.NewClusterApiClient("", "./data/local.kubeconfig")

	patch := []byte("{\"imagePullSecrets\": [{\"name\": \"lyridlocaltest.key\"}]}")
	namespace := "lyrid-9cc8b789-e6df-434a-afbb-371e8280ec1a"

	sa, err := capi.Clientset.CoreV1().ServiceAccounts(namespace).Patch(context.Background(), "default", types.StrategicMergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		t.Fatal("Error patch service accounts:", error.Error(err))
	}

	t.Log(sa)
}

// go test ./test -v -run ^TestPatchConfigMap$
func TestPatchConfigMap(t *testing.T) {
	capi, _ := api.NewClusterApiClient("", "./data/zzz-lyrid-local.kubeconfig")
	manifestUrl := "https://storage.beta.lyrid.io/client/vega-configs-template/config-domain-updated.yaml"
	configYaml, _ := model.ReadYamlFromUrl(manifestUrl)

	y := map[string]interface{}{}
	yaml.Unmarshal([]byte(configYaml), &y)

	jsonInterface, err := utils.ConvertYAMLToJSON(y)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(jsonInterface)

	jsonByte, _ := json.Marshal(jsonInterface)

	t.Log(string(jsonByte))
	_, err = capi.PatchConfigMap("config-domain", "knative-serving", jsonByte)
	if err != nil {
		t.Fatal(err)
	}
}

// go test ./test -v -run ^TestAssignSecretCert$
func TestAssignSecretCert(t *testing.T) {
	privatekey, _ := os.ReadFile("./data/9cc8b789e6df434aafbb371e8280ec1a__.cesc.lyr.id_2023-05-23 08_43_25 +0000 UTC_cesc.lyr.id.key")
	certificate, _ := os.ReadFile("./data/9cc8b789e6df434aafbb371e8280ec1a__.cesc.lyr.id_2023-05-23 08_43_25 +0000 UTC_cesc.lyr.id.crt")
	capi, _ := api.NewClusterApiClient("", "./data/lyrid-client-demo.kubeconfig")
	req := &v1.Secret{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: "vega.cesc.lyr.id",
		},
		Immutable: nil,
		Data: map[string][]byte{
			"tls.key": privatekey,
			"tls.crt": certificate,
		},
		StringData: nil,
		Type:       "kubernetes.io/tls",
	}

	res, err := capi.Clientset.CoreV1().Secrets("lyrid-9cc8b789-e6df-434a-afbb-371e8280ec1a").Create(context.Background(), req, metav1.CreateOptions{})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(res)
}

// go test ./test -v -run ^TestDeleteCluster$
func TestDeleteCluster(t *testing.T) {
	capi, _ := api.NewClusterApiClient("", "./data/az-vega.kubeconfig")

	discoveryClient := capi.Clientset.Discovery()
	apiResourceList, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		t.Fatalf(err.Error())
	}

	for _, apiResource := range apiResourceList {
		if apiResource.GroupVersion == "cluster.x-k8s.io/v1beta1" {
			for _, resource := range apiResource.APIResources {
				if resource.Name == "clusters" {
					result := capi.Clientset.RESTClient().Delete().
						AbsPath("apis/cluster.x-k8s.io/v1beta1/namespaces/default/clusters/yyy-lyrid-dev").
						VersionedParams(&metav1.GetOptions{}, metav1.ParameterCodec).
						Do(context.TODO())
					b, _ := result.Raw()
					t.Log(string(b))
				}
			}
		}
	}
}

// go test ./test -v -run ^TestGetRevision$
func TestGetRevision(t *testing.T) {
	capi, _ := api.NewClusterApiClient("", "./data/biznet.config")
	b, err := capi.GetKNativeRevision("lyrid-cf8bf2c4-650a-4550-be9b-ca7557a13adb-00001", "lyrid-6a4b7713-b06d-4c31-9d96-ce81e96bd72c")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(b))
}

// go test ./test -v -run ^TestGetConfiguration$
func TestGetConfiguration(t *testing.T) {
	capi, _ := api.NewClusterApiClient("", "./data/biznet.config")
	b, err := capi.GetKNativeConfiguration("lyrid-ee249691-5b2e-44c8-a032-f3dc1998caad", "lyrid-3250c3b9-3f96-417e-82f6-ac0678002f61")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(b))
}

// go test ./test -v -run ^TestGetClusterK8sResource$
func TestGetClusterK8sResource(t *testing.T) {
	capi, _ := api.NewClusterApiClient("", "./data/az-vega.kubeconfig")
	b, err := capi.GetClusterK8sResource("lyrid-patch-yaml-midn", "lyrid-9cc8b789-e6df-434a-afbb-371e8280ec1a")
	if err != nil {
		t.Fatal(err)
	}

	var cluster map[string]interface{}
	json.Unmarshal(b, &cluster)

	t.Log(cluster["status"])

	var clusterStatus map[string]interface{} = cluster["status"].(map[string]interface{})
	t.Log(clusterStatus["phase"])
}

// go test ./test -v -run ^TestUpdateClusterK8sResource$
func TestUpdateClusterK8sResource(t *testing.T) {
	capi, _ := api.NewClusterApiClient("", "./data/az-vega.kubeconfig")
	un, err := capi.UpdateClusterK8sResourceAnnotations("lyrid-patch-yaml-midn", "lyrid-9cc8b789-e6df-434a-afbb-371e8280ec1a", "abc")
	if err != nil {
		t.Fatal(err)
	}

	b, _ := json.Marshal(un.Object)

	t.Log(string(b))
}

// go test ./test -v -run ^TestDialContext$
func TestDialContext(t *testing.T) {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.DialContext(ctx, network, fmt.Sprintf("%s:%d", "205.234.21.72", 443))
			},
		},
	}

	req, _ := http.NewRequest("GET", "https://api.lyrid.io/api/vega/list", nil)
	req.Header.Add("Authorization", "Bearer <token>")
	req.Header.Add("Content-Type", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(res)
}

// go test ./test -v -run ^TestGetDeployment$
func TestGetDeployment(t *testing.T) {
	capi, _ := api.NewClusterApiClient("", "./data/certificatetest-yahv.kubeconfig")
	deployment, err := capi.GetDeployment("vega", "lyrid-9cc8b789-e6df-434a-afbb-371e8280ec1a")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(deployment.Spec.Template.Spec.Containers[0].Image)
}

// go test ./test -v -run ^TestNodeShell$
func TestNodeShell(t *testing.T) {
	t.Run("execute command in machine node", func(t *testing.T) {
		var nodeName string = "az-vega-md-0-xgzfb"
		capi, _ := api.NewClusterApiClient("", "./data/az-vega.kubeconfig")
		err := capi.ExecuteNodeShellCommand(nodeName, "apt-get install open-iscsi -y")
		if err != nil {
			t.Fatal(err)
		}
	})
}
