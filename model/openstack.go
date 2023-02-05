package model

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
)
