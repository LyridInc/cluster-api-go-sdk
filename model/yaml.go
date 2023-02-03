package model

import (
	"os"

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
	}
	Clouds struct {
		Openstack Openstack `yaml:"openstack" json:"openstack"`
	}
	CloudsYaml struct {
		Clouds Clouds `yaml:"clouds" json:"clouds"`
	}

	OpenstackGenerateClusterOptions struct {
		ControlPlaneMachineFlavor string
		NodeMachineFlavor         string
		ExternalNetworkId         string
		ImageName                 string
		SshKeyName                string
		DnsNameServers            string
		FailureDomain             string
	}
)

func (y *CloudsYaml) Parse(yamlByte []byte) error {
	err := yaml.Unmarshal(yamlByte, &y)
	if err != nil {
		return err
	}
	return nil
}

func (y *CloudsYaml) SetEnvironment(options OpenstackGenerateClusterOptions) {

	// set env.rc variables
	if authUrl, ok := y.Clouds.Openstack.Auth["auth_url"]; ok {
		os.Setenv("CAPO_AUTH_URL", authUrl.(string))
	}
	if username, ok := y.Clouds.Openstack.Auth["username"]; ok {
		os.Setenv("CAPO_USERNAME", username.(string))
	}
	if password, ok := y.Clouds.Openstack.Auth["password"]; ok {
		os.Setenv("CAPO_PASSWORD", password.(string))
	}
	if region, ok := y.Clouds.Openstack.Auth["region"]; ok {
		os.Setenv("CAPO_REGION", region.(string))
	}
	if projectId, ok := y.Clouds.Openstack.Auth["project_id"]; ok {
		os.Setenv("CAPO_PROJECT_ID", projectId.(string))
	}
	if projectName, ok := y.Clouds.Openstack.Auth["project_name"]; ok {
		os.Setenv("CAPO_PROJECT_NAME", projectName.(string))
	}
	if domainName, ok := y.Clouds.Openstack.Auth["domain_name"]; ok {
		os.Setenv("CAPO_DOMAIN_NAME", domainName.(string))
	}
	if applicationCredentialName, ok := y.Clouds.Openstack.Auth["application_credential_name"]; ok {
		os.Setenv("CAPO_APPLICATION_CREDENTIAL_NAME", applicationCredentialName.(string))
	}

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
