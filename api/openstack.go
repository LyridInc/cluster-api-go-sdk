package api

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"reflect"
	"strings"

	"github.com/LyridInc/cluster-api-go-sdk/model"
	"github.com/LyridInc/cluster-api-go-sdk/option"
	"github.com/LyridInc/cluster-api-go-sdk/utils"
	yamlmodel "github.com/LyridInc/cluster-api-go-sdk/yaml"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	yamlserializer "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
)

type (
	OpenstackClient struct {
		MagnumEndpoint       string
		NetworkEndpoint      string
		LoadBalancerEndpoint string
		AuthEndpoint         string
		ImageEndpoint        string
		ComputeEndpoint      string
		AuthToken            string
		ProjectId            string
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

	OpenstackProject struct {
		Name   *string                 `json:"name,omitempty"`
		Domain *map[string]interface{} `json:"domain,omitempty"`
	}

	OpenstackScope struct {
		Project *OpenstackProject `json:"project,omitempty"`
	}

	OpenstackAuth struct {
		Identity OpenstackIdentity `json:"identity"`
		Scope    *OpenstackScope   `json:"scope,omitempty"`
	}
)

func (c *OpenstackClient) doHttpRequest(request *http.Request) ([]byte, error) {
	logRequest, _ := httputil.DumpRequest(request, true)
	log.Printf("Request: %q", logRequest)

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)
	// log.Printf("Response: %q", string(b))
	return b, err
}

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

	authMap := map[string]interface{}{}
	authByte, _ := json.Marshal(auth)
	json.Unmarshal(authByte, &authMap)

	if identity, ok := authMap["identity"]; ok {
		if len(auth.Identity.Methods) > 0 {
			identityMap := identity.(map[string]interface{})
			switch auth.Identity.Methods[0] {
			case "password":
				delete(identityMap, "application_credential")

			case "application_credential":
				delete(identityMap, "password")
			}
			authMap["identity"] = identityMap
		}
	}

	url := c.AuthEndpoint + "/v3/auth/tokens"
	b, _ := json.Marshal(map[string]interface{}{
		"auth": authMap,
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

	body, _ := c.doHttpRequest(request)

	jsonResponse := map[string]interface{}{}
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

	body, _ := c.doHttpRequest(request)

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
	if floatingIpQuota.Limit != -1 && floatingIpQuota.Limit != 0 && floatingIpQuota.Used >= floatingIpQuota.Limit {
		return false, fmt.Errorf("floating ip quota limit exceeded")
	}
	networkQuota := quotas.Quota.Network
	if networkQuota.Limit != -1 && networkQuota.Limit != 0 && networkQuota.Used >= networkQuota.Limit {
		return false, fmt.Errorf("network quota limit exceeded")
	}
	routerQuota := quotas.Quota.Router
	if routerQuota.Limit != -1 && routerQuota.Limit != 0 && routerQuota.Used >= routerQuota.Limit {
		return false, fmt.Errorf("router quota limit exceeded")
	}
	securityGroupQuota := quotas.Quota.SecurityGroup
	if securityGroupQuota.Limit != -1 && securityGroupQuota.Limit != 0 && securityGroupQuota.Used+2 >= securityGroupQuota.Limit {
		return false, fmt.Errorf("security group quota limit exceeded")
	}
	securityGroupRuleQuota := quotas.Quota.SecurityGroupRule
	if securityGroupRuleQuota.Limit != -1 && securityGroupRuleQuota.Limit != 0 && securityGroupRuleQuota.Used+10 >= securityGroupRuleQuota.Limit {
		return false, fmt.Errorf("security group rules quota limit exceeded")
	}
	portQuota := quotas.Quota.Port
	if portQuota.Limit != -1 && portQuota.Limit != 0 && portQuota.Used+3 >= portQuota.Limit {
		return false, fmt.Errorf("port quota limit exceeded")
	}

	return true, nil
}

func (c *OpenstackClient) GetLoadBalancer(loadBalancerId string) (*model.LoadBalancerResponse, error) {
	url := c.LoadBalancerEndpoint + "/v2/lbaas/loadbalancers/" + loadBalancerId
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Token", c.AuthToken)

	body, _ := c.doHttpRequest(request)

	loadBalancer := model.LoadBalancerResponse{}
	json.Unmarshal(body, &loadBalancer)

	if loadBalancer.LoadBalancer.ID == "" {
		return nil, nil
	}

	return &loadBalancer, nil
}

func (c *OpenstackClient) GetImageList(filter map[string]string) (*model.ImageListResponse, error) {
	url := c.ImageEndpoint + "/v2/images"
	if filter != nil {
		queryString := utils.MapToQueryString(filter)
		url = fmt.Sprintf("%s?%s", url, queryString)
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Token", c.AuthToken)

	body, _ := c.doHttpRequest(request)

	imageListResponse := model.ImageListResponse{}
	json.Unmarshal(body, &imageListResponse)

	return &imageListResponse, nil
}

func (c *OpenstackClient) GetImageByID(id string) (*model.Image, error) {
	url := c.ImageEndpoint + "/v2/images/" + id

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Token", c.AuthToken)

	body, _ := c.doHttpRequest(request)

	imageResponse := model.Image{}
	json.Unmarshal(body, &imageResponse)

	return &imageResponse, nil
}

func (c *OpenstackClient) GetNetworkList(filter map[string]string) (*model.NetworkListResponse, error) {
	url := c.NetworkEndpoint + "/v2.0/networks"
	if filter != nil {
		queryString := utils.MapToQueryString(filter)
		url = fmt.Sprintf("%s?%s", url, queryString)
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Token", c.AuthToken)

	body, _ := c.doHttpRequest(request)

	networkListResponse := model.NetworkListResponse{}
	json.Unmarshal(body, &networkListResponse)

	return &networkListResponse, nil
}

func (c *OpenstackClient) GetNetworkByID(id string) (*model.NetworkResponse, error) {
	url := c.NetworkEndpoint + "/v2.0/networks/" + id

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Token", c.AuthToken)

	body, _ := c.doHttpRequest(request)

	networkResponse := model.NetworkResponse{}
	json.Unmarshal(body, &networkResponse)

	return &networkResponse, nil
}

func (c *OpenstackClient) GetSubnetList(filter map[string]string) (*model.SubnetListResponse, error) {
	url := c.NetworkEndpoint + "/v2.0/subnets"
	if filter != nil {
		queryString := utils.MapToQueryString(filter)
		url = fmt.Sprintf("%s?%s", url, queryString)
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Token", c.AuthToken)

	body, _ := c.doHttpRequest(request)

	subnetListResponse := model.SubnetListResponse{}
	json.Unmarshal(body, &subnetListResponse)

	return &subnetListResponse, nil
}

func (c *OpenstackClient) GetSubnetByID(id string) (*model.SubnetResponse, error) {
	url := c.NetworkEndpoint + "/v2.0/subnets/" + id

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Token", c.AuthToken)

	body, _ := c.doHttpRequest(request)

	subnetResponse := model.SubnetResponse{}
	json.Unmarshal(body, &subnetResponse)

	return &subnetResponse, nil
}

func (c *OpenstackClient) GetFlavorList(filter map[string]string) (*model.FlavorListResponse, error) {
	url := c.ComputeEndpoint + "/flavors/detail"
	if filter != nil {
		queryString := utils.MapToQueryString(filter)
		url = fmt.Sprintf("%s?%s", url, queryString)
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Token", c.AuthToken)

	body, _ := c.doHttpRequest(request)

	flavorListResponse := model.FlavorListResponse{}
	json.Unmarshal(body, &flavorListResponse)

	return &flavorListResponse, nil
}

func (c *OpenstackClient) GetFlavorByID(id string) (*model.FlavorResponse, error) {
	url := c.ComputeEndpoint + "/flavors/" + id

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Token", c.AuthToken)

	body, _ := c.doHttpRequest(request)

	flavorResponse := model.FlavorResponse{}
	json.Unmarshal(body, &flavorResponse)

	return &flavorResponse, nil
}

func (c *OpenstackClient) GetKeypairList(filter map[string]string) (*model.KeypairListResponse, error) {
	url := c.ComputeEndpoint + "/os-keypairs"
	if filter != nil {
		queryString := utils.MapToQueryString(filter)
		url = fmt.Sprintf("%s?%s", url, queryString)
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Token", c.AuthToken)

	body, _ := c.doHttpRequest(request)

	keypairListResponse := model.KeypairListResponse{}
	json.Unmarshal(body, &keypairListResponse)

	return &keypairListResponse, nil
}

func (c *OpenstackClient) CreateKeypair(keypairName string) (*model.KeypairResponse, error) {
	keypairRequest := map[string]interface{}{}
	keypairRequest["name"] = keypairName
	requestBody, _ := json.Marshal(map[string]interface{}{
		"keypair": keypairRequest,
	})

	url := c.ComputeEndpoint + "/os-keypairs"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Token", c.AuthToken)

	body, _ := c.doHttpRequest(request)

	keypairResponse := model.KeypairResponse{}
	json.Unmarshal(body, &keypairResponse)

	return &keypairResponse, nil
}

func (c *OpenstackClient) DeleteKeypair(keypairName string) (*model.KeypairResponse, error) {
	url := c.ComputeEndpoint + "/os-keypairs/" + keypairName
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Token", c.AuthToken)

	body, _ := c.doHttpRequest(request)

	keypairResponse := model.KeypairResponse{}
	json.Unmarshal(body, &keypairResponse)

	return &keypairResponse, nil
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

// magnum client - start

func (c *OpenstackClient) MagnumListClusters() ([]byte, error) {
	url := c.MagnumEndpoint + "/clusters"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Token", c.AuthToken)

	return c.doHttpRequest(request)
}

func (c *OpenstackClient) MagnumCreateCluster(args model.MagnumCreateClusterRequest) ([]byte, error) {
	url := c.MagnumEndpoint + "/clusters"
	b, _ := json.Marshal(args)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Token", c.AuthToken)

	return c.doHttpRequest(request)
}

func (c *OpenstackClient) MagnumListClusterTemplates() ([]byte, error) {
	url := c.MagnumEndpoint + "/clustertemplates"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Token", c.AuthToken)

	return c.doHttpRequest(request)
}

func (c *OpenstackClient) MagnumCreateClusterTemplate(args model.MagnumCreateClusterTemplateRequest) ([]byte, error) {
	url := c.MagnumEndpoint + "/clustertemplates"
	b, _ := json.Marshal(args)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Auth-Token", c.AuthToken)

	return c.doHttpRequest(request)
}

func (c *OpenstackClient) MagnumGenerateKubeconfig(clusterID string) ([]byte, error) {
	// sign CA and CSR
	csrBytes, csrPrivateBytes, err := c.CreateClusterCertificateSigningRequest(clusterID)
	if err != nil {
		return nil, err
	}

	requestBody := map[string]string{
		"cluster_uuid": clusterID,
		"csr":          string(csrBytes),
	}
	b, _ := json.Marshal(requestBody)
	payload := bytes.NewBufferString(string(b))

	url := c.MagnumEndpoint + "/certificates"
	request, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Auth-Token", c.AuthToken)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", "None")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	certResponse, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var jsonRes map[string]interface{}
	jwtBytes := []byte(certResponse)
	if err := json.Unmarshal(jwtBytes, &jsonRes); err != nil {
		return nil, err
	}

	log.Println(jsonRes)
	encodedCsr := base64.StdEncoding.EncodeToString(csrBytes)
	encodedPrivateCsr := base64.StdEncoding.EncodeToString(csrPrivateBytes)
	encodedPem := base64.StdEncoding.EncodeToString([]byte(jsonRes["pem"].(string)))

	config := model.KubeconfigConfig{
		ApiVersion: "v1",
		Kind:       "Config",
		Clusters: []model.KubeconfigCluster{{
			CertificateAuthorityData: encodedPem,
			Server:                   "https://103.176.45.244:6443",
		}},
		Users: []model.KubeconfigUser{{
			ClientCertificateData: encodedCsr,
			ClientKeyData:         encodedPrivateCsr,
			// Token:                 "jwt_token",
		}},
		Contexts: []model.KubeconfigContext{{
			Cluster: "lyra-dev",
			User:    "admin",
		}},
		CurrentContext: "default",
	}

	// TODO: put it as YAML
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(yamlData))

	return certResponse, nil
}

// magnum client - end

func (c *OpenstackClient) CreateClusterCertificateSigningRequest(clusterIdentification string) ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	keyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privatePemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyBytes,
	}
	privatePemBytes := pem.EncodeToMemory(privatePemBlock)

	template := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: clusterIdentification,
		},
		DNSNames: []string{fmt.Sprintf("%s.lyr.id", clusterIdentification)},
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, privateKey)
	if err != nil {
		return nil, nil, err
	}

	pemBlock := &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	}
	pemBytes := pem.EncodeToMemory(pemBlock)

	return pemBytes, privatePemBytes, nil
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
			if len(opt.InfrastructureKindSpecOption.ManagedSubnets) > 0 {
				infrastructureSpec.ManagedSubnets = opt.InfrastructureKindSpecOption.ManagedSubnets
			}
			// infrastructureSpec.AllowAllInClusterTraffic = opt.InfrastructureKindSpecOption.AllowAllInClusterTraffic
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
