package test

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/LyridInc/cluster-api-go-sdk/api"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	"k8s.io/client-go/tools/clientcmd"
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

	release, err := hc.CliInstall("bitnami/redis", "redis", namespace, "v18.8.0", []string{
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
	kubeconfig := "./data/certificatetest-yahv.kubeconfig"
	hc, err := api.NewHelmClient(kubeconfig, namespace)
	if err != nil {
		t.Fatal(error.Error(err))
	}

	release, err := hc.CliStatus("vega")
	if err != nil {
		t.Fatal(error.Error(err))
	}
	t.Log(release)
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

// go test ./test -v -run ^TestReadHelmChartFromStorage$
func TestReadHelmChartFromStorage(t *testing.T) {
	accessKeyId := os.Getenv("R2_ACCESSKEYID")
	accessKeySecret := os.Getenv("R2_ACCESSKEYSECRET")
	host := os.Getenv("R2_HOST")
	releaseName := "vega"
	clusterName := "certificatetest-yahv"

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s", host),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
	)
	if err != nil {
		t.Fatal(err)
	}

	r2Client := s3.NewFromConfig(cfg)
	baseDir := os.Getenv("R2_BUCKET_ENV") + "-vegaconfigs"

	object, err := r2Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(baseDir),
		Key:    aws.String("chart/vega/latest/chart.zip"),
	})
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := io.ReadAll(object.Body)
	r, err := zip.NewReader(bytes.NewReader(buf), object.ContentLength)
	if err != nil {
		t.Fatal(err)
	}

	chartFolder := "./data/" + fmt.Sprintf("%s-%s-chart", clusterName, releaseName)

	var valuesYaml *zip.File
	for _, f := range r.File {
		path := chartFolder + "/" + f.Name
		if strings.HasPrefix(f.Name, "values.yaml") {
			valuesYaml = f
		}

		fileReader, err := f.Open()
		if err != nil {
			t.Fatal(err)
		}
		defer fileReader.Close()

		buffer, err := io.ReadAll(fileReader)
		if err != nil {
			t.Fatal(err)
		}

		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			t.Fatal(err)
		}

		if f.FileInfo().IsDir() {
			continue
		}

		if err := os.WriteFile(path, buffer, f.Mode()); err != nil {
			t.Fatal(err)
		}
	}

	if valuesYaml == nil {
		err := fmt.Errorf("values.yaml not found")
		t.Fatal(err)
	}

	valuesYamlReader, err := valuesYaml.Open()
	if err != nil {
		t.Fatal(err)
	}

	valuesYamlContent, _ := io.ReadAll(valuesYamlReader)
	newContent := string(valuesYamlContent)

	namespace := "lyrid-9cc8b789-e6df-434a-afbb-371e8280ec1a"
	kubeconfig := "./data/certificatetest-yahv.kubeconfig"
	hc, err := api.NewHelmClient(kubeconfig, namespace)
	if err != nil {
		t.Fatal(error.Error(err))
	}

	conf, err := clientcmd.RESTConfigFromKubeConfig(hc.KubeConfClientOptions.KubeConfig)
	if err != nil {
		t.Fatal(err)
	}

	vegaValues, _ := os.ReadFile("./data/vega-values.json")
	valuesMap := map[string]interface{}{}
	json.Unmarshal(vegaValues, &valuesMap)

	hostName := strings.Split(valuesMap["vegaHostname"].(string), ".")
	subdomain := fmt.Sprintf("%s%s", hostName[1], ".beta.lyr.id")
	caData := base64.StdEncoding.EncodeToString(conf.CAData)

	var (
		knativeEndpoint string = "istio-ingressgateway.istio-system.svc.cluster.local"
		prometheusPort  string = "80"
	)

	values := map[string]string{
		"KUBE_CA_DATA":        caData,
		"KUBE_TOKEN":          valuesMap["kubeToken"].(string),
		"VEGA_TAG":            valuesMap["vegaTag"].(string),
		"SERVICE_PORT":        "8081",
		"SUBDOMAIN":           subdomain,
		"UIID":                valuesMap["uuid"].(string),
		"GRPC_HOST":           os.Getenv("GRPC_HOSTNAME"),
		"GRPC_PORT":           os.Getenv("GRPC_PORT"),
		"REDIS_ENDPOINT":      "redis-headless",
		"REDIS_PORT":          "6379",
		"REDIS_PW":            "",
		"KNATIVE_ENDPOINT":    knativeEndpoint,
		"PROMETHEUS_PORT":     prometheusPort,
		"PROMETHEUS_ENDPOINT": valuesMap["prometheus"].(string),
		"INGRESS_NAME":        valuesMap["ingress"].(string),
		"ENV_CONFIG":          os.Getenv("ENV_CONFIG"),
		"ACCOUNT_ID":          valuesMap["accountId"].(string),
	}

	for k, v := range values {
		newContent = strings.ReplaceAll(newContent, k, v)
	}

	t.Log(newContent)

	if err := os.WriteFile(chartFolder+"/values.yaml", []byte(newContent), valuesYaml.Mode()); err != nil {
		t.Fatal(err)
	}

	timeout := time.Second * (5 * 60)
	hc, _ = api.NewHelmClientFromConfigBytes(hc.KubeConfClientOptions.KubeConfig, namespace)
	if _, err := hc.CliUpgrade(chartFolder, releaseName, namespace, nil, timeout, false, true); err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(chartFolder)
}
