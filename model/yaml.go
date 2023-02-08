package model

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/LyridInc/cluster-api-go-sdk/option"
	"gopkg.in/yaml.v2"
)

type (
	Openstack struct {
		Auth               map[string]interface{} `yaml:"auth" json:"auth"`
		Regions            []string               `yaml:"regions" json:"regions"`
		Interface          string                 `yaml:"interface" json:"interface"`
		IdentityApiVersion int                    `yaml:"identity_api_version" json:"identity_api_version"`
		RegionName         string                 `yaml:"region_name" json:"region_name"`
		Verify             bool                   `yaml:"verify" json:"verify"`
		LbMethod           string                 `yaml:"lb_method" json:"lb_method"`
		CreateMonitor      bool                   `yaml:"create_monitor" json:"create_monitor"`
		MonitorDelay       string                 `yaml:"monitor_delay" json:"monitor_delay"`
		MonitorMaxRetries  int                    `yaml:"monitor_max_retries" json:"monitor_max_retries"`
		MonitorTimeout     string                 `yaml:"monitor_timeout" json:"monitor_timeout"`
		CaCert             string                 `yaml:"cacert" json:"cacert"`
	}
	Clouds struct {
		Openstack Openstack `yaml:"openstack" json:"openstack"`
	}
	CloudsYaml struct {
		Clouds Clouds `yaml:"clouds" json:"clouds"`
	}
)

var yamlContent string

func (y *CloudsYaml) Parse(yamlByte []byte) error {
	err := yaml.Unmarshal(yamlByte, &y)
	if err != nil {
		return err
	}
	yamlContent = string(yamlByte)
	return nil
}

func (y *CloudsYaml) SetEnvironment(options option.OpenstackGenerateClusterOptions) {

	cloud := "openstack"
	cloudOs := y.Clouds.Openstack
	authOs := cloudOs.Auth
	openstackConf := "[Global]\n"

	// set env.rc variables
	if authUrl, ok := authOs["auth_url"]; ok {
		os.Setenv("CAPO_AUTH_URL", authUrl.(string))
		openstackConf = openstackConf + `auth-url="` + authUrl.(string) + "\"\n"
	}
	if username, ok := authOs["username"]; ok {
		os.Setenv("CAPO_USERNAME", username.(string))
		openstackConf = openstackConf + `username="` + username.(string) + "\"\n"
	}
	if password, ok := authOs["password"]; ok {
		os.Setenv("CAPO_PASSWORD", password.(string))
		openstackConf = openstackConf + `password="` + password.(string) + "\"\n"
	}
	if region, ok := authOs["region"]; ok {
		os.Setenv("CAPO_REGION", region.(string))
		openstackConf = openstackConf + `region="` + region.(string) + "\"\n"
	}
	if projectId, ok := authOs["project_id"]; ok {
		os.Setenv("CAPO_PROJECT_ID", projectId.(string))
		openstackConf = openstackConf + `tenant-id="` + projectId.(string) + "\"\n"
	}
	if projectName, ok := authOs["project_name"]; ok {
		os.Setenv("CAPO_PROJECT_NAME", projectName.(string))
		openstackConf = openstackConf + `tenant-name="` + projectName.(string) + "\"\n"
	}
	if domainName, ok := authOs["user_domain_name"]; ok {
		os.Setenv("CAPO_DOMAIN_NAME", domainName.(string))
		openstackConf = openstackConf + `domain-name="` + domainName.(string) + "\"\n"
	}

	caCertB64 := base64.StdEncoding.EncodeToString([]byte(cloudOs.CaCert + "\n"))
	os.Setenv("OPENSTACK_CLOUD_CACERT_B64", caCertB64)
	if cloudOs.CaCert != "" {
		openstackConf = openstackConf + `ca-file="/etc/certs/cacert"` + "\n"
	}

	if applicationCredentialName, ok := authOs["application_credential_name"]; ok {
		os.Setenv("CAPO_APPLICATION_CREDENTIAL_NAME", applicationCredentialName.(string))
		openstackConf = openstackConf + `application-credential-name="` + applicationCredentialName.(string) + "\"\n"
	}
	if applicationCredentialId, ok := authOs["application_credential_id"]; ok {
		os.Setenv("CAPO_APPLICATION_CREDENTIAL_ID", applicationCredentialId.(string))
		openstackConf = openstackConf + `application-credential-id="` + applicationCredentialId.(string) + "\"\n"
	}
	if applicationCredentialSecret, ok := authOs["application_credential_secret"]; ok {
		os.Setenv("CAPO_APPLICATION_CREDENTIAL_SECRET", applicationCredentialSecret.(string))
		openstackConf = openstackConf + `application-credential-secret="` + applicationCredentialSecret.(string) + "\"\n"
	}

	if cloudOs.LbMethod != "" ||
		cloudOs.CreateMonitor ||
		cloudOs.MonitorDelay != "" ||
		cloudOs.MonitorMaxRetries != 0 ||
		cloudOs.MonitorTimeout != "" {
		openstackConf = openstackConf + "\n[LoadBalancer]\n"
	}

	if cloudOs.LbMethod != "" {
		os.Setenv("CAPO_LB_METHOD", cloudOs.LbMethod)
		openstackConf = openstackConf + `lb-method="` + cloudOs.LbMethod + "\"\n"
	}
	if cloudOs.CreateMonitor {
		os.Setenv("CAPO_CREATE_MONITOR", fmt.Sprint(cloudOs.CreateMonitor))
		openstackConf = openstackConf + `create-monitor="` + fmt.Sprint(cloudOs.CreateMonitor) + "\"\n"
	}
	if cloudOs.MonitorDelay != "" {
		os.Setenv("CAPO_MONITOR_DELAY", cloudOs.MonitorDelay)
		openstackConf = openstackConf + `monitor-delay="` + cloudOs.MonitorDelay + "\"\n"
	}
	if cloudOs.MonitorMaxRetries != 0 {
		os.Setenv("CAPO_MONITOR_MAX_RETRIES", fmt.Sprint(cloudOs.MonitorMaxRetries))
		openstackConf = openstackConf + `monitor-max-retries="` + fmt.Sprint(cloudOs.MonitorMaxRetries) + "\"\n"
	}
	if cloudOs.MonitorTimeout != "" {
		os.Setenv("CAPO_MONITOR_TIMEOUT", cloudOs.MonitorTimeout)
		openstackConf = openstackConf + `monitor-timeout="` + cloudOs.MonitorTimeout + "\"\n"
	}

	if options.IgnoreVolumeAZ {
		openstackConf = openstackConf + "\n[BlockStorage]\n"
		openstackConf = openstackConf + "ignore-volume-az=true\n"
	}

	openstackConfB64 := base64.StdEncoding.EncodeToString([]byte(openstackConf))
	cloudYamlB64 := base64.StdEncoding.EncodeToString([]byte(yamlContent))

	os.Setenv("OPENSTACK_CLOUD", cloud)
	os.Setenv("OPENSTACK_CLOUD_PROVIDER_CONF", openstackConf)
	os.Setenv("OPENSTACK_CLOUD_PROVIDER_CONF_B64", openstackConfB64)
	os.Setenv("OPENSTACK_CLOUD_YAML_B64", cloudYamlB64)

	// set generate cluster options
	if options.ControlPlaneMachineFlavor != "" {
		os.Setenv("OPENSTACK_CONTROL_PLANE_MACHINE_FLAVOR", options.ControlPlaneMachineFlavor)
	}
	if options.NodeMachineFlavor != "" {
		os.Setenv("OPENSTACK_NODE_MACHINE_FLAVOR", options.NodeMachineFlavor)
	}
	if options.ExternalNetworkId != "" {
		os.Setenv("OPENSTACK_EXTERNAL_NETWORK_ID", options.ExternalNetworkId)
	}
	if options.ImageName != "" {
		os.Setenv("OPENSTACK_IMAGE_NAME", options.ImageName)
	}
	if options.SshKeyName != "" {
		os.Setenv("OPENSTACK_SSH_KEY_NAME", options.SshKeyName)
	}
	if options.DnsNameServers != "" {
		os.Setenv("OPENSTACK_DNS_NAMESERVERS", options.DnsNameServers)
	}
	if options.FailureDomain != "" {
		os.Setenv("OPENSTACK_FAILURE_DOMAIN", options.FailureDomain)
	}
}

func ReadYamlFromUrl(url string) (string, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
