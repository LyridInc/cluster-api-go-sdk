package servercore

type Config struct {
	Username    string
	Password    string
	Domain      string
	ProjectName string

	ApiKey      string
	ApiUrl      string
	CloudApiUrl string
}

type AuthResponse struct {
	Token Token `json:"token"`
}

type Token struct {
	Methods   []string  `json:"methods"`
	User      User      `json:"user"`
	AuditIDs  []string  `json:"audit_ids"`
	ExpiresAt string    `json:"expires_at"`
	IssuedAt  string    `json:"issued_at"`
	Project   Project   `json:"project"`
	IsDomain  bool      `json:"is_domain"`
	Roles     []Role    `json:"roles"`
	Catalog   []Catalog `json:"catalog"`
}

type User struct {
	Domain            Domain  `json:"domain"`
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	PasswordExpiresAt *string `json:"password_expires_at"` // nullable
}

type Project struct {
	Domain Domain `json:"domain"`
	ID     string `json:"id"`
	Name   string `json:"name"`
}

type Domain struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Catalog struct {
	Endpoints []Endpoint `json:"endpoints"`
	ID        string     `json:"id"`
	Type      string     `json:"type"`
	Name      string     `json:"name"`
}

type Endpoint struct {
	ID        string `json:"id"`
	Interface string `json:"interface"`
	RegionID  string `json:"region_id"`
	URL       string `json:"url"`
	Region    string `json:"region"`
}
