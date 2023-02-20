package test

import (
	"testing"

	"github.com/LyridInc/cluster-api-go-sdk/api"
	"helm.sh/helm/v3/pkg/repo"
)

// export $(< test.env)

// go test ./test -v -run ^TestHelmGetRelease$
func TestHelmGetRelease(t *testing.T) {
	namespace := "lyrid-9cc8b789-e6df-434a-afbb-371e8280ec1a"
	hc, err := api.NewHelmClient("./data/capi-testing.kubeconfig", namespace)
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
		release, err := hc.Client.GetRelease("vega")
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
	kubeconfig := "./data/capi-helm-testing.kubeconfig"
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
		release, err := hc.Install("istio/base", "istio-base", "istio-system", false)
		if err != nil {
			t.Fatal(error.Error(err))
		}
		t.Log(release)
	})

	t.Run("install istiod", func(t *testing.T) {
		release, err := hc.Install("istio/istiod", "istiod", "istio-system", true)
		if err != nil {
			t.Fatal(error.Error(err))
		}
		t.Log(release)
	})

	t.Run("install ingress gateway", func(t *testing.T) {
		capi, _ := api.NewClusterApiClient("", kubeconfig)

		if _, err := capi.CreateNamespace("istio-ingress"); err != nil {
			t.Fatal(err)
		}
		if _, err := capi.AddLabelNamespace("istio-ingress", "istio-injection", "enabled"); err != nil {
			t.Fatal(err)
		}

		release, err := hc.Install("istio/gateway", "istio-ingressgateway", "istio-ingress", true)
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
