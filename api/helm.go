package api

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	helm "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
)

type (
	HelmClient struct {
		KubeconfigFile        string
		HelmOptions           helm.Options
		KubeConfClientOptions helm.KubeConfClientOptions
		Client                helm.Client
		Timeout               time.Duration
		ActionConfig          *action.Configuration
		EnvSettings           *cli.EnvSettings
	}
)

func NewHelmClient(kubeconfigFile, namespace string) (*HelmClient, error) {
	helmOption := helm.Options{
		Debug:     true,
		Linting:   true,
		Namespace: namespace,
	}
	kubeconfig, err := os.ReadFile(kubeconfigFile)
	if err != nil {
		return nil, err
	}
	kubeConfClientOption := helm.KubeConfClientOptions{
		Options:     &helmOption,
		KubeContext: "",
		KubeConfig:  kubeconfig,
	}

	helmClient, err := helm.NewClientFromKubeConf(&kubeConfClientOption, helm.Burst(100), helm.Timeout(10e9))
	if err != nil {
		return nil, err
	}

	settings := cli.New()
	settings.KubeConfig = kubeconfigFile
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {}); err != nil {
		return nil, err
	}

	return &HelmClient{
		HelmOptions:           helmOption,
		KubeConfClientOptions: kubeConfClientOption,
		KubeconfigFile:        kubeconfigFile,
		Client:                helmClient,
		Timeout:               60 * time.Second,
		ActionConfig:          actionConfig,
		EnvSettings:           settings,
	}, nil
}

func (c *HelmClient) AddRepo(entry repo.Entry) error {
	return c.Client.AddOrUpdateChartRepo(entry)
}

func (c *HelmClient) Install(chartName, releaseName, namespace string, wait bool) (*release.Release, error) {
	spec := helm.ChartSpec{
		ReleaseName:     releaseName,
		ChartName:       chartName,
		Namespace:       namespace,
		CreateNamespace: true,
		Wait:            wait,
		Timeout:         c.Timeout,
	}
	opt := helm.GenericHelmOptions{}
	return c.Client.InstallOrUpgradeChart(context.Background(), &spec, &opt)
}

func (c *HelmClient) CliInstall(chartName, releaseName, namespace string, settingValues []string) (*release.Release, error) {
	installAction := action.NewInstall(c.ActionConfig)
	installAction.Namespace = namespace
	chartPath, err := installAction.LocateChart(chartName, c.EnvSettings)
	if err != nil {
		return nil, err
	}

	values := make(map[string]interface{})
	for _, s := range settingValues {
		values, err = setHelmValue(values, s)
		if err != nil {
			return nil, err
		}
	}

	installAction.ReleaseName = releaseName
	chart, err := loader.Load(chartPath)
	if err != nil {
		return nil, err
	}

	return installAction.Run(chart, values)
}

// setHelmValue is a helper function to convert a Helm-style value string to a map.
func setHelmValue(vals map[string]interface{}, v string) (map[string]interface{}, error) {
	split := strings.Split(v, "=")
	if len(split) != 2 {
		return vals, fmt.Errorf("incorrect format for helm value '%s'", v)
	}
	var val interface{}
	if strings.HasSuffix(split[0], ".enabled") {
		vv, err := strconv.ParseBool(split[1])
		if err != nil {
			return vals, fmt.Errorf("incorrect data type for '%s'", v)
		}
		val = vv
	} else {
		val = split[1]
	}

	parts := strings.Split(split[0], ".")
	m := vals
	for i, p := range parts {
		if i == len(parts)-1 {
			m[p] = val
			break
		}
		_, ok := m[p]
		if !ok {
			m[p] = make(map[string]interface{})
		}
		m = m[p].(map[string]interface{})
	}
	return vals, nil
}

// helm upgrade --install "${INGRESS}" -n "${LYRID_NAMESPACE}" ingress-nginx/ingress-nginx -f ${INSTALLER_PATH}/ingress-values.yaml
// helm upgrade --install --wait-for-jobs --timeout 5m0s vega -n "${LYRID_NAMESPACE}" ${CHART_PATH}/
// helm status
// http request function to be hit to Lyra API
// replace values function in values.yaml
// grab value from kubeconfig and config
