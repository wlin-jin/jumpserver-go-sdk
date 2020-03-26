package jumpserver

import (
	"fmt"
)

// UsersService handles communication with the user related
// methods of the GitHub API.
//
// GitHub API docs: https://developer.github.com/v3/users/
type AssetsService struct {
	client *Client

	// Authentication type
	authType int

	// Basic auth username
	username string

	// Basic auth password
	password string
}

type Asset struct {
	AdminUser    string `json:"admin_user"`
	Comment      string `json:"comment"`
	Connectivity struct {
		Datetime string `json:"datetime"`
		Status   int64  `json:"status"`
	} `json:"connectivity"`
	CPUCores     interface{}   `json:"cpu_cores"`
	CPUCount     interface{}   `json:"cpu_count"`
	CPUModel     interface{}   `json:"cpu_model"`
	CPUVcpus     interface{}   `json:"cpu_vcpus"`
	CreatedBy    string        `json:"created_by"`
	DateCreated  string        `json:"date_created"`
	DiskInfo     interface{}   `json:"disk_info"`
	DiskTotal    interface{}   `json:"disk_total"`
	Domain       interface{}   `json:"domain"`
	HardwareInfo string        `json:"hardware_info"`
	Hostname     string        `json:"hostname"`
	HostnameRaw  interface{}   `json:"hostname_raw"`
	ID           string        `json:"id"`
	IP           string        `json:"ip"`
	IsActive     bool          `json:"is_active"`
	Labels       []interface{} `json:"labels"`
	Memory       interface{}   `json:"memory"`
	Model        interface{}   `json:"model"`
	Nodes        []string      `json:"nodes"`
	Number       interface{}   `json:"number"`
	OrgID        string        `json:"org_id"`
	OrgName      string        `json:"org_name"`
	Os           interface{}   `json:"os"`
	OsArch       interface{}   `json:"os_arch"`
	OsVersion    interface{}   `json:"os_version"`
	Platform     string        `json:"platform"`
	Protocols    []string      `json:"protocols"`
	PublicIP     string        `json:"public_ip"`
	Sn           interface{}   `json:"sn"`
	Vendor       interface{}   `json:"vendor"`
}

func (s *AssetsService) GetList() ([]*Asset, *Response, error) {
	apiEndpoint := "/api/assets/v1/assets/"
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	var assets []*Asset
	resp, err := s.client.Do(req, &assets)
	if err != nil {
		return nil, resp, err
	}
	return assets, resp, nil
}

func (s *AssetsService) Search(assetName string) ([]*Asset, *Response, error) {
	apiEndpoint := fmt.Sprintf("/api/assets/v1/assets/?search=%s", assetName)
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	var assets []*Asset
	resp, err := s.client.Do(req, &assets)
	if err != nil {
		return nil, resp, err
	}
	return assets, resp, nil
}

func (s *AssetsService) Get(AssetId string) (*Asset, *Response, error) {
	apiEndpoint := fmt.Sprintf("/api/assets/v1/assets/%s/", AssetId)
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	var asset *Asset
	resp, err := s.client.Do(req, &asset)
	if err != nil {
		return nil, resp, err
	}
	return asset, resp, nil
}

func (s *AssetsService) Delete(ip string) error {
	apiEndpoint := fmt.Sprintf("/api/v1/assets/assets/?ip=%s", ip)
	req, err := s.client.NewRequest("DELETE", apiEndpoint, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		return err
	}
	return nil
}

// asset body
type AssetBody struct {
	AdminUser string   `json:"admin_user"`
	Comment   string   `json:"comment"`
	Domain    string   `json:"domain"`
	Hostname  string   `json:"hostname"`
	IP        string   `json:"ip"`
	IsActive  bool     `json:"is_active"`
	Labels    []string `json:"labels"`
	Nodes     []string `json:"nodes"`
	Number    string   `json:"number"`
	Platform  string   `json:"platform"`
	Port      int64    `json:"port"`
	Protocol  string   `json:"protocol"`
	Protocols []string `json:"protocols"`
	PublicIP  string   `json:"public_ip"`
}

func (s *AssetsService) Create(body *AssetBody) (*Asset, error) {
	apiEndpoint := fmt.Sprintf("/api/v1/assets/assets/")
	req, err := s.client.NewRequest("POST", apiEndpoint, body)
	if err != nil {
		return nil, err
	}

	var asset *Asset
	resp, err := s.client.Do(req, &asset)
	fmt.Println(resp)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return asset, nil
}
