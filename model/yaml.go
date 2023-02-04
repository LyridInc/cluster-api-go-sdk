package model

import (
	"encoding/base64"
	"fmt"
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
	openstackConf := `[Global]\n`

	// set env.rc variables
	if authUrl, ok := y.Clouds.Openstack.Auth["auth_url"]; ok {
		os.Setenv("CAPO_AUTH_URL", authUrl.(string))
		openstackConf = openstackConf + `auth-url="` + authUrl.(string) + "\"\n"
	}
	if username, ok := y.Clouds.Openstack.Auth["username"]; ok {
		os.Setenv("CAPO_USERNAME", username.(string))
		openstackConf = openstackConf + `username="` + username.(string) + "\"\n"
	}
	if password, ok := y.Clouds.Openstack.Auth["password"]; ok {
		os.Setenv("CAPO_PASSWORD", password.(string))
		openstackConf = openstackConf + `password="` + password.(string) + "\"\n"
	}
	if region, ok := y.Clouds.Openstack.Auth["region"]; ok {
		os.Setenv("CAPO_REGION", region.(string))
		openstackConf = openstackConf + `region="` + region.(string) + "\"\n"
	}
	if projectId, ok := y.Clouds.Openstack.Auth["project_id"]; ok {
		os.Setenv("CAPO_PROJECT_ID", projectId.(string))
		openstackConf = openstackConf + `tenant-id="` + projectId.(string) + "\"\n"
	}
	if projectName, ok := y.Clouds.Openstack.Auth["project_name"]; ok {
		os.Setenv("CAPO_PROJECT_NAME", projectName.(string))
		openstackConf = openstackConf + `tenant-name="` + projectName.(string) + "\"\n"
	}
	if domainName, ok := y.Clouds.Openstack.Auth["user_domain_name"]; ok {
		os.Setenv("CAPO_DOMAIN_NAME", domainName.(string))
		openstackConf = openstackConf + `domain-name="` + domainName.(string) + "\"\n"
	}

	caCertB64 := base64.StdEncoding.EncodeToString([]byte(y.Clouds.Openstack.CaCert + "\n"))
	os.Setenv("OPENSTACK_CLOUD_CACERT_B64", caCertB64)
	if y.Clouds.Openstack.CaCert != "" {
		openstackConf = openstackConf + `ca-file="/etc/certs/cacert"` + "\n"
	}

	if applicationCredentialName, ok := y.Clouds.Openstack.Auth["application_credential_name"]; ok {
		os.Setenv("CAPO_APPLICATION_CREDENTIAL_NAME", applicationCredentialName.(string))
		openstackConf = openstackConf + `application-credential-name="` + applicationCredentialName.(string) + "\"\n"
	}
	if applicationCredentialId, ok := y.Clouds.Openstack.Auth["application_credential_id"]; ok {
		os.Setenv("CAPO_APPLICATION_CREDENTIAL_ID", applicationCredentialId.(string))
		openstackConf = openstackConf + `application-credential-id="` + applicationCredentialId.(string) + "\"\n"
	}
	if applicationCredentialSecret, ok := y.Clouds.Openstack.Auth["application_credential_secret"]; ok {
		os.Setenv("CAPO_APPLICATION_CREDENTIAL_SECRET", applicationCredentialSecret.(string))
		openstackConf = openstackConf + `application-credential-secret="` + applicationCredentialSecret.(string) + "\"\n"
	}
	if y.Clouds.Openstack.LbMethod != "" {
		os.Setenv("CAPO_LB_METHOD", y.Clouds.Openstack.LbMethod)
		openstackConf = openstackConf + `lb-method="` + y.Clouds.Openstack.LbMethod + "\"\n"
	}
	if y.Clouds.Openstack.CreateMonitor {
		os.Setenv("CAPO_CREATE_MONITOR", fmt.Sprint(y.Clouds.Openstack.CreateMonitor))
		openstackConf = openstackConf + `create-monitor="` + fmt.Sprint(y.Clouds.Openstack.CreateMonitor) + "\"\n"
	}
	if y.Clouds.Openstack.MonitorDelay != "" {
		os.Setenv("CAPO_MONITOR_DELAY", y.Clouds.Openstack.MonitorDelay)
		openstackConf = openstackConf + `monitor-delay="` + y.Clouds.Openstack.MonitorDelay + "\"\n"
	}
	if y.Clouds.Openstack.MonitorMaxRetries != 0 {
		os.Setenv("CAPO_MONITOR_MAX_RETRIES", fmt.Sprint(y.Clouds.Openstack.MonitorMaxRetries))
		openstackConf = openstackConf + `monitor-max-retries="` + fmt.Sprint(y.Clouds.Openstack.MonitorMaxRetries) + "\"\n"
	}
	if y.Clouds.Openstack.MonitorTimeout != "" {
		os.Setenv("CAPO_MONITOR_TIMEOUT", y.Clouds.Openstack.MonitorTimeout)
		openstackConf = openstackConf + `monitor-timeout="` + y.Clouds.Openstack.MonitorTimeout + "\"\n"
	}

	openstackConfB64 := base64.StdEncoding.EncodeToString([]byte(openstackConf))
	cloudYamlB64 := base64.StdEncoding.EncodeToString([]byte(yamlContent))

	os.Setenv("OPENSTACK_CLOUD", cloud)
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
