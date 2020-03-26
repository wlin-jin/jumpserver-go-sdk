package jumpserver

import (
	"bytes"
	"encoding/json"

	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	defaultBaseURL = "https://api.github.com/"
	uploadBaseURL  = "https://uploads.github.com/"
	userAgent      = "go-github"

	headerRateLimit     = "X-RateLimit-Limit"
	headerRateRemaining = "X-RateLimit-Remaining"
	headerRateReset     = "X-RateLimit-Reset"
	headerOTP           = "X-GitHub-OTP"
)

// httpClient defines an interface for an http.Client implementation so that alternative
// http Clients can be passed in for making requests
type httpClient interface {
	Do(request *http.Request) (response *http.Response, err error)
}

// A Client manages communication with the JIRA API.
type Client struct {
	// HTTP client used to communicate with the API.
	client httpClient

	// Base URL for API requests.
	baseURL *url.URL

	// Services used for talking to different parts of the JIRA API.
	Users  *UsersService
	Assets *AssetsService
	Perms  *PermsService
}

func NewClient(httpClient httpClient, baseURL string) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	// ensure the baseURL contains a trailing slash so that all paths are preserved in later calls
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		client:  httpClient,
		baseURL: parsedBaseURL,
	}

	// services
	c.Users = &UsersService{client: c}
	c.Assets = &AssetsService{client: c}
	c.Perms = &PermsService{client: c}
	return c, nil
}

// NewRequest creates an API request.
// A relative URL can be provided in urlStr, in which case it is resolved relative to the baseURL of the Client.
// If specified, the value pointed to by body is JSON encoded and included as the request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	// Relative URLs should be specified without a preceding slash since baseURL will have the trailing slash
	rel.Path = strings.TrimLeft(rel.Path, "/")

	u := c.baseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// Do sends an API request and returns the API response.
// The API response is JSON decoded and stored in the value pointed to by v, or returned as an error if an API error has occurred.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	httpResp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	err = CheckResponse(httpResp)
	if err != nil {
		// Even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return newResponse(httpResp, nil), err
	}

	//bodyBytes, err := ioutil.ReadAll(httpResp.Body)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//bodyString := string(bodyBytes)
	//fmt.Print(bodyString)

	if v != nil {
		// Open a NewDecoder and defer closing the reader only if there is a provided interface to decode to
		defer httpResp.Body.Close()
		err = json.NewDecoder(httpResp.Body).Decode(v)
	}

	resp := newResponse(httpResp, v)

	return resp, err
}

// CheckResponse checks the API response for errors, and returns them if present.
// A response is considered an error if it has a status code outside the 200 range.
// The caller is responsible to analyze the response body.
// The body can contain JSON (if the error is intended) or xml (sometimes JIRA just failes).
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	err := fmt.Errorf("Request failed. Please analyze the request body for more details. Status code: %d", r.StatusCode)
	return err
}

// GetBaseURL will return you the Base URL.
// This is the same URL as in the NewClient constructor
func (c *Client) GetBaseURL() url.URL {
	return *c.baseURL
}

// Response represents JIRA API response. It wraps http.Response returned from
// API and provides information about paging.
type Response struct {
	*http.Response

	StartAt    int
	MaxResults int
	Total      int

	NextPage  int
	PrevPage  int
	FirstPage int
	LastPage  int
}

func newResponse(r *http.Response, v interface{}) *Response {
	resp := &Response{Response: r}
	resp.populatePageValues()
	return resp
}

// populatePageValues parses the HTTP Link response headers and populates the
// various pagination link values in the Response.
func (r *Response) populatePageValues() {
	if links, ok := r.Response.Header["Link"]; ok && len(links) > 0 {
		for _, link := range strings.Split(links[0], ",") {
			segments := strings.Split(strings.TrimSpace(link), ";")

			// link must at least have href and rel
			if len(segments) < 2 {
				continue
			}

			// ensure href is properly formatted
			if !strings.HasPrefix(segments[0], "<") || !strings.HasSuffix(segments[0], ">") {
				continue
			}

			// try to pull out page parameter
			url, err := url.Parse(segments[0][1 : len(segments[0])-1])
			if err != nil {
				continue
			}
			page := url.Query().Get("page")
			if page == "" {
				continue
			}

			for _, segment := range segments[1:] {
				switch strings.TrimSpace(segment) {
				case `rel="next"`:
					r.NextPage, _ = strconv.Atoi(page)
				case `rel="prev"`:
					r.PrevPage, _ = strconv.Atoi(page)
				case `rel="first"`:
					r.FirstPage, _ = strconv.Atoi(page)
				case `rel="last"`:
					r.LastPage, _ = strconv.Atoi(page)
				}

			}
		}
	}
}

// BasicAuthTransport is an http.RoundTripper that authenticates all requests
// using HTTP Basic Authentication with the provided username and password. It
// additionally supports users who have two-factor authentication enabled on
// their GitHub account.
type BasicAuthTransport struct {
	Username string // GitHub username
	Password string // GitHub password
	OTP      string // one-time password for users with two-factor auth enabled

	// Transport is the underlying HTTP transport to use when making requests.
	// It will default to http.DefaultTransport if nil.
	Transport http.RoundTripper
}

// RoundTrip implements the RoundTripper interface.
func (t *BasicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// To set extra headers, we must make a copy of the Request so
	// that we don't modify the Request we were given. This is required by the
	// specification of http.RoundTripper.
	//
	// Since we are going to modify only req.Header here, we only need a deep copy
	// of req.Header.
	req2 := new(http.Request)
	*req2 = *req
	req2.Header = make(http.Header, len(req.Header))
	for k, s := range req.Header {
		req2.Header[k] = append([]string(nil), s...)
	}

	if t.OTP != "" {
		req2.Header.Set(headerOTP, t.OTP)
	}
	return t.transport().RoundTrip(req2)
}

// Client returns an *http.Client that makes requests that are authenticated
// using HTTP Basic Authentication.
func (t *BasicAuthTransport) Client() *http.Client {
	return &http.Client{Transport: t}
}

func (t *BasicAuthTransport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}

// BasicAuthTransport is an http.RoundTripper that authenticates all requests
// using HTTP Basic Authentication with the provided username and password. It
// additionally supports users who have two-factor authentication enabled on
// their GitHub account.
type TokenAuthTransport struct {
	Username string // JumpServer username
	Password string // JumpServer password
	AuthURL  string // JumpServer auth url

	Token string

	OTP string // one-time password for users with two-factor auth enabled

	// Transport is the underlying HTTP transport to use when making requests.
	// It will default to http.DefaultTransport if nil.
	Transport http.RoundTripper
}

// RoundTrip implements the RoundTripper interface.
func (t *TokenAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// To set extra headers, we must make a copy of the Request so
	// that we don't modify the Request we were given. This is required by the
	// specification of http.RoundTripper.
	//
	// Since we are going to modify only req.Header here, we only need a deep copy
	// of req.Header.
	req2 := new(http.Request)
	*req2 = *req
	req2.Header = make(http.Header, len(req.Header))
	for k, s := range req.Header {
		req2.Header[k] = append([]string(nil), s...)
	}

	if t.Token == "" {
		err := t.setTokenAuth()
		if err != nil {
			return nil, errors.Wrap(err, "token auth: set token auth failed")
		}
	}

	err := t.signRequest(req2)
	if err != nil {
		return nil, err
	}

	if t.OTP != "" {
		req2.Header.Set(headerOTP, t.OTP)
	}
	return t.transport().RoundTrip(req2)
}
func (t *TokenAuthTransport) signRequest(req *http.Request) error {
	if t.Token == "" {
		return errors.New("token auth: token auth is empty")
	}

	token := fmt.Sprintf("Bearer %s", t.Token)
	req.Header.Set("Authorization", token)

	return nil
}

func (t *TokenAuthTransport) buildAuthRequest() (*http.Request, error) {
	body := struct {
		Username   string `json:"username"`
		Password   string `json:"password"`
		PublicKey  string `json:"public_key"`
		RemoteAddr string `json:"remote_addr"`
		LoginType  string `json:"login_type"`
	}{
		Username:   t.Username,
		Password:   t.Password,
		PublicKey:  "",
		RemoteAddr: "",
		LoginType:  "T",
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(body)

	req, err := http.NewRequest("POST", t.AuthURL, b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return req, nil

}

func (t *TokenAuthTransport) setTokenAuth() error {
	req, err := t.buildAuthRequest()
	if err != nil {
		return err
	}

	var authClient = &http.Client{
		Timeout: time.Second * 60,
	}
	resp, err := authClient.Do(req)
	if err != nil {
		return err
	}

	authInfo := new(AuthenticateInfo)
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(authInfo)
	if err != nil {
		return err
	}

	t.Token = authInfo.Token
	return nil
}

// Client returns an *http.Client that makes requests that are authenticated
// using HTTP Basic Authentication.
func (t *TokenAuthTransport) Client() *http.Client {
	return &http.Client{Transport: t}
}

func (t *TokenAuthTransport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}
