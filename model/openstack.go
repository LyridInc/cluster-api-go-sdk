package model

import "time"

type (
	QuotasResponse struct {
		Quota Quota `json:"quota"`
	}

	Quota struct {
		Network           QuotaValues `json:"network"`
		Subnet            QuotaValues `json:"subnet"`
		SubnetPool        QuotaValues `json:"subnetpool"`
		Port              QuotaValues `json:"port"`
		Router            QuotaValues `json:"router"`
		FloatingIp        QuotaValues `json:"floatingip"`
		RbacPolicy        QuotaValues `json:"rbac_policy"`
		SecurityGroup     QuotaValues `json:"security_group"`
		SecurityGroupRule QuotaValues `json:"security_group_rule"`
		Trunk             QuotaValues `json:"trunk"`
	}

	QuotaValues struct {
		Limit    int `json:"limit"`
		Used     int `json:"used"`
		Reserved int `json:"reserved"`
	}

	LoadBalancerResponse struct {
		LoadBalancer LoadBalancer `json:"loadbalancer"`
	}

	LoadBalancer struct {
		ID                 string        `json:"id"`
		AdminStateUp       bool          `json:"admin_state_up"`
		Description        string        `json:"description"`
		ProjectID          string        `json:"project_id"`
		ProvisioningStatus string        `json:"provisioning_status"`
		FlavorID           string        `json:"flavor_id"`
		VipSubnetID        string        `json:"vip_subnet_id"`
		VipAddress         string        `json:"vip_address"`
		VipNetworkID       string        `json:"vip_network_id"`
		VipPortID          string        `json:"vip_port_id"`
		AdditionalVips     []interface{} `json:"additional_vips"`
		Provider           string        `json:"provider"`
		CreatedAt          *time.Time    `json:"created_at"`
		UpdatedAt          *time.Time    `json:"updated_at"`
		OperatingStatus    string        `json:"operating_status"`
		Name               string        `json:"name"`
		VipQosPolicyID     string        `json:"vip_qos_policy_id"`
		AvailabilityZone   string        `json:"availability_zone"`
		Tags               []interface{} `json:"tags"`
	}
)
