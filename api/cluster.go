package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"math/rand"

	"github.com/LyridInc/cluster-api-go-sdk/model"
	"github.com/LyridInc/cluster-api-go-sdk/option"
	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	yamlserializer "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/tree"
)

type (
	ClusterApiClient struct {
		Client           client.Client
		Clientset        *kubernetes.Clientset
		DynamicInterface dynamic.Interface
		InitOptions      client.InitOptions
		Config           *rest.Config
		ConfigFile       string
		KubeconfigFile   string
		LabelSelector    *metav1.LabelSelector
		ConfigBytes      []byte
	}
)

func NewClusterApiClient(configFile, kubeconfigFile string) (*ClusterApiClient, error) {
	cl, err := client.New(context.Background(), configFile)
	if err != nil {
		log.Fatal("Client config error:", err)
		return nil, err
	}

	conf, err := clientcmd.BuildConfigFromFlags("", kubeconfigFile)
	if err != nil {
		log.Fatal("Build config from flags error:", err)
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		log.Fatal("Clientset config error:", err)
		return nil, err
	}

	dd, err := dynamic.NewForConfig(conf)
	if err != nil {
		log.Fatal("Dynamic interface config error:", err)
		return nil, err
	}

	return &ClusterApiClient{
		Client:           cl,
		Clientset:        clientset,
		Config:           conf,
		DynamicInterface: dd,
		ConfigFile:       configFile,
		KubeconfigFile:   kubeconfigFile,
		LabelSelector:    nil,
	}, nil
}

func (c *ClusterApiClient) SetRateLimit(burst int, qps float32) error {
	cl, err := client.New(context.Background(), c.ConfigFile)
	if err != nil {
		log.Fatal("Client config error:", err)
		return err
	}

	var conf *rest.Config
	if c.ConfigBytes == nil {
		conf, err = clientcmd.BuildConfigFromFlags("", c.KubeconfigFile)
		if err != nil {
			log.Fatal("Build config from flags error:", err)
			return err
		}
	} else {
		conf, err = clientcmd.RESTConfigFromKubeConfig(c.ConfigBytes)
		if err != nil {
			return err
		}
	}

	conf.Burst = burst
	conf.QPS = qps

	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		log.Fatal("Clientset config error:", err)
		return err
	}

	dd, err := dynamic.NewForConfig(conf)
	if err != nil {
		log.Fatal("Dynamic interface config error:", err)
		return err
	}

	c.DynamicInterface = dd
	c.Clientset = clientset
	c.Client = cl

	return nil
}

func (c *ClusterApiClient) GetConfigValues(configBytes []byte) (map[string]interface{}, error) {
	conf, err := clientcmd.RESTConfigFromKubeConfig(configBytes)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"certificate_authority_data": string(conf.CAData),
		"cert_data":                  string(conf.CertData),
		"server":                     string(conf.Host),
		"bearer_token":               conf.BearerToken,
	}, nil
}

func (c *ClusterApiClient) SetKubernetesClientsetFromConfigBytes(configBytes []byte) error {
	conf, err := clientcmd.RESTConfigFromKubeConfig(configBytes)
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		return err
	}

	dd, err := dynamic.NewForConfig(conf)
	if err != nil {
		log.Fatal("Dynamic interface config error:", err)
		return nil
	}

	c.Clientset = clientset
	c.DynamicInterface = dd
	c.ConfigBytes = configBytes
	c.Config = conf
	return nil
}

func (c *ClusterApiClient) SetRateLimitFromConfigBytes(burst int, qps float32, configBytes []byte) error {
	conf, err := clientcmd.RESTConfigFromKubeConfig(configBytes)
	if err != nil {
		return err
	}

	conf.Burst = burst
	conf.QPS = qps

	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		log.Fatal("Clientset config error:", err)
		return err
	}

	dd, err := dynamic.NewForConfig(conf)
	if err != nil {
		log.Fatal("Dynamic interface config error:", err)
		return err
	}

	c.DynamicInterface = dd
	c.Clientset = clientset
	c.ConfigBytes = configBytes
	c.Config = conf

	return nil
}

func (c *ClusterApiClient) SetKubernetesClientset(kubeconfigFile string) error {
	conf, err := clientcmd.BuildConfigFromFlags("", kubeconfigFile)
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		return err
	}

	dd, err := dynamic.NewForConfig(conf)
	if err != nil {
		log.Fatal("Dynamic interface config error:", err)
		return nil
	}

	c.Clientset = clientset
	c.DynamicInterface = dd
	return nil
}

func (c *ClusterApiClient) InitInfrastructure(infrastructure string) ([]client.Components, error) {
	c.InitOptions = client.InitOptions{
		Kubeconfig:              client.Kubeconfig{Path: c.KubeconfigFile},
		CoreProvider:            "",
		InfrastructureProviders: []string{infrastructure},
		BootstrapProviders:      nil,
		ControlPlaneProviders:   nil,
		TargetNamespace:         "",
		LogUsageInstructions:    true,
		WaitProviders:           false,
		WaitProviderTimeout:     time.Duration(5*60) * time.Second,
	}

	result, err := c.Client.Init(context.Background(), c.InitOptions)
	if err != nil {
		return nil, err
	}

	if ready, err := c.InfrastructureReadiness(infrastructure); !ready || err != nil {
		return nil, err
	}
	return result, nil
}

func (c *ClusterApiClient) DeleteInfrastructure(infrastructure string) error {
	return c.Client.Delete(context.Background(), client.DeleteOptions{
		Kubeconfig:              client.Kubeconfig{Path: c.KubeconfigFile},
		IncludeNamespace:        false,
		IncludeCRDs:             false,
		CoreProvider:            "",
		BootstrapProviders:      nil,
		InfrastructureProviders: []string{infrastructure},
		ControlPlaneProviders:   nil,
		DeleteAll:               false,
	})
}

func (c *ClusterApiClient) GenerateWorkloadClusterYaml(opt option.GenerateWorkloadClusterOptions) (string, error) {
	templateOptions := client.GetClusterTemplateOptions{
		Kubeconfig:               client.Kubeconfig{Path: c.KubeconfigFile},
		ClusterName:              opt.ClusterName,
		TargetNamespace:          opt.TargetNamespace,
		KubernetesVersion:        opt.KubernetesVersion,
		ListVariablesOnly:        false,
		WorkerMachineCount:       &opt.WorkerMachineCount,
		ControlPlaneMachineCount: &opt.ControlPlaneMachineCount,
	}

	if opt.Flavor != "" {
		templateOptions.ProviderRepositorySource = &client.ProviderRepositorySourceOptions{
			InfrastructureProvider: opt.InfrastructureProvider,
			Flavor:                 opt.Flavor,
		}
	}

	if opt.URL != "" {
		templateOptions.URLSource = &client.URLSourceOptions{
			URL: opt.URL,
		}
	}

	template, err := c.Client.GetClusterTemplate(context.Background(), templateOptions)
	if err != nil {
		return "", err
	}

	yaml, err := template.Yaml()
	if err != nil {
		return "", err
	}

	return string(yaml), nil
}

func (c *ClusterApiClient) GenerateOciWorkloadClusterYaml(opt option.GenerateOciWorkloadClusterOption) (string, error) {
	os.Setenv("OCI_COMPARTMENT_ID", opt.CompartmentID)
	os.Setenv("OCI_MANAGED_NODE_IMAGE_ID", opt.ImageID)
	os.Setenv("OCI_MANAGED_NODE_SHAPE", opt.Shape)
	os.Setenv("OCI_SSH_KEY", opt.SSHKey)
	os.Setenv("OCI_REGION", opt.Region)
	os.Setenv("OCI_WORKLOAD_REGION", opt.WorkloadRegion)
	os.Setenv("KUBERNETES_VERSION", opt.KubernetesVersion)
	os.Setenv("NAMESPACE", opt.Namespace)
	os.Setenv("NODE_MACHINE_COUNT", fmt.Sprintf("%d", opt.MachineCount))

	if opt.MachineTypeOCPU > 0 {
		os.Setenv("OCI_MANAGED_NODE_MACHINE_TYPE_OCPUS", fmt.Sprintf("%d", opt.MachineTypeOCPU))
	}
	if opt.BootVolumeSize != 0 {
		os.Setenv("OCI_MANAGED_NODE_BOOT_VOLUME_SIZE", fmt.Sprintf("%d", opt.BootVolumeSize))
	}

	var controlMachineCount int64 = 1
	templateOptions := client.GetClusterTemplateOptions{
		Kubeconfig: client.Kubeconfig{
			Path: c.KubeconfigFile,
		},
		ClusterName:              opt.ClusterName,
		TargetNamespace:          opt.Namespace,
		KubernetesVersion:        opt.KubernetesVersion,
		ListVariablesOnly:        false,
		WorkerMachineCount:       &opt.MachineCount,
		ControlPlaneMachineCount: &controlMachineCount,
	}

	if opt.URL != "" {
		templateOptions.URLSource = &client.URLSourceOptions{
			URL: opt.URL,
		}
	}

	template, err := c.Client.GetClusterTemplate(context.Background(), templateOptions)
	if err != nil {
		return "", err
	}

	yaml, err := template.Yaml()
	if err != nil {
		return "", err
	}

	return string(yaml), nil
}

func (c *ClusterApiClient) GenerateCloudStackWorkloadClusterYaml(opt option.GenerateCloudStackWorkloadClusterOption) (string, error) {
	os.Setenv("CLOUDSTACK_ZONE_NAME", opt.ZoneName)
	os.Setenv("CLOUDSTACK_NETWORK_NAME", opt.NetworkName)
	os.Setenv("CLUSTER_ENDPOINT_IP", opt.ClusterEndpointIP)
	os.Setenv("CLUSTER_ENDPOINT_PORT", opt.ClusterEndpointPort)
	os.Setenv("CLOUDSTACK_CONTROL_PLANE_MACHINE_OFFERING", opt.ControlPlaneMachineOffering)
	os.Setenv("CLOUDSTACK_WORKER_MACHINE_OFFERING", opt.WorkerMachineOffering)
	os.Setenv("CLOUDSTACK_TEMPLATE_NAME", opt.TemplateName)
	os.Setenv("CLOUDSTACK_SSH_KEY_NAME", opt.SshKeyName)
	os.Setenv("CLOUDSTACK_AFFINITY_GROUP_ID", opt.AffinityGroupId)
	os.Setenv("CLOUDSTACK_VPC_NETWORK_NAME_WITH_CUSTOM_OFFERING", opt.NetworkName)
	os.Setenv("CLOUDSTACK_CUSTOM_VPC_NETWORK_OFFERING_NAME", opt.NetworkOffering)
	os.Setenv("CLOUDSTACK_VPC_CIDR", opt.VpcCidr)
	os.Setenv("CLOUDSTACK_VPC_NAME_WITH_CUSTOM_OFFERING", opt.VpcName)
	os.Setenv("CLOUDSTACK_CUSTOM_VPC_OFFERING_NAME", opt.VpcOffering)
	os.Setenv("CLOUDSTACK_GATEWAY", opt.Gateway)
	os.Setenv("CLOUDSTACK_NETMASK", opt.Netmask)

	templateOptions := client.GetClusterTemplateOptions{
		Kubeconfig: client.Kubeconfig{
			Path: c.KubeconfigFile,
		},
		ClusterName:              opt.ClusterName,
		TargetNamespace:          opt.Namespace,
		KubernetesVersion:        opt.KubernetesVersion,
		ListVariablesOnly:        false,
		WorkerMachineCount:       &opt.WorkerMachineCount,
		ControlPlaneMachineCount: &opt.ControlPlaneMachineCount,
	}

	if opt.URL != "" {
		templateOptions.URLSource = &client.URLSourceOptions{
			URL: opt.URL,
		}
	}

	template, err := c.Client.GetClusterTemplate(context.Background(), templateOptions)
	if err != nil {
		return "", err
	}

	yaml, err := template.Yaml()
	if err != nil {
		return "", err
	}

	return string(yaml), nil
}

func (c *ClusterApiClient) GenerateAwsWorkloadClusterYaml(opt option.GenerateAwsWorkloadClusterOption) (string, error) {
	os.Setenv("AWS_REGION", opt.Region)
	os.Setenv("AWS_SSH_KEY_NAME", opt.SshKeyName)
	os.Setenv("KUBERNETES_VERSION", opt.KubernetesVersion)
	os.Setenv("FLAVOR", opt.Flavor)
	os.Setenv("WORKER_MACHINE_COUNT", fmt.Sprintf("%d", opt.WorkerMachineCount))

	var controlMachineCount int64 = 1
	templateOptions := client.GetClusterTemplateOptions{
		Kubeconfig: client.Kubeconfig{
			Path: c.KubeconfigFile,
		},
		ClusterName:              opt.ClusterName,
		TargetNamespace:          opt.Namespace,
		KubernetesVersion:        opt.KubernetesVersion,
		ListVariablesOnly:        false,
		WorkerMachineCount:       &opt.WorkerMachineCount,
		ControlPlaneMachineCount: &controlMachineCount,
	}

	if opt.URL != "" {
		templateOptions.URLSource = &client.URLSourceOptions{
			URL: opt.URL,
		}
	}

	template, err := c.Client.GetClusterTemplate(context.Background(), templateOptions)
	if err != nil {
		return "", err
	}

	yaml, err := template.Yaml()
	if err != nil {
		return "", err
	}

	return string(yaml), nil
}

func (c *ClusterApiClient) GenerateGkeWorkloadClusterYaml(opt option.GenerateGkeWorkloadClusterOption) (string, error) {
	os.Setenv("GCP_CONTROL_PLANE_MACHINE_TYPE", opt.ControlPlaneMachineType)
	os.Setenv("GCP_NODE_MACHINE_TYPE", opt.WorkerMachineType)
	os.Setenv("KUBERNETES_VERSION", opt.KubernetesVersion)
	os.Setenv("CLUSTER_NAME", opt.ClusterName)
	os.Setenv("GCP_NETWORK_NAME", opt.NetworkName)
	os.Setenv("GCP_PROJECT", opt.Project)
	os.Setenv("GCP_REGION", opt.Region)

	var controlMachineCount int64 = 1
	templateOptions := client.GetClusterTemplateOptions{
		Kubeconfig: client.Kubeconfig{
			Path: c.KubeconfigFile,
		},
		ClusterName:              opt.ClusterName,
		TargetNamespace:          opt.Namespace,
		KubernetesVersion:        opt.KubernetesVersion,
		ListVariablesOnly:        false,
		WorkerMachineCount:       &opt.WorkerMachineCount,
		ControlPlaneMachineCount: &controlMachineCount,
	}

	if opt.URL != "" {
		templateOptions.URLSource = &client.URLSourceOptions{
			URL: opt.URL,
		}
	}

	template, err := c.Client.GetClusterTemplate(context.Background(), templateOptions)
	if err != nil {
		return "", err
	}

	yaml, err := template.Yaml()
	if err != nil {
		return "", err
	}

	return string(yaml), nil
}

func (c *ClusterApiClient) GetWorkloadClusterKubeconfig(clusterName, namespace string) (*string, error) {
	opt := client.GetKubeconfigOptions{
		Kubeconfig:          client.Kubeconfig{Path: c.KubeconfigFile},
		WorkloadClusterName: clusterName,
		Namespace:           namespace,
	}

	out, err := c.Client.GetKubeconfig(context.Background(), opt)
	if err != nil {
		return nil, err
	}

	return &out, err
}

func (c *ClusterApiClient) ApplyYaml(yamlString string) error {
	var err error
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(yamlString)), 100)
	for {
		var rawObj runtime.RawExtension
		if err = decoder.Decode(&rawObj); err != nil {
			break
		}

		dri, unstructuredObj, err := c.createDynamicResourceInterface(rawObj, "apply")
		if err != nil {
			c.LabelSelector = nil
			return err
		}
		if dri == nil {
			continue
		}

		if _, err := (*dri).Create(context.Background(), unstructuredObj, metav1.CreateOptions{}); err != nil {
			c.LabelSelector = nil
			if strings.Contains(error.Error(err), ` already exists`) {
				continue
			}
			return err
		}
	}
	if err != io.EOF {
		c.LabelSelector = nil
		return err
	}

	c.LabelSelector = nil
	return nil
}

func (c *ClusterApiClient) DeleteYaml(yamlString string) error {
	var err error
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(yamlString)), 100)
	for {
		var rawObj runtime.RawExtension
		if err = decoder.Decode(&rawObj); err != nil {
			break
		}

		dri, unstructuredObj, err := c.createDynamicResourceInterface(rawObj, "delete")
		if err != nil {
			c.LabelSelector = nil
			if strings.Contains(error.Error(err), "no matches ") {
				continue
			}
			return err
		}
		if dri == nil {
			continue
		}

		if err := (*dri).Delete(context.Background(), unstructuredObj.GetName(), metav1.DeleteOptions{}); err != nil {
			c.LabelSelector = nil
			if strings.Contains(error.Error(err), ` not found`) || strings.Contains(error.Error(err), "no matches ") {
				continue
			}
			return err
		}
	}
	if err != io.EOF {
		c.LabelSelector = nil
		return err
	}

	c.LabelSelector = nil
	return nil
}

func (c *ClusterApiClient) ClusterApiReadiness() (bool, error) {
	namespaces := []string{
		"capi-kubeadm-bootstrap-system",
		"capi-kubeadm-control-plane-system",
		"capi-system",
	}
	readiness := true

	for _, ns := range namespaces {
		pods, err := c.Clientset.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return false, err
		}

		for _, pod := range pods.Items {
			for _, container := range pod.Status.ContainerStatuses {
				if !container.Ready {
					readiness = false
					return readiness, fmt.Errorf("container %s's status in pod %s is not ready", container.Name, pod.Name)
				}
			}

			if pod.Status.Phase != v1.PodRunning {
				readiness = false
				return readiness, fmt.Errorf("%s is currently %s", pod.Name, pod.Status.Phase)
			}
		}
	}

	return readiness, nil
}

func (c *ClusterApiClient) InfrastructureReadiness(infrastructure string) (bool, error) {
	namespace, ok := option.Namespaces[infrastructure]
	if !ok {
		return false, fmt.Errorf("this infrastructure is not available")
	}

	ns, err := c.Clientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	pods, err := c.Clientset.CoreV1().Pods(ns.Name).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	readiness := false
	for _, pod := range pods.Items {
		readiness = true
		for _, container := range pod.Status.ContainerStatuses {
			if !container.Ready {
				readiness = false
				return readiness, fmt.Errorf("container %s's status in pod %s is not ready", container.Name, pod.Name)
			}
		}

		if pod.Status.Phase != v1.PodRunning {
			readiness = false
			return readiness, fmt.Errorf("%s is currently %s", pod.Name, pod.Status.Phase)
		}
	}

	return readiness, nil
}

func (c *ClusterApiClient) createDynamicResourceInterface(rawObj runtime.RawExtension, action string) (*dynamic.ResourceInterface, *unstructured.Unstructured, error) {
	obj, gvk, err := yamlserializer.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, nil, err
	}

	unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}
	if c.LabelSelector != nil {
		match := false
		labels := unstructuredObj.GetLabels()
		selector := c.LabelSelector.MatchLabels
		for k, v := range selector {
			if val, ok := labels[k]; ok {
				if val == v {
					match = true
					break
				}
			}
		}
		if !match {
			return nil, nil, nil
		}
	}

	gr, err := restmapper.GetAPIGroupResources(c.Clientset.Discovery())
	if err != nil {
		return nil, nil, err
	}

	mapper := restmapper.NewDiscoveryRESTMapper(gr)
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, nil, err
	}

	if mapping.Resource.Resource == "lists" {
		if items, ok := unstructuredObj.Object["items"]; ok {
			mapItems := items.([]interface{})
			yamlResult := ""
			var err error
			for _, item := range mapItems {
				b, _ := yaml.Marshal(item)
				yamlResult = yamlResult + "---\n" + string(b)
			}
			switch action {
			case "apply":
				err = c.ApplyYaml(yamlResult)
			case "delete":
				err = c.DeleteYaml(yamlResult)
			}
			return nil, nil, err
		}
	}

	var dri dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		if unstructuredObj.GetNamespace() == "" {
			unstructuredObj.SetNamespace("default")
		}
		dri = c.DynamicInterface.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
	} else {
		dri = c.DynamicInterface.Resource(mapping.Resource)
	}

	return &dri, unstructuredObj, nil
}

func (c *ClusterApiClient) CreateSecret(secret v1.Secret) (*v1.Secret, error) {
	secretValue, err := c.Clientset.CoreV1().Secrets(secret.ObjectMeta.Namespace).Create(context.TODO(), &secret, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return secretValue, nil
}

func (c *ClusterApiClient) CreateDockerRegistrySecret(secretName, namespace string, args model.CreateDockerRegistrySecretArgs) (*v1.Secret, error) {
	secretObj := v1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: v1.SchemeGroupVersion.String(),
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        secretName,
			Namespace:   namespace,
			Annotations: args.Annotations,
		},
		Type: v1.SecretTypeDockerConfigJson,
		Data: map[string][]byte{},
	}

	dockerConfigAuth := model.DockerConfigEntry{
		Username: args.Username,
		Password: args.Password,
		Email:    args.Email,
		Auth:     base64.StdEncoding.EncodeToString([]byte(args.Username + ":" + args.Password)),
	}

	dockerConfigJSON := model.DockerConfigJSON{
		Auths: map[string]model.DockerConfigEntry{args.Server: dockerConfigAuth},
	}

	b, err := json.Marshal(dockerConfigJSON)
	if err != nil {
		return nil, err
	}

	secretObj.Data[v1.DockerConfigJsonKey] = b

	secretValue, err := c.Clientset.CoreV1().Secrets(secretObj.ObjectMeta.Namespace).Create(context.TODO(), &secretObj, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return secretValue, nil
}

func (c *ClusterApiClient) CreateNamespace(namespace string) (*v1.Namespace, error) {
	return c.Clientset.CoreV1().Namespaces().Create(context.Background(), &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}, metav1.CreateOptions{})
}

func (c *ClusterApiClient) AddLabelNamespace(namespace, label, value string) (*v1.Namespace, error) {
	type PatchStringValue struct {
		Op    string `json:"op"`
		Path  string `json:"path"`
		Value string `json:"value"`
	}
	payload := []PatchStringValue{{
		Op:    "add",
		Path:  "/metadata/labels/" + label,
		Value: value,
	}}
	b, _ := json.Marshal(payload)
	return c.Clientset.CoreV1().Namespaces().Patch(context.Background(), namespace, types.JSONPatchType, b, metav1.PatchOptions{})
}

func (c *ClusterApiClient) GetService(serviceName, namespace string) (*v1.Service, error) {
	service, err := c.Clientset.CoreV1().Services(namespace).Get(context.Background(), serviceName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (c *ClusterApiClient) GetSecret(secretName, namespace string) (*v1.Secret, error) {
	secret, err := c.Clientset.CoreV1().Secrets(namespace).Get(context.Background(), secretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (c *ClusterApiClient) PatchServiceAccount(name, namespace string, patch []byte) (*v1.ServiceAccount, error) {
	sa, err := c.Clientset.CoreV1().ServiceAccounts(namespace).Patch(context.Background(), name, types.StrategicMergePatchType, patch, metav1.PatchOptions{})

	return sa, err
}

func (c *ClusterApiClient) PatchConfigMap(name, namespace string, patch []byte) (*v1.ConfigMap, error) {
	cm, err := c.Clientset.CoreV1().ConfigMaps(namespace).Patch(context.Background(), name, types.StrategicMergePatchType, patch, metav1.PatchOptions{})

	return cm, err
}

func (c *ClusterApiClient) UpdateSecret(namespace string, secret *v1.Secret) (*v1.Secret, error) {
	sc, err := c.Clientset.CoreV1().Secrets(namespace).Update(context.Background(), secret, metav1.UpdateOptions{})

	return sc, err
}

func (c *ClusterApiClient) DeleteCluster(clusterName, namespace string) ([]byte, error) {
	return c.Clientset.RESTClient().Delete().
		AbsPath("apis/cluster.x-k8s.io/v1beta1/namespaces/"+namespace+"/clusters/"+clusterName).
		VersionedParams(&metav1.GetOptions{}, metav1.ParameterCodec).
		DoRaw(context.TODO())
}

func (c *ClusterApiClient) GetKNativeRevision(revisionName, namespace string) ([]byte, error) {
	return c.Clientset.RESTClient().Get().
		AbsPath("apis/serving.knative.dev/v1/namespaces/"+namespace+"/revisions/"+revisionName).
		VersionedParams(&metav1.GetOptions{}, metav1.ParameterCodec).
		DoRaw(context.TODO())
}

func (c *ClusterApiClient) GetKNativeConfiguration(configurationName, namespace string) ([]byte, error) {
	return c.Clientset.RESTClient().Get().
		AbsPath("apis/serving.knative.dev/v1/namespaces/"+namespace+"/configurations/"+configurationName).
		VersionedParams(&metav1.GetOptions{}, metav1.ParameterCodec).
		DoRaw(context.TODO())
}

func (c *ClusterApiClient) GetClusterK8sResource(clusterName, namespace string) ([]byte, error) {
	return c.Clientset.RESTClient().Get().
		AbsPath("apis/cluster.x-k8s.io/v1beta1/namespaces/"+namespace+"/clusters/"+clusterName).
		VersionedParams(&metav1.GetOptions{}, metav1.ParameterCodec).
		DoRaw(context.TODO())
}

func (c *ClusterApiClient) GetDeployment(deploymentName, namespace string) (*appsv1.Deployment, error) {
	deployment, err := c.Clientset.AppsV1().Deployments(namespace).Get(context.Background(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return deployment, nil
}

func (c *ClusterApiClient) RestartDeployment(deploymentName, namespace string) (*appsv1.Deployment, error) {
	deployment, err := c.Clientset.AppsV1().Deployments(namespace).Get(context.Background(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var restartTimestampFound bool
	container := &deployment.Spec.Template.Spec.Containers[0]
	for i, env := range container.Env {
		if env.Name == "RESTART_TIMESTAMP" {
			container.Env[i].Value = strconv.FormatInt(time.Now().Unix(), 10)
			restartTimestampFound = true
			break
		}
	}

	if !restartTimestampFound {
		container.Env = append(container.Env,
			v1.EnvVar{
				Name:  "RESTART_TIMESTAMP",
				Value: strconv.FormatInt(time.Now().Unix(), 10),
			})
	}

	_, err = c.Clientset.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		panic(err.Error())
	}

	return deployment, nil
}

// https://stackoverflow.com/questions/65927298/patching-a-pvc-using-go-client
func (c *ClusterApiClient) UpdateClusterK8sResourceAnnotations(clusterName, namespace string, patchValues interface{}) (*unstructured.Unstructured, error) {
	patch := []struct {
		Op    string      `json:"op"`
		Path  string      `json:"path"`
		Value interface{} `json:"value"`
	}{
		{
			Op:    "add",
			Path:  "/metadata/annotations",
			Value: struct{}{},
		},
		{
			Op:    "add",
			Path:  "/metadata/annotations/accountId",
			Value: "this-is-account-id",
		},
		{
			Op:    "add",
			Path:  "/metadata/annotations/region",
			Value: "banten-1",
		},
		{
			Op:    "add",
			Path:  "/metadata/annotations/vendor",
			Value: "biznet",
		},
	}

	b, _ := json.Marshal(patch)

	resource := schema.GroupVersionResource{
		Group:    "cluster.x-k8s.io",
		Version:  "v1beta1",
		Resource: "clusters",
	}

	return c.DynamicInterface.Resource(resource).
		Namespace(namespace).
		Patch(context.TODO(), clusterName, types.JSONPatchType, b, metav1.PatchOptions{})
}

func (c *ClusterApiClient) ExecuteNodeShellCommand(nodeName, command string) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	suffix := make([]byte, 4)
	for i := range suffix {
		suffix[i] = charset[seededRand.Intn(len(charset))]
	}
	var (
		terminationGracePeriodSeconds int64 = 0
		privilegedSecurityContext     bool  = true
		jobName                             = "node-shell-job-" + nodeName + "-" + string(suffix)
		namespace                           = "kube-system"
	)

	log.Printf("Executing command %s on %s", command, nodeName)

	var ttlSecondsAfterFinished int32 = 150

	jobSpec := batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: &ttlSecondsAfterFinished,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					RestartPolicy:                 "Never",
					TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
					HostPID:                       true,
					HostIPC:                       true,
					HostNetwork:                   true,
					Tolerations: []v1.Toleration{
						{
							Operator: "Exists",
						},
					},
					Containers: []v1.Container{
						{
							Name:  "shell",
							Image: "ubuntu:24.04",
							SecurityContext: &v1.SecurityContext{
								Privileged: &privilegedSecurityContext,
							},
							Command: []string{"sh"},
							Args: []string{
								"-c",
								"nsenter -t 1 -m -u -i -n sh -c '" + command + "'",
							},
						},
					},
					NodeSelector: map[string]string{
						"kubernetes.io/hostname": nodeName,
					},
				},
			},
		},
	}

	// Create the Job
	_, err := c.Clientset.BatchV1().Jobs(namespace).Create(context.Background(), &jobSpec, metav1.CreateOptions{})
	if err != nil {
		if !strings.Contains(err.Error(), " already exists") {
			return "", err
		}
	}

	// Wait for the Job to complete
	err = wait.PollImmediate(3*time.Second, 2*time.Minute, func() (bool, error) {
		job, err := c.Clientset.BatchV1().Jobs(namespace).Get(context.Background(), jobName, metav1.GetOptions{})
		if err != nil {
			if k8serrors.IsNotFound(err) {
				fmt.Printf("Job %s not created yet...\n", jobName)
				return false, nil // Continue polling
			}
			return false, err
		}

		if job.Status.Succeeded > 0 {
			return true, nil
		}

		if job.Status.Failed > 0 {
			err := fmt.Errorf("Job %s is failed: %v", jobName, job.Status)
			fmt.Println(err)
			return false, err
		}

		fmt.Printf("Waiting for job %s to complete...\n", jobName)
		return false, nil
	})
	if err != nil {
		return "", err
	}

	// Retrieve logs from the job's pod
	pods, err := c.Clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("job-name=%s", jobName),
	})
	if err != nil || len(pods.Items) == 0 {
		return "", fmt.Errorf("failed to retrieve job pod: %v", err)
	}

	sort.Slice(pods.Items, func(i, j int) bool {
		return pods.Items[i].CreationTimestamp.Time.After(pods.Items[j].CreationTimestamp.Time)
	})

	podName := pods.Items[0].Name
	logs, err := c.GetPodLogs(namespace, podName)
	if err != nil {
		return "", err
	}

	// Cleanup the job after execution
	deletePolicy := metav1.DeletePropagationForeground
	err = c.Clientset.BatchV1().Jobs(namespace).Delete(context.Background(), jobName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		fmt.Printf("Failed to delete job %s: %v\n", jobName, err)
	}

	return logs, nil
}

func (c *ClusterApiClient) GetPodLogs(namespace, podName string) (string, error) {
	const maxAttempts = 3
	const retryDelay = 5 * time.Second

	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		req := c.Clientset.CoreV1().Pods(namespace).GetLogs(podName, &v1.PodLogOptions{})
		logs, err := req.Stream(ctx)
		if err != nil {
			log.Printf("Attempt %d: Error getting stream for pod %s: %v\n", attempt, podName, err)
			lastErr = err
			time.Sleep(retryDelay)
			continue
		}
		defer logs.Close()

		var buf bytes.Buffer
		_, err = io.Copy(&buf, logs)
		if err != nil {
			log.Printf("Attempt %d: Error copying log buffer for pod %s: %v\n", attempt, podName, err)
			lastErr = err
			time.Sleep(retryDelay)
			continue
		}

		// Return logs on successful retrieval
		return buf.String(), nil
	}

	// Return the last error after exhausting all attempts
	log.Printf("Failed to get logs for pod %s after %d attempts: %v\n", podName, maxAttempts, lastErr)
	return "", fmt.Errorf("failed to get logs for pod %s after %d attempts: %w", podName, maxAttempts, lastErr)
}

func (c *ClusterApiClient) DescribeCluster(clusterName, namespace string) (*tree.ObjectTree, error) {
	objTree, err := c.Client.DescribeCluster(context.Background(), client.DescribeClusterOptions{
		Namespace:   clusterName,
		ClusterName: namespace,
		Kubeconfig:  client.Kubeconfig{Path: c.KubeconfigFile},
	})

	return objTree, err
}

func (c *ClusterApiClient) GetWorkloadClusterMachines(clusterName, namespace string) ([]model.CRDResource, error) {
	b, err := c.Clientset.RESTClient().Get().
		AbsPath("apis/cluster.x-k8s.io/v1beta1/namespaces/"+namespace+"/machines").
		VersionedParams(&metav1.GetOptions{}, metav1.ParameterCodec).
		DoRaw(context.TODO())
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	b, _ = json.Marshal(m["items"])

	machines := []model.CRDResource{}
	if err := json.Unmarshal(b, &machines); err != nil {
		return nil, err
	}

	filteredMachines := []model.CRDResource{}
	for _, machine := range machines {
		if machine.Metadata.Labels["cluster.x-k8s.io/cluster-name"] == clusterName {
			filteredMachines = append(filteredMachines, machine)
		}
	}

	return filteredMachines, nil
}
