package easypost

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

var apiBaseURL = &url.URL{
	Scheme: "https", Host: "api.easypost.com", Path: "/v2/",
}

var defaultUserAgent string
var defaultTimeout int

func init() {
	defaultUserAgent = fmt.Sprintf(
		"EasyPost/v2 GoClient/%s Go/%s OS/%s",
		Version, runtime.Version(), runtime.GOOS)
	defaultTimeout = 60000
}

// A Client provides an HTTP client for EasyPost API operations.
type Client struct {
	// BaseURL specifies the location of the API. It is used with
	// ResolveReference to create request URLs. (If 'Path' is specified, it
	// should end with a trailing slash.) If nil, the default will be used.
	BaseURL *url.URL
	// Client is an HTTP client used to make API requests. If nil,
	// http.DefaultClient will be used.
	Client *http.Client
	// APIKey is the user's API key. It is required.
	// Note: Treat your API Keys as passwords—keep them secret. API Keys give
	// full read/write access to your account, so they should not be included in
	// public repositories, emails, client side code, etc.
	APIKey string
	// UserAgent is a User-Agent to be sent with API HTTP requests. If empty,
	// a default will be used.
	UserAgent string
	// Timeout specifies the time limit (in milliseconds) for requests made by this Client. The
	// timeout includes connection time, any redirects, and reading the
	// response body.
	Timeout int
}

// New returns a new Client with the given API key.
func New(apiKey string) *Client {
	return &Client{APIKey: apiKey}
}

func (c *Client) baseURL() *url.URL {
	if c.BaseURL != nil {
		return c.BaseURL
	}
	return apiBaseURL
}

func (c *Client) userAgent() string {
	if c.UserAgent != "" {
		return c.UserAgent
	}
	return defaultUserAgent
}

func (c *Client) timeout() time.Duration {
	// return timeout duration in milliseconds
	timeout := c.Timeout
	if c.Timeout <= 0 {
		timeout = defaultTimeout
	}
	return time.Duration(timeout) * time.Millisecond
}

func (c *Client) client() *http.Client {
	client := c.Client
	if client == nil {
		client = http.DefaultClient
	}
	client.Timeout = c.timeout()
	return client
}

func (c *Client) convertOptsToURLValues(v interface{}) url.Values {
	values, _ := query.Values(v)
	return values
}

func (c *Client) setBody(req *http.Request, in interface{}) error {
	switch in := in.(type) {
	case url.Values:
		buf := []byte(in.Encode())
		req.Body = ioutil.NopCloser(bytes.NewReader(buf))
		req.GetBody = func() (io.ReadCloser, error) {
			return ioutil.NopCloser(bytes.NewReader(buf)), nil
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		// Setting Content-Length avoids chunked encoding, which the API
		// server doesn't currently support.
		req.ContentLength = int64(len(buf))
	case nil: // nop
	default:
		buf, err := json.Marshal(in)
		if err != nil {
			return err
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(buf))
		req.GetBody = func() (io.ReadCloser, error) {
			return ioutil.NopCloser(bytes.NewReader(buf)), nil
		}
		req.Header.Set("Content-Type", "application/json")
		// Setting Content-Length avoids chunked encoding, which the API
		// server doesn't currently support.
		req.ContentLength = int64(len(buf))
	}
	return nil
}

func (c *Client) do(ctx context.Context, method, path string, in, out interface{}) error {
	req := &http.Request{
		Method: method,
		URL:    c.baseURL().ResolveReference(&url.URL{Path: path}),
		Header: make(http.Header, 2),
	}
	req.Header.Set("User-Agent", c.userAgent())
	if err := c.setBody(req, in); err != nil {
		return err
	}
	req.SetBasicAuth(c.APIKey, "")
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	res, err := c.client().Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		if out != nil {
			return json.NewDecoder(res.Body).Decode(out)
		}
		return nil
	}

	buf, _ := ioutil.ReadAll(res.Body)
	apiErr := &APIError{Status: res.Status, StatusCode: res.StatusCode}
	if json.Unmarshal(buf, &apiErrorResponse{Error: apiErr}) != nil {
		apiErr.Message = string(buf)
	}
	return apiErr
}

func (c *Client) get(ctx context.Context, path string, out interface{}) error {
	return c.do(ctx, http.MethodGet, path, nil, out)
}

func (c *Client) post(ctx context.Context, path string, in, out interface{}) error {
	return c.do(ctx, http.MethodPost, path, in, out)
}

func (c *Client) put(ctx context.Context, path string, in, out interface{}) error {
	return c.do(ctx, http.MethodPut, path, in, out)
}

func (c *Client) del(ctx context.Context, path string) error {
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}
