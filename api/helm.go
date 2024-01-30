package api

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	helm "github.com/mittwald/go-helm-client"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/registry"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
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
	settings.SetNamespace(namespace)
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {}); err != nil {
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

func NewHelmClientFromConfigBytes(configBytes []byte, namespace string) (*HelmClient, error) {
	helmOption := helm.Options{
		Debug:     true,
		Linting:   true,
		Namespace: namespace,
	}
	kubeConfClientOption := helm.KubeConfClientOptions{
		Options:     &helmOption,
		KubeContext: "",
		KubeConfig:  configBytes,
	}

	helmClient, err := helm.NewClientFromKubeConf(&kubeConfClientOption, helm.Burst(100), helm.Timeout(10e9))
	if err != nil {
		return nil, err
	}

	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(&SimpleRESTClientGetter{
		Namespace:  namespace,
		KubeConfig: string(configBytes),
	}, namespace, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {}); err != nil {
		return nil, err
	}

	return &HelmClient{
		HelmOptions:           helmOption,
		KubeConfClientOptions: kubeConfClientOption,
		Client:                helmClient,
		Timeout:               60 * time.Second,
		ActionConfig:          actionConfig,
		EnvSettings:           settings,
	}, nil
}

func (c *HelmClient) AddRepo(entry repo.Entry) error {
	return c.Client.AddOrUpdateChartRepo(entry)
}

func (c *HelmClient) Pull(chartRef, registryUrl, username, password string) error {
	client, err := registry.NewClient()
	if err != nil {
		return err
	}
	if err := client.Login(registryUrl, registry.LoginOptBasicAuth(username, password)); err != nil {
		return err
	}

	result, err := client.Pull(registryUrl+"/"+chartRef, registry.PullOptWithChart(true))
	if err != nil {
		return err
	}

	log.Println(result)
	log.Println(string(result.Chart.Data))

	// var buf bytes.Buffer
	// gzWriter := gzip.NewWriter(buf)

	return nil
}

func (c *HelmClient) Install(chartName, releaseName, version, namespace string, values map[string]interface{}, wait bool) (*release.Release, error) {
	var valuesYaml string
	if values != nil {
		yamlData, err := yaml.Marshal(values)
		if err != nil {
			return nil, err
		}
		valuesYaml = string(yamlData)
	}

	spec := helm.ChartSpec{
		ReleaseName:     releaseName,
		ChartName:       chartName,
		Namespace:       namespace,
		CreateNamespace: true,
		Wait:            wait,
		Timeout:         c.Timeout,
		Version:         version,
		ValuesYaml:      valuesYaml,
	}

	opt := helm.GenericHelmOptions{}
	var finalError error

	retry := 0
	for {
		if retry >= 3 {
			break
		}

		release, err := c.Client.InstallOrUpgradeChart(context.Background(), &spec, &opt)
		finalError = err
		if err != nil {
			fmt.Println("Retry installing ", releaseName+"...")
			retry = retry + 1
		} else {
			return release, nil
		}
	}

	return nil, finalError
}

func (c *HelmClient) CliAddRepo(entry repo.Entry) error {
	r, err := repo.NewChartRepository(&entry, getter.All(c.EnvSettings))
	if err != nil {
		return err
	}

	if _, err := r.DownloadIndexFile(); err != nil {
		return err
	}

	var f repo.File
	f.Update(&entry)

	if err := f.WriteFile(c.EnvSettings.RepositoryConfig, 0644); err != nil {
		return err
	}

	return nil
}

func (c *HelmClient) CliInstall(chartName, releaseName, namespace, version string, settingValues []string) (*release.Release, error) {
	installAction := action.NewInstall(c.ActionConfig)
	installAction.Namespace = namespace
	chartPath, err := installAction.LocateChart(chartName, c.EnvSettings)
	if err != nil {
		return nil, err
	}

	if version != "" {
		installAction.Version = version
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

func (c *HelmClient) CliUpgrade(chartPath, releaseName, namespace string, values map[string]interface{}, timeout time.Duration, reuseValues, waitForJobs bool) (*release.Release, error) {
	upgradeAction := action.NewUpgrade(c.ActionConfig)
	upgradeAction.WaitForJobs = waitForJobs
	upgradeAction.Timeout = timeout
	upgradeAction.Namespace = namespace
	upgradeAction.Install = !reuseValues
	upgradeAction.ReuseValues = reuseValues

	chart, err := loader.Load(chartPath)
	if err != nil {
		chartPath, err = upgradeAction.LocateChart(chartPath, c.EnvSettings)
		if err != nil {
			return nil, err
		}
		chart, err = loader.Load(chartPath)
		if err != nil {
			return nil, err
		}
	}

	release, err := upgradeAction.Run(releaseName, chart, values)
	if err != nil {
		errMsg := error.Error(err)
		if strings.HasSuffix(errMsg, " has no deployed releases") {
			installAction := action.NewInstall(c.ActionConfig)
			installAction.WaitForJobs = waitForJobs
			installAction.Timeout = timeout
			installAction.Namespace = namespace
			installAction.ReleaseName = releaseName
			return installAction.Run(chart, values)
		} else {
			return nil, err
		}
	}

	return release, nil
}

func (c *HelmClient) CliDelete(releaseName string) (*release.UninstallReleaseResponse, error) {
	deleteAction := action.NewUninstall(c.ActionConfig)

	response, err := deleteAction.Run(releaseName)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *HelmClient) CliStatus(releaseName string) (map[string]interface{}, error) {
	actionStatus := action.NewStatus(c.ActionConfig)
	release, err := actionStatus.Run(releaseName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"name":        release.Name,
		"status":      release.Info.Status.String(),
		"description": release.Info.Description,
		"version":     release.Version,
	}, nil
}

func (c *HelmClient) ReplaceYamlPlaceholder(yaml, placeholder, value string) string {
	return strings.ReplaceAll(yaml, placeholder, value)
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

type SimpleRESTClientGetter struct {
	Namespace  string
	KubeConfig string
}

func NewRESTClientGetter(namespace, kubeConfig string) *SimpleRESTClientGetter {
	return &SimpleRESTClientGetter{
		Namespace:  namespace,
		KubeConfig: kubeConfig,
	}
}

func (c *SimpleRESTClientGetter) ToRESTConfig() (*rest.Config, error) {
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(c.KubeConfig))
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *SimpleRESTClientGetter) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	config, err := c.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	// The more groups you have, the more discovery requests you need to make.
	// given 25 groups (our groups + a few custom conf) with one-ish version each, discovery needs to make 50 requests
	// double it just so we don't end up here again for a while.  This config is only used for discovery.
	config.Burst = 100

	discoveryClient, _ := discovery.NewDiscoveryClientForConfig(config)
	return memory.NewMemCacheClient(discoveryClient), nil
}

func (c *SimpleRESTClientGetter) ToRESTMapper() (meta.RESTMapper, error) {
	discoveryClient, err := c.ToDiscoveryClient()
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)
	expander := restmapper.NewShortcutExpander(mapper, discoveryClient)
	return expander, nil
}

func (c *SimpleRESTClientGetter) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	// use the standard defaults for this client command
	// DEPRECATED: remove and replace with something more accurate
	loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig

	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}
	overrides.Context.Namespace = c.Namespace

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides)
}
