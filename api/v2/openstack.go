package v2

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/tools/clientcmd"
	infrav1alpha1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha1"
	infrav1beta1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type OpenstackK8sClient struct {
	K8sClient client.Client
}

func NewOpenstackK8sClient(kubeconfigPath string) (*OpenstackK8sClient, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("Error BuildConfigFromFlags: %v", err)
	}

	scheme := runtime.NewScheme()
	_ = infrav1beta1.AddToScheme(scheme)
	_ = infrav1alpha1.AddToScheme(scheme)

	k8sClient, err := client.New(cfg, client.Options{Scheme: scheme})
	if err != nil {
		return nil, fmt.Errorf("Error initializing kubernetes client: %v", err)
	}

	return &OpenstackK8sClient{
		K8sClient: k8sClient,
	}, nil
}

func NewOpenstackK8sClientFromKubeconfig(kubeconfig []byte) (*OpenstackK8sClient, error) {
	rawConfig, err := clientcmd.Load(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse kubeconfig: %v", err)
	}

	clientConfig := clientcmd.NewDefaultClientConfig(*rawConfig, &clientcmd.ConfigOverrides{})
	cfg, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("Error BuildConfigFromFlags: %v", err)
	}

	scheme := runtime.NewScheme()
	_ = infrav1beta1.AddToScheme(scheme)
	_ = infrav1alpha1.AddToScheme(scheme)

	k8sClient, err := client.New(cfg, client.Options{Scheme: scheme})
	if err != nil {
		return nil, fmt.Errorf("Error initializing kubernetes client: %v", err)
	}

	return &OpenstackK8sClient{
		K8sClient: k8sClient,
	}, nil
}

func (c *OpenstackK8sClient) GetClusterCRD(name, namespace string) (*infrav1beta1.OpenStackCluster, error) {
	crd := &infrav1beta1.OpenStackCluster{}
	key := client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}

	if err := c.K8sClient.Get(context.Background(), key, crd); err != nil {
		return nil, fmt.Errorf("Error get OpenStackCluster CRD: %v", err)
	}

	return crd, nil
}

func (c *OpenstackK8sClient) AddClusterCRDAnnotation(name, namespace, field, value string) (*infrav1beta1.OpenStackCluster, error) {
	crd, err := c.GetClusterCRD(name, namespace)
	if err != nil {
		return nil, err
	}

	if crd.Annotations == nil {
		crd.Annotations = map[string]string{}
	}

	crd.Annotations[field] = value
	if err := c.K8sClient.Update(context.Background(), crd); err != nil {
		return nil, fmt.Errorf("Add annotation - Error update OpenStackCluster CRD: %v", err)
	}

	return crd, nil
}

func (c *OpenstackK8sClient) RemoveClusterCRDAnnotation(name, namespace, field string) (*infrav1beta1.OpenStackCluster, error) {
	crd, err := c.GetClusterCRD(name, namespace)
	if err != nil {
		return nil, err
	}

	if crd.Annotations == nil {
		crd.Annotations = map[string]string{}
	} else {
		delete(crd.Annotations, field)
	}

	if err := c.K8sClient.Update(context.Background(), crd); err != nil {
		return nil, fmt.Errorf("Remove annotation - Error update OpenStackCluster CRD: %v", err)
	}

	return crd, nil
}

func (c *OpenstackK8sClient) GetMachineCRDsByClusterName(clusterName, namespace string) (*infrav1beta1.OpenStackMachineList, error) {
	req, err := labels.NewRequirement("cluster.x-k8s.io/cluster-name", selection.Equals, []string{clusterName})
	if err != nil {
		return nil, fmt.Errorf("Failed to build label selector: %v", err)
	}

	selector := labels.NewSelector().Add(*req)
	list := &infrav1beta1.OpenStackMachineList{}
	listOpts := &client.ListOptions{
		Namespace:     namespace,
		LabelSelector: selector,
	}
	if err := c.K8sClient.List(context.Background(), list, listOpts); err != nil {
		return nil, fmt.Errorf("Error get OpenStackMachineList CRD: %v", err)
	}

	return list, nil
}

func (c *OpenstackK8sClient) GetMachineCRD(name, namespace string) (*infrav1beta1.OpenStackMachine, error) {
	crd := &infrav1beta1.OpenStackMachine{}
	key := client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}

	if err := c.K8sClient.Get(context.Background(), key, crd); err != nil {
		return nil, fmt.Errorf("Error get OpenStackMachine CRD: %v", err)
	}

	return crd, nil
}

func (c *OpenstackK8sClient) AddMachineCRDAnnotation(name, namespace, field, value string) (*infrav1beta1.OpenStackMachine, error) {
	crd, err := c.GetMachineCRD(name, namespace)
	if err != nil {
		return nil, err
	}

	if crd.Annotations == nil {
		crd.Annotations = map[string]string{}
	}

	crd.Annotations[field] = value
	if err := c.K8sClient.Update(context.Background(), crd); err != nil {
		return nil, fmt.Errorf("Add annotation - Error update OpenStackMachine CRD: %v", err)
	}

	return crd, nil
}

func (c *OpenstackK8sClient) RemoveMachineCRDAnnotation(name, namespace, field string) (*infrav1beta1.OpenStackMachine, error) {
	crd, err := c.GetMachineCRD(name, namespace)
	if err != nil {
		return nil, err
	}

	if crd.Annotations == nil {
		crd.Annotations = map[string]string{}
	} else {
		delete(crd.Annotations, field)
	}

	if err := c.K8sClient.Update(context.Background(), crd); err != nil {
		return nil, fmt.Errorf("Remove annotation - Error update OpenStackMachine CRD: %v", err)
	}

	return crd, nil
}

func (c *OpenstackK8sClient) GetServerCRDsByClusterName(clusterName, namespace string) (*infrav1alpha1.OpenStackServerList, error) {
	req, err := labels.NewRequirement("cluster.x-k8s.io/cluster-name", selection.Equals, []string{clusterName})
	if err != nil {
		return nil, fmt.Errorf("Failed to build label selector: %v", err)
	}

	selector := labels.NewSelector().Add(*req)
	list := &infrav1alpha1.OpenStackServerList{}
	listOpts := &client.ListOptions{
		Namespace:     namespace,
		LabelSelector: selector,
	}
	if err := c.K8sClient.List(context.Background(), list, listOpts); err != nil {
		return nil, fmt.Errorf("Error get OpenStackServerList CRD: %v", err)
	}

	return list, nil
}

func (c *OpenstackK8sClient) GetServerCRD(name, namespace string) (*infrav1alpha1.OpenStackServer, error) {
	crd := &infrav1alpha1.OpenStackServer{}
	key := client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}

	if err := c.K8sClient.Get(context.Background(), key, crd); err != nil {
		return nil, fmt.Errorf("Error get OpenStackServer CRD: %v", err)
	}

	return crd, nil
}

func (c *OpenstackK8sClient) AddServerCRDAnnotation(name, namespace, field, value string) (*infrav1alpha1.OpenStackServer, error) {
	crd, err := c.GetServerCRD(name, namespace)
	if err != nil {
		return nil, err
	}

	if crd.Annotations == nil {
		crd.Annotations = map[string]string{}
	}

	crd.Annotations[field] = value
	if err := c.K8sClient.Update(context.Background(), crd); err != nil {
		return nil, fmt.Errorf("Add annotation - Error update OpenStackServer CRD: %v", err)
	}

	return crd, nil
}

func (c *OpenstackK8sClient) RemoveServerCRDAnnotation(name, namespace, field string) (*infrav1alpha1.OpenStackServer, error) {
	crd, err := c.GetServerCRD(name, namespace)
	if err != nil {
		return nil, err
	}

	if crd.Annotations == nil {
		crd.Annotations = map[string]string{}
	} else {
		delete(crd.Annotations, field)
	}

	if err := c.K8sClient.Update(context.Background(), crd); err != nil {
		return nil, fmt.Errorf("Remove annotation - Error update OpenStackServer CRD: %v", err)
	}

	return crd, nil
}
