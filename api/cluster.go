package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/LyridInc/cluster-api-go-sdk/option"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	yamlserializer "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
)

type (
	ClusterApiClient struct {
		Client           client.Client
		Clientset        *kubernetes.Clientset
		DynamicInterface dynamic.Interface
		InitOptions      client.InitOptions
		ConfigFile       string
		KubeconfigFile   string
	}
)

func NewClusterApiClient(configFile, kubeconfigFile string) *ClusterApiClient {
	cl, err := client.New(configFile)
	if err != nil {
		log.Fatal("Client config error:", err)
		return nil
	}

	conf, err := clientcmd.BuildConfigFromFlags("", kubeconfigFile)
	if err != nil {
		log.Fatal("Build config from flags error:", err)
		return nil
	}

	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		log.Fatal("Clientset config error:", err)
		return nil
	}

	dd, err := dynamic.NewForConfig(conf)
	if err != nil {
		log.Fatal("Dynamic interface config error:", err)
		return nil
	}

	return &ClusterApiClient{
		Client:           cl,
		Clientset:        clientset,
		DynamicInterface: dd,
		ConfigFile:       configFile,
		KubeconfigFile:   kubeconfigFile,
	}
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

	result, err := c.Client.Init(c.InitOptions)
	if err != nil {
		return nil, err
	}

	if ready, err := c.InfrastructureReadiness(infrastructure); !ready || err != nil {
		return nil, err
	}
	return result, nil
}

func (c *ClusterApiClient) DeleteInfrastructure(infrastructure string) error {
	return c.Client.Delete(client.DeleteOptions{
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
		TargetNamespace:          "",
		KubernetesVersion:        opt.KubernetesVersion,
		ListVariablesOnly:        false,
		WorkerMachineCount:       &opt.WorkerMachineCount,
		ControlPlaneMachineCount: &opt.ControlPlaneMachineCount,
		ProviderRepositorySource: &client.ProviderRepositorySourceOptions{
			InfrastructureProvider: opt.InfrastructureProvider,
			Flavor:                 opt.Flavor,
		},
	}
	template, err := c.Client.GetClusterTemplate(templateOptions)
	if err != nil {
		return "", err
	}

	yaml, err := template.Yaml()
	if err != nil {
		return "", err
	}

	return string(yaml), nil
}

func (c *ClusterApiClient) Apply(yamlString string) error {
	var err error
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(yamlString)), 100)
	for {
		var rawObj runtime.RawExtension
		if err = decoder.Decode(&rawObj); err != nil {
			break
		}

		obj, gvk, err := yamlserializer.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if err != nil {
			return err
		}

		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			return err
		}

		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

		gr, err := restmapper.GetAPIGroupResources(c.Clientset.Discovery())
		if err != nil {
			return err
		}

		mapper := restmapper.NewDiscoveryRESTMapper(gr)
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return err
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

		if _, err := dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{}); err != nil {
			return err
		}
	}
	if err != io.EOF {
		return err
	}

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
