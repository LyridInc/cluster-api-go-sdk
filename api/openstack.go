package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
		ApplicationCredentialName   string
		ApplicationCredentialId     string
		ApplicationCredentialSecret string
	}
)

func (c *OpenstackClient) Authenticate(credential OpenstackCredential) error {
	token := os.Getenv("OS_TOKEN")
	if token != "" {
		c.AuthToken = token
		response, err := c.CheckAuthToken()
		if err != nil {
			return err
		}
		responseError, ok := response["error"]
		if !ok {
			return nil
		}
		responseErrorMap := responseError.(map[string]interface{})
		code, ok := responseErrorMap["code"]
		if ok {
			statusCode := code.(float64)
			if statusCode != 401 {
				return fmt.Errorf("%v", responseError)
			}
		}
	}

	url := c.AuthEndpoint + "/v3/auth/tokens"
	requestBody := []byte(`{
		"auth": {
			"identity": {
				"methods": ["application_credential"],
				"application_credential": {
					"id": "` + credential.ApplicationCredentialId + `",
					"name": "` + credential.ApplicationCredentialName + `",
					"secret": "` + credential.ApplicationCredentialSecret + `"
				}
			}
		}
	}`)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	for key, value := range response.Header {
		if key == "X-Subject-Token" && len(value) > 0 {
			c.AuthToken = value[0]
			os.Setenv("OS_TOKEN", c.AuthToken)
			break
		}
	}

	return nil
}

func (c *OpenstackClient) CheckAuthToken() (map[string]interface{}, error) {
	url := c.AuthEndpoint + "/v3/auth/tokens"
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
	if floatingIpQuota.Used >= floatingIpQuota.Limit {
		return false, fmt.Errorf("floating ip quota limit exceeded")
	}
	networkQuota := quotas.Quota.Network
	if networkQuota.Used >= networkQuota.Limit {
		return false, fmt.Errorf("network quota limit exceeded")
	}
	routerQuota := quotas.Quota.Router
	if routerQuota.Used >= routerQuota.Limit {
		return false, fmt.Errorf("router quota limit exceeded")
	}
	securityGroupQuota := quotas.Quota.SecurityGroup
	if securityGroupQuota.Used+2 >= securityGroupQuota.Limit {
		return false, fmt.Errorf("security group quota limit exceeded")
	}
	securityGroupRuleQuota := quotas.Quota.SecurityGroupRule
	if securityGroupRuleQuota.Used+10 >= securityGroupRuleQuota.Limit {
		return false, fmt.Errorf("security group rules quota limit exceeded")
	}
	portQuota := quotas.Quota.Port
	if portQuota.Used+3 >= portQuota.Limit {
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

		apiVersion := unstructuredObj.GetAPIVersion()
		kind := unstructuredObj.GetKind()
		if spec, ok := unstructuredObj.Object["spec"]; ok {
			specByte, _ := json.Marshal(spec)
			if strings.HasPrefix(apiVersion, "infrastructure.cluster.x-k8s.io") && kind == "OpenStackCluster" {
				infrastructureSpec := yamlmodel.InfrastructureSpec{}
				json.Unmarshal(specByte, &infrastructureSpec)
				infrastructureSpec.AllowAllInClusterTraffic = opt.InfrastructureKindSpecOption.AllowAllInClusterTraffic
				unstructuredObj.Object["spec"] = infrastructureSpec
			} else if strings.HasPrefix(apiVersion, "cluster.x-k8s.io") && kind == "Cluster" {
				clusterSpec := yamlmodel.ClusterSpec{}
				json.Unmarshal(specByte, &clusterSpec)
				if opt.ClusterKindSpecOption.CidrBlocks != nil {
					clusterSpec.ClusterNetwork.Pods.CidrBlocks = opt.ClusterKindSpecOption.CidrBlocks
				}
				unstructuredObj.Object["spec"] = clusterSpec
			}
		} else {
			unstructuredObjByte, _ := json.Marshal(unstructuredObj.Object)
			kindMap := map[string]interface{}{}
			if strings.HasPrefix(apiVersion, "storage.k8s.io") && kind == "StorageClass" {
				storageClass := yamlmodel.StorageClass{}
				json.Unmarshal(unstructuredObjByte, &storageClass)
				if opt.StorageClassKindOption.Parameters != nil {
					storageClass.Parameters = opt.StorageClassKindOption.Parameters
				}
				b, _ := json.Marshal(storageClass)
				json.Unmarshal(b, &kindMap)
				unstructuredObj.Object = kindMap
			} else if kind == "Secret" {
				secret := yamlmodel.Secret{}
				json.Unmarshal(unstructuredObjByte, &secret)
				if opt.SecretKindOption.Data != nil {
					secret.Data = opt.SecretKindOption.Data
				}
				b, _ := json.Marshal(secret)
				json.Unmarshal(b, &kindMap)
				unstructuredObj.Object = kindMap
			}
		}

		x, _ := yaml.Marshal(unstructuredObj.Object)
		yamlResult = yamlResult + "---\n" + string(x)
	}

	if err != io.EOF {
		return "", err
	}

	return yamlResult, nil
}
