package jumpserver

import (
	"fmt"
)

// UsersService handles communication with the user related
// methods of the GitHub API.
//
// GitHub API docs: https://developer.github.com/v3/users/
type PermsService struct {
	client *Client

	// Authentication type
	authType int

	// Basic auth username
	username string

	// Basic auth password
	password string
}

type RefreshResp struct {
	Msg bool
}

func (s *PermsService) RefreshCache() (*RefreshResp, error) {
	apiEndpoint := fmt.Sprintf("/api/v1/perms/asset-permissions/cache/refresh/")
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, err
	}

	var resp *RefreshResp
	_, err = s.client.Do(req, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
