package servercore

type FloatingIPsResponse struct {
	FloatingIps []FloatingIP `json:"floatingips"`
}

type FloatingIPResponse struct {
	FloatingIP FloatingIP `json:"floatingip"`
}

type FloatingIP struct {
	FixedIPAddress    string       `json:"fixed_ip_address"`
	FloatingIPAddress string       `json:"floating_ip_address"`
	ID                string       `json:"id"`
	Loadbalancer      Loadbalancer `json:"loadbalancer"`
	PortID            string       `json:"port_id"`
	ProjectID         string       `json:"project_id"`
	Region            string       `json:"region"`
	Servers           []Server     `json:"servers"`
	Status            string       `json:"status"`
}

type CreateFloatingIPRequest struct {
	FloatingIPs []FloatingIPRequest `json:"floatingips"`
}

type FloatingIPRequest struct {
	IPs      string `json:"ips,omitempty"`
	Quantity uint   `json:"quantity"`
	Region   string `json:"region"`
}

type Loadbalancer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Server struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Updated string `json:"updated"`
}
