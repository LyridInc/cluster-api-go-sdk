package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/LyridInc/cluster-api-go-sdk/model"
	"github.com/LyridInc/cluster-api-go-sdk/option"
	yamlmodel "github.com/LyridInc/cluster-api-go-sdk/yaml"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	yamlserializer "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
)

type (
	OpenstackClient struct {
		NetworkEndpoint string
		AuthEndpoint    string
		AuthToken       string
		ProjectId       string
	}

	OpenstackCredential struct {
		ApplicationCredentialName   string `json:"name"`
		ApplicationCredentialId     string `json:"id"`
		ApplicationCredentialSecret string `json:"secret"`
	}

	OpenstackPassword struct {
		User map[string]interface{} `json:"user"`
	}

	OpenstackIdentity struct {
		Methods               []string            `json:"methods"`
		Password              OpenstackPassword   `json:"password"`
		ApplicationCredential OpenstackCredential `json:"application_credential"`
	}

	OpenstackAuth struct {
		Identity OpenstackIdentity `json:"identity"`
	}
)

func (c *OpenstackClient) Authenticate(auth OpenstackAuth) (*http.Response, error) {
	token := os.Getenv("OS_TOKEN")
	if token != "" {
		c.AuthToken = token
		response, err := c.CheckAuthToken()
		if err != nil {
			return nil, err
		}
		responseError, ok := response["error"]
		if !ok {
			return nil, nil
		}
		responseErrorMap := responseError.(map[string]interface{})
		code, ok := responseErrorMap["code"]
		if ok {
			statusCode := code.(float64)
			if statusCode != 401 {
				return nil, fmt.Errorf("%v", responseError)
			}
		}
	}

	url := c.AuthEndpoint + "/v3/auth/tokens"
	b, _ := json.Marshal(map[string]interface{}{
		"auth": auth,
	})
	requestBody := b

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return response, err
	}
	defer response.Body.Close()

	for key, value := range response.Header {
		if key == "X-Subject-Token" && len(value) > 0 {
			c.AuthToken = value[0]
			os.Setenv("OS_TOKEN", c.AuthToken)
			break
		}
	}

	return response, nil
}

func (c *OpenstackClient) CheckAuthToken() (map[string]interface{}, error) {
	url := c.AuthEndpoint + "/v3/auth/tokens"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Token", c.AuthToken)
	request.Header.Set("X-Subject-Token", c.AuthToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	jsonResponse := map[string]interface{}{}
	body, _ := io.ReadAll(response.Body)
	json.Unmarshal(body, &jsonResponse)

	return jsonResponse, nil
}

func (c *OpenstackClient) GetProjectQuotas() (*model.QuotasResponse, error) {
	url := c.NetworkEndpoint + "/v2.0/quotas/" + c.ProjectId + "/details.json"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Token", c.AuthToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	quotas := model.QuotasResponse{}
	json.Unmarshal(body, &quotas)

	return &quotas, nil
}

func (c *OpenstackClient) ValidateQuotas() (bool, error) {
	quotas, err := c.GetProjectQuotas()
	if err != nil {
		return false, err
	}

	floatingIpQuota := quotas.Quota.FloatingIp
	if floatingIpQuota.Limit != -1 && floatingIpQuota.Used >= floatingIpQuota.Limit {
		return false, fmt.Errorf("floating ip quota limit exceeded")
	}
	networkQuota := quotas.Quota.Network
	if networkQuota.Limit != -1 && networkQuota.Used >= networkQuota.Limit {
		return false, fmt.Errorf("network quota limit exceeded")
	}
	routerQuota := quotas.Quota.Router
	if routerQuota.Limit != -1 && routerQuota.Used >= routerQuota.Limit {
		return false, fmt.Errorf("router quota limit exceeded")
	}
	securityGroupQuota := quotas.Quota.SecurityGroup
	if securityGroupQuota.Limit != -1 && securityGroupQuota.Used+2 >= securityGroupQuota.Limit {
		return false, fmt.Errorf("security group quota limit exceeded")
	}
	securityGroupRuleQuota := quotas.Quota.SecurityGroupRule
	if securityGroupRuleQuota.Limit != -1 && securityGroupRuleQuota.Used+10 >= securityGroupRuleQuota.Limit {
		return false, fmt.Errorf("security group rules quota limit exceeded")
	}
	portQuota := quotas.Quota.Port
	if portQuota.Limit != -1 && portQuota.Used+3 >= portQuota.Limit {
		return false, fmt.Errorf("port quota limit exceeded")
	}

	return true, nil
}

func (c *OpenstackClient) UpdateYamlManifest(yamlString string, opt option.ManifestOption) (string, error) {
	var (
		err        error
		yamlResult string
	)
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(yamlString)), 100)
	for {
		var rawObj runtime.RawExtension
		if err = decoder.Decode(&rawObj); err != nil {
			break
		}

		obj, _, err := yamlserializer.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if err != nil {
			return "", err
		}

		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			return "", err
		}

		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}
		unstructuredObj = UpdateUnstructuredObject(unstructuredObj, opt)

		x, _ := yaml.Marshal(unstructuredObj.Object)
		yamlResult = yamlResult + "---\n" + string(x)
	}

	if err != io.EOF {
		return "", err
	}

	return yamlResult, nil
}

func UpdateUnstructuredObject(unstructuredObj *unstructured.Unstructured, opt option.ManifestOption) *unstructured.Unstructured {
	apiVersion := unstructuredObj.GetAPIVersion()
	kind := unstructuredObj.GetKind()

	spec, ok := unstructuredObj.Object["spec"]
	var (
		specByte      []byte
		objByte       []byte
		specInterface interface{}
	)
	if ok {
		specByte, _ = json.Marshal(spec)
	}
	unstructuredObjByte, _ := json.Marshal(unstructuredObj.Object)
	kindMap := map[string]interface{}{}

	if strings.HasPrefix(apiVersion, "storage.k8s.io") && kind == "StorageClass" {
		storageClass := yamlmodel.StorageClass{}
		json.Unmarshal(unstructuredObjByte, &storageClass)
		if opt.StorageClassKindOption.Parameters != nil {
			storageClass.Parameters = opt.StorageClassKindOption.Parameters
		}
		objByte, _ = json.Marshal(storageClass)
	} else if kind == "Secret" {
		secret := yamlmodel.Secret{}
		json.Unmarshal(unstructuredObjByte, &secret)
		if opt.SecretKindOption.Data != nil {
			secret.Data = opt.SecretKindOption.Data
		}
		if opt.SecretKindOption.Metadata != nil {
			secret.Metadata = opt.SecretKindOption.Metadata
		}
		objByte, _ = json.Marshal(secret)
	} else if kind == "PersistentVolumeClaim" {
		pvc := yamlmodel.PersistentVolumeClaim{}
		json.Unmarshal(unstructuredObjByte, &pvc)
		if opt.PersistentVolumeClaimKindOption.Metadata != nil {
			pvc.Metadata = opt.PersistentVolumeClaimKindOption.Metadata
		}
		objByte, _ = json.Marshal(pvc)
		if ok {
			if opt.PersistentVolumeClaimKindOption.Storage != "" {
				pvc.Spec.Resources.Requests.Storage = opt.PersistentVolumeClaimKindOption.Storage
			}
			if opt.PersistentVolumeClaimKindOption.StorageClassName != "" {
				pvc.Spec.StorageClassName = opt.PersistentVolumeClaimKindOption.StorageClassName
			}
			if opt.PersistentVolumeClaimKindOption.VolumeMode != "" {
				pvc.Spec.VolumeMode = opt.PersistentVolumeClaimKindOption.VolumeMode
			}
			specInterface = pvc.Spec
		}
	}

	if len(objByte) > 0 {
		json.Unmarshal(objByte, &kindMap)
		unstructuredObj.Object = kindMap
	}

	if ok {
		if strings.HasPrefix(apiVersion, "infrastructure.cluster.x-k8s.io") && kind == "OpenStackCluster" {
			infrastructureSpec := yamlmodel.InfrastructureSpec{}
			json.Unmarshal(specByte, &infrastructureSpec)
			if opt.InfrastructureKindSpecOption.NodeCidr != "" {
				infrastructureSpec.NodeCidr = opt.InfrastructureKindSpecOption.NodeCidr
			}
			infrastructureSpec.AllowAllInClusterTraffic = opt.InfrastructureKindSpecOption.AllowAllInClusterTraffic
			specInterface = infrastructureSpec
		} else if strings.HasPrefix(apiVersion, "cluster.x-k8s.io") && kind == "Cluster" {
			clusterSpec := yamlmodel.ClusterSpec{}
			json.Unmarshal(specByte, &clusterSpec)
			if opt.ClusterKindSpecOption.CidrBlocks != nil {
				clusterSpec.ClusterNetwork.Pods.CidrBlocks = opt.ClusterKindSpecOption.CidrBlocks
			}
			specInterface = clusterSpec
		} else if strings.HasPrefix(apiVersion, "apps") && kind == "DaemonSet" {
			daemonSetSpec := yamlmodel.DaemonSetSpec{}
			json.Unmarshal(specByte, &daemonSetSpec)
			if opt.DaemonSetKindOption.VolumeSecretName != "" {
				for i, v := range daemonSetSpec.Template.Spec.Volumes {
					vv := v.(map[string]interface{})
					if vv["name"] == "cloud-config-volume" || vv["name"] == "secret-cinderplugin" {
						vvs := vv["secret"].(map[string]interface{})
						vvs["secretName"] = opt.DaemonSetKindOption.VolumeSecretName
						vv["secret"] = vvs
						daemonSetSpec.Template.Spec.Volumes[i] = vv
					}
				}
			}
			specInterface = daemonSetSpec
		} else if strings.HasPrefix(apiVersion, "apps") && kind == "Deployment" {
			deploymentSpec := yamlmodel.DeploymentSpec{}
			json.Unmarshal(specByte, &deploymentSpec)
			if opt.DeploymentKindOption.VolumeSecretName != "" {
				for i, v := range deploymentSpec.Template.Spec.Volumes {
					vv := v.(map[string]interface{})
					if vv["name"] == "cloud-config-volume" || vv["name"] == "secret-cinderplugin" {
						vvs := vv["secret"].(map[string]interface{})
						vvs["secretName"] = opt.DeploymentKindOption.VolumeSecretName
						vv["secret"] = vvs
						deploymentSpec.Template.Spec.Volumes[i] = vv
					}
				}
			}
			specInterface = deploymentSpec
		}

		if !reflect.ValueOf(&specInterface).Elem().IsNil() {
			m := map[string]interface{}{}
			specByte, _ = json.Marshal(specInterface)
			json.Unmarshal(specByte, &m)
			removeNulls(m)
			unstructuredObj.Object["spec"] = m
		}
	}

	return unstructuredObj
}

func CleanUpMapFromNullValues(m *map[string]interface{}) {
	for key, value := range *m {
		if value == nil {
			delete(*m, key)
		}
	}
}

func removeNulls(m map[string]interface{}) {
	val := reflect.ValueOf(m)
	for _, e := range val.MapKeys() {
		v := val.MapIndex(e)
		if v.IsNil() {
			delete(m, e.String())
			continue
		}
		switch t := v.Interface().(type) {
		// If key is a JSON object (Go Map), use recursion to go deeper
		case map[string]interface{}:
			removeNulls(t)

		case string:
			if string(t) == "" {
				delete(m, e.String())
				continue
			}

		case []interface{}:
			for _, vv := range t {
				vvr := reflect.ValueOf(vv)
				switch tvv := vvr.Interface().(type) {
				case map[string]interface{}:
					removeNulls(tvv)
				}
			}
		}
	}
}
