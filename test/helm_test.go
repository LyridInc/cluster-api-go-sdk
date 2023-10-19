package test

import (
	"os"
	"testing"
	"time"

	"github.com/LyridInc/cluster-api-go-sdk/api"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
)

// export $(< test.env)

// go test ./test -v -run ^TestHelmGetRelease$
func TestHelmGetRelease(t *testing.T) {
	namespace := "default"
	hc, err := api.NewHelmClient("./data/experiment.kubeconfig", namespace)
	if err != nil {
		t.Fatal(error.Error(err))
	}

	t.Run("get list of releases", func(t *testing.T) {
		releases, err := hc.Client.ListDeployedReleases()
		if err != nil {
			t.Fatal(error.Error(err))
		}
		for _, r := range releases {
			t.Log(r.Name)
		}
	})

	t.Run("get release by name", func(t *testing.T) {
		release, err := hc.Client.GetRelease("redis")
		if err != nil {
			t.Fatal(error.Error(err))
		}
		t.Log(release.Name)
	})
}

// go test ./test -v -run ^TestHelmGetReleaseValues$
func TestHelmGetReleaseValues(t *testing.T) {
	namespace := "lyrid-9cc8b789-e6df-434a-afbb-371e8280ec1a"
	hc, err := api.NewHelmClient("./data/capi-testing.kubeconfig", namespace)
	if err != nil {
		t.Fatal(error.Error(err))
	}

	values, err := hc.Client.GetReleaseValues("vega", true)
	if err != nil {
		t.Fatal(error.Error(err))
	}
	t.Log(values)
}

// go test ./test -v -run ^TestHelmInstallChart$
func TestHelmInstallChart(t *testing.T) {
	namespace := "default"
	version := "v1.16.3"
	kubeconfig := "./data/lyr-debug.kubeconfig"
	hc, err := api.NewHelmClient(kubeconfig, namespace)
	if err != nil {
		t.Fatal(error.Error(err))
	}

	if err := hc.AddRepo(repo.Entry{
		Name: "istio",
		URL:  "https://istio-release.storage.googleapis.com/charts",
	}); err != nil {
		t.Fatal(error.Error(err))
	}

	t.Run("install istio-base", func(t *testing.T) {
		release, err := hc.Install("istio/base", "istio-base", version, "istio-system", nil, false)
		if err != nil {
			t.Fatal(error.Error(err))
		}
		t.Log(release)
	})

	t.Run("install istiod", func(t *testing.T) {
		release, err := hc.Install("istio/istiod", "istiod", version, "istio-system", nil, true)
		if err != nil {
			t.Fatal(error.Error(err))
		}
		t.Log(release)
	})

	t.Run("install ingress gateway", func(t *testing.T) {
		capi, _ := api.NewClusterApiClient("", kubeconfig)

		if _, err := capi.AddLabelNamespace("istio-system", "istio-injection", "enabled"); err != nil {
			t.Fatal(err)
		}

		values := map[string]interface{}{
			"service": map[string]interface{}{
				"type": "ClusterIP",
			},
		}

		release, err := hc.Install("istio/gateway", "istio-ingressgateway", version, "istio-system", values, true)
		if err != nil {
			t.Fatal(error.Error(err))
		}
		t.Log(release)
	})
}

// go test ./test -v -run ^TestHelmInstallChartWithSet$
func TestHelmInstallChartWithSet(t *testing.T) {
	namespace := "default"
	kubeconfig := "./data/capi-helm-testing.kubeconfig"
	hc, err := api.NewHelmClient(kubeconfig, namespace)
	if err != nil {
		t.Fatal(error.Error(err))
	}

	release, err := hc.CliInstall("bitnami/redis", "redis", namespace, []string{
		"architecture=standalone",
		"auth.enabled=false",
		"master.persistence.enabled=false",
	})
	if err != nil {
		t.Fatal(error.Error(err))
	}
	t.Log(release)
}

// go test ./test -v -run ^TestHelmUpgradeChart$
func TestHelmUpgradeChart(t *testing.T) {
	namespace := "default"
	kubeconfig := "./data/capi-helm-testing.kubeconfig"
	hc, err := api.NewHelmClient(kubeconfig, namespace)
	if err != nil {
		t.Fatal(error.Error(err))
	}

	timeout := time.Second * (5 * 60)

	release, err := hc.CliUpgrade("./data/chart", "vega", namespace, nil, timeout, false, true)
	if err != nil {
		t.Fatal(error.Error(err))
	}
	t.Log(release)
}

// go test ./test -v -run ^TestHelmReleaseStatus$
func TestHelmReleaseStatus(t *testing.T) {
	namespace := "lyrid-9cc8b789-e6df-434a-afbb-371e8280ec1a"
	kubeconfig := "./data/capi-delta-vega.kubeconfig"
	hc, err := api.NewHelmClient(kubeconfig, namespace)
	if err != nil {
		t.Fatal(error.Error(err))
	}

	release, err := hc.CliStatus("redis")
	if err != nil {
		t.Fatal(error.Error(err))
	}
	t.Log(release["status"] == "deployed")
}

// go test ./test -v -run ^TestHelmReplaceChartValues$
func TestHelmReplaceChartValues(t *testing.T) {
	namespace := "default"
	kubeconfig := "./data/capi-helm-testing.kubeconfig"
	hc, err := api.NewHelmClient(kubeconfig, namespace)
	if err != nil {
		t.Fatal(error.Error(err))
	}

	yamlByte, _ := os.ReadFile("./data/chart-beta/values.yaml")
	yaml := string(yamlByte)
	yaml = hc.ReplaceYamlPlaceholder(yaml, "VEGA_TAG", "bigbang")

	t.Log(yaml)
}

// go test ./test -v -run ^TestHelmUpgradeChartValues$
func TestHelmUpgradeChartValues(t *testing.T) {
	namespace := "lyrid-9cc8b789-e6df-434a-afbb-371e8280ec1a"
	kubeconfig := "./data/lyr-kube-qx4dhr.kubeconfig"
	hc, err := api.NewHelmClient(kubeconfig, namespace)
	if err != nil {
		t.Fatal(error.Error(err))
	}

	if err := hc.CliAddRepo(repo.Entry{
		Name: "ingress-nginx",
		URL:  "https://kubernetes.github.io/ingress-nginx/",
	}); err != nil {
		t.Fatal(error.Error(err))
	}

	yamlByte, err := os.ReadFile("./data/vega/ingress-values.yaml")
	if err != nil {
		t.Fatal(error.Error(err))
	}
	yaml := hc.ReplaceYamlPlaceholder(string(yamlByte), "{{ $INGRESS }}", "test-ingress")
	file, err := os.Create("./data/vega/ingress-values-f.yaml")
	if err != nil {
		t.Fatal(error.Error(err))
	}

	file.Write([]byte(yaml))
	file.Close()

	timeout := time.Second * (5 * 60)

	option := &values.Options{
		ValueFiles: []string{"./data/vega/ingress-values-f.yaml"},
	}
	provider := getter.All(hc.EnvSettings)
	values, err := option.MergeValues(provider)
	if err != nil {
		t.Fatal(error.Error(err))
	}
	release, err := hc.CliUpgrade("ingress-nginx/ingress-nginx", "ingress-tp36ousx", namespace, values, timeout, false, true)
	if err != nil {
		t.Fatal(error.Error(err))
	}
	t.Log(release)
}

// go test ./test -v -run ^TestHelmDeleteRelease$
func TestHelmDeleteRelease(t *testing.T) {
	namespace := "lyrid-9cc8b789-e6df-434a-afbb-371e8280ec1a"
	kubeconfig := "./data/delta-test.kubeconfig"
	hc, err := api.NewHelmClient(kubeconfig, namespace)
	if err != nil {
		t.Fatal(error.Error(err))
	}

	response, err := hc.CliDelete("x-vega")
	if err != nil {
		t.Fatal(error.Error(err))
	}
	t.Log(response)
}

// go test ./test -v -run ^TestChangeHelmDeployment$
func TestChangeHelmDeployment(t *testing.T) {

	// bigbang-v2.0.0-2023-10-16051712-dc92e31
	// master-v0.0.1-2023-10-16051307-e146dc3
	// helm -n lyrid-9cc8b789-e6df-434a-afbb-371e8280ec1a upgrade --reuse-values --set vega.image.tag=master-v0.0.1-2023-10-16051307-e146dc3 vega .

	namespace := "lyrid-9cc8b789-e6df-434a-afbb-371e8280ec1a"
	kubeconfig := "./data/certificatetest-yahv.kubeconfig"
	releaseName := "vega"
	hc, err := api.NewHelmClient(kubeconfig, namespace)
	if err != nil {
		t.Fatal(error.Error(err))
	}

	timeout := time.Second * (5 * 60)
	values := map[string]interface{}{
		"vega": map[string]interface{}{
			"image": map[string]interface{}{
				"tag": "bigbang-v2.0.0-2023-10-16051712-dc92e31",
			},
		},
	}
	_, err = hc.CliUpgrade("./data/latest", releaseName, namespace, values, timeout, true, true)
	if err != nil {
		t.Fatal(error.Error(err))
	}

	t.Log("success")

}
