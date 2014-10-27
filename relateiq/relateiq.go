// Package relateiq is a client for the RelateIQ API
// (https://api.relateiq.com/).
package relateiq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

const (
	libraryVersion = "0.0.1"
	userAgent      = "go-relateiq/" + libraryVersion
)

// Credentials for access to the RelateIQ API.
type Credentials struct {
	APIKey    string
	APISecret string
}

// A Client communicates with the RelateIQ HTTP API.
type Client struct {
	// BaseURL is the base URL for all HTTP requests; by default,
	// "https://api.relateiq.com/v2/".
	BaseURL *url.URL

	// UserAgent is the HTTP User-Agent to send with all requests.
	UserAgent string

	cred Credentials // API credentials

	httpClient *http.Client // HTTP client to use when contacting API

	Accounts *AccountsService // Accounts service
}

// NewClient creates a new client for communicating with the RelateIQ
// HTTP API.
//
// Authentication using API credentials is required to access the
// RelateIQ API. You can obtain an API key and secret from your
// organization's integration settings screen.
//
//  c := NewClient(nil, Credentials{APIKey: "x", APISecret: "y"})
//  // ...
func NewClient(httpClient *http.Client, cred Credentials) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	c := &Client{
		BaseURL:    &url.URL{Scheme: "https", Host: "api.relateiq.com", Path: "/v2/"},
		UserAgent:  userAgent,
		cred:       cred,
		httpClient: httpClient,
	}
	c.Accounts = &AccountsService{c}
	return c
}

// NewRequest creates an API request. A relative URL can be provided
// in urlStr, in which case it is resolved relative to the BaseURL of
// the Client. Relative URL paths should always be specified without a
// preceding slash. If opt is specified, its encoding (using
// go-querystring) is used as the request URL's querystring. If body
// is specified, the value pointed to by body is JSON encoded and
// included as the request body.
func (c *Client) NewRequest(method, urlPath string, opt interface{}, body interface{}) (*http.Request, error) {
	u := c.BaseURL.ResolveReference(&url.URL{Path: urlPath})

	if opt != nil {
		qs, err := query.Values(opt)
		if err != nil {
			return nil, err
		}
		u.RawQuery = qs.Encode()
	}

	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.cred.APIKey, c.cred.APISecret)

	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

// Do sends an API request and returns the API response. The API
// response is decoded and stored in the value pointed to by v, or
// returned as an error if an API error has occurred.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		// even though there was an error, we still return the
		// response in case the caller wants to inspect it further
		return resp, err
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return nil, fmt.Errorf("reading response from %s %s: %s", req.Method, req.URL.RequestURI(), err)
		}
	}
	return resp, nil
}

// ListOptions specifies common options for endpoints that return
// lists.
type ListOptions struct {
	Start int `url:"_start,omitempty"`
	Limit int `url:"_limit,omitempty"`
}

// CheckResponse checks the API response for errors, and returns them
// if present. A response is considered an error if it has a status
// code outside the 200 range. API error responses are expected to
// have either no response body, or a JSON response body that maps to
// ErrorResponse. Any other response body will be silently ignored.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}

// An ErrorResponse reports errors caused by an API request.
type ErrorResponse struct {
	Response *http.Response `json:",omitempty"` // HTTP response that caused this error
	Message  string         // error message
}

// Error returns a string describing the API error response.
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message)
}

// HTTPStatusCode returns the HTTP status code of the API error
// response.
func (r *ErrorResponse) HTTPStatusCode() int { return r.Response.StatusCode }
