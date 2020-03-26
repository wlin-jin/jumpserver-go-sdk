package jumpserver

import (
	"fmt"
)

// UsersService handles communication with the user related
// methods of the GitHub API.
//
// GitHub API docs: https://developer.github.com/v3/users/
type UsersService struct {
	client *Client

	// Authentication type
	authType int

	// Basic auth username
	username string

	// Basic auth password
	password string
}

type AuthenticateInfo struct {
	Token string `json:"token,omitempty"`
	User  *User  `json:"user"`
}

// User represents a GitHub user.
type User struct {
	Id                      string    `json:"id,omitempty"`
	Name                    string    `json:"name,omitempty"`
	Username                string    `json:"username,omitempty"`
	Email                   string    `json:"email,omitempty"`
	Groups                  *[]string `json:"groups,omitempty"`
	GroupsDisplay           string    `json:"groups_display,omitempty"`
	Role                    string    `json:"role,omitempty"`
	RoleDisplay             string    `json:"role_display,omitempty"`
	AvatarUrl               string    `json:"avatar_url,omitempty"`
	Wechat                  string    `json:"wechat,omitempty"`
	Phone                   string    `json:"phone,omitempty"`
	OtpLevel                *int      `json:"otp_level,omitempty"`
	Comment                 string    `json:"comment,omitempty"`
	Source                  string    `json:"source,omitempty"`
	SourceDisplay           string    `json:"source_display,omitempty"`
	IsValid                 *bool     `json:"is_valid,omitempty"`
	IsExpired               *bool     `json:"is_expired,omitempty"`
	IsActive                *bool     `json:"is_active,omitempty"`
	CreatedBy               string    `json:"created_by,omitempty"`
	IsFirstLogin            *bool     `json:"is_first_login,omitempty"`
	DatePasswordLastUpdated string    `json:"date_password_last_updated,omitempty"`
	DateExpired             string    `json:"date_expired,omitempty"`

	// API URLs
	URL               string `json:"url,omitempty"`
	EventsURL         string `json:"events_url,omitempty"`
	FollowingURL      string `json:"following_url,omitempty"`
	FollowersURL      string `json:"followers_url,omitempty"`
	GistsURL          string `json:"gists_url,omitempty"`
	OrganizationsURL  string `json:"organizations_url,omitempty"`
	ReceivedEventsURL string `json:"received_events_url,omitempty"`
	ReposURL          string `json:"repos_url,omitempty"`
	StarredURL        string `json:"starred_url,omitempty"`
	SubscriptionsURL  string `json:"subscriptions_url,omitempty"`

	// Permissions identifies the permissions that a user has on a given
	// repository. This is only populated when calling Repositories.ListCollaborators.
	Permissions *map[string]bool `json:"permissions,omitempty"`
}

func (s *UsersService) GetList() ([]*User, *Response, error) {
	apiEndpoint := "/api/users/v1/users/"
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	var userList []*User
	resp, err := s.client.Do(req, &userList)
	if err != nil {
		return nil, resp, err
	}
	return userList, resp, nil
}

func (s *UsersService) Search(username string) ([]*User, error) {
	apiEndpoint := fmt.Sprintf("/api/v1/users/users/?username=%s", username)
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, err
	}

	var userList []*User
	_, err = s.client.Do(req, &userList)
	if err != nil {
		return nil, err
	}
	return userList, nil
}

type UserGroup struct {
	Comment     string `json:"comment"`
	CreatedBy   string `json:"created_by"`
	DateCreated string `json:"date_created"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	OrgID       string `json:"org_id"`
	OrgName     string `json:"org_name"`
	UsersAmount int64  `json:"users_amount"`
}

func (s *UsersService) Groups(name string) ([]*UserGroup, error) {
	apiEndpoint := fmt.Sprintf("/api/v1/users/groups/?name=test")
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, err
	}

	groups := make([]*UserGroup, 0)
	_, err = s.client.Do(req, &groups)
	if err != nil {
		return nil, err
	}
	return groups, nil
}
