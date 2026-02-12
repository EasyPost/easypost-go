package easypost

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
)

var apiBaseURL = &url.URL{
	Scheme: "https", Host: "api.easypost.com", Path: "/v2/",
}

var defaultUserAgent string
var defaultTimeout int

func init() {
	// We skip grabbing the OS version (for now) as there is not a reliable way to do so across OS's
	defaultUserAgent = fmt.Sprintf(
		"EasyPost/v2 GoClient/%s Go/%s OS/%s OSVersion/%s OSArch/%s",
		Version, runtime.Version(), runtime.GOOS, "NA", runtime.GOARCH)
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
	// MockRequests is a list of requests that will be mocked by the client.
	MockRequests []MockRequest
	// Hooks is a collection of HookEventSubscriber instances for various hooks available in the client
	Hooks Hooks
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

func (c *Client) setParameters(req *http.Request, params interface{}) error {
	switch req.Method {
	case http.MethodGet, http.MethodDelete:
		// Convert interface into query parameters and set as request URL query
		var values url.Values
		if urlValues, ok := params.(url.Values); ok {
			// If params is already url.Values, use it directly
			values = urlValues
		} else {
			// Otherwise, use the query package to convert the struct
			values, _ = query.Values(params)
		}
		req.URL.RawQuery = values.Encode()
		return nil

	case http.MethodPost, http.MethodPut, http.MethodPatch:
		// Convert interface into JSON and set as request body
		buf, err := json.Marshal(params)
		if err != nil {
			return err
		}
		req.Body = io.NopCloser(bytes.NewReader(buf))
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(buf)), nil
		}
		req.Header.Set("Content-Type", "application/json")
		// Setting Content-Length avoids chunked encoding, which the API
		// server doesn't currently support.
		req.ContentLength = int64(len(buf))

		return nil
	default:
		return fmt.Errorf("unsupported method: %s", req.Method)
	}
}

func (c *Client) do(ctx context.Context, method, path string, params interface{}, out interface{}) error {
	if c.APIKey == "" {
		return newMissingPropertyError("APIKey")
	}

	req := &http.Request{
		Method: method,
		URL:    c.baseURL().ResolveReference(&url.URL{Path: path}),
		Header: make(http.Header, 2),
	}

	req.Header.Set("User-Agent", c.userAgent())
	if err := c.setParameters(req, params); err != nil {
		return err
	}

	req.SetBasicAuth(c.APIKey, "")
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	var res *http.Response
	var err error

	// prepare and execute request hook(s)
	requestId := uuid.New()
	requestTimestamp := time.Now()
	requestEvent := &RequestHookEvent{
		Method:           req.Method,
		Url:              req.URL,
		RequestBody:      req.Body,
		Headers:          req.Header,
		RequestTimestamp: requestTimestamp,
		Id:               requestId,
	}

	// loop over each request hook and execute it
	for _, hook := range c.Hooks.RequestHookEventSubscriptions {
		hook.Execute(ctx, *requestEvent)
	}

	if len(c.MockRequests) > 0 {
		// If there are mock requests set, this client will ONLY make mock requests
		res = c.findMatchingMockRequest(req)
		if res == nil {
			return errors.New("no matching mock request found")
		}
	} else {
		// Otherwise, make a real request
		res, err = c.client().Do(req)
	}

	if err != nil {
		// prepare and execute response hook(s) for failed requests
		responseEvent := &ResponseHookEvent{
			HttpStatus:        0,
			Method:            req.Method,
			Url:               req.URL,
			ResponseBody:      nil,
			Headers:           nil,
			RequestTimestamp:  requestTimestamp,
			ResponseTimestamp: time.Now(),
			Id:                requestId,
		}
		// loop over each response hook and execute it
		for _, hook := range c.Hooks.ResponseHookEventSubscriptions {
			hook.Execute(ctx, *responseEvent)
		}

		return err
	}

	// prepare and execute response hook(s) for successful requests
	responseEvent := &ResponseHookEvent{
		HttpStatus:        res.StatusCode,
		Method:            req.Method,
		Url:               req.URL,
		ResponseBody:      res.Body,
		Headers:           res.Header,
		RequestTimestamp:  requestTimestamp,
		ResponseTimestamp: time.Now(),
		Id:                requestId,
	}

	// loop over each response hook and execute it
	for _, hook := range c.Hooks.ResponseHookEventSubscriptions {
		hook.Execute(ctx, *responseEvent)
	}

	defer func() { _ = res.Body.Close() }()

	// status code is 2xx, no error occurred
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		if out != nil {
			return json.NewDecoder(res.Body).Decode(out)
		}
		return nil
	}

	// status code is not 2xx, an error occurred
	apiErr := BuildErrorFromResponse(res)

	return apiErr
}

// MakeAPICall makes an API call to the EasyPost API.
//
// This public, generic interface is useful for making arbitrary API calls to the EasyPost API that
// are not yet supported by the client library's services. When possible, the service for your use case
// should be used instead as it provides a more convenient and higher-level interface depending on the endpoint.
func (c *Client) MakeAPICall(method, endpoint string, params map[string]interface{}) (out map[string]interface{}, err error) {
	// For GET/DELETE requests, convert map to url.Values since query package doesn't handle maps
	var paramsToUse interface{} = params
	if method == http.MethodGet || method == http.MethodDelete {
		values := url.Values{}
		for key, value := range params {
			values.Set(key, fmt.Sprintf("%v", value))
		}
		paramsToUse = values
	}

	err = c.do(context.TODO(), method, endpoint, paramsToUse, &out)
	return
}
