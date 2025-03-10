package easypost

import (
	"io"
	"net/http"
	"regexp"
	"strings"
)

type MockRequestMatchRule struct {
	Method          string
	UrlRegexPattern string
}

type MockRequestResponseInfo struct {
	StatusCode int
	Body       string
}

func (r *MockRequestResponseInfo) MockBody() io.ReadCloser {
	return io.NopCloser(strings.NewReader(r.Body))
}

type MockRequest struct {
	MatchRule    MockRequestMatchRule
	ResponseInfo MockRequestResponseInfo
}

func (r *MockRequestResponseInfo) AsResponse() *http.Response {
	return &http.Response{
		Status:        http.StatusText(r.StatusCode),
		StatusCode:    r.StatusCode,
		Body:          r.MockBody(),
		ContentLength: int64(len(r.Body)),
		Header:        nil,
		Request:       nil,
	}
}

func (c *Client) findMatchingMockRequest(req *http.Request) *http.Response {
	for _, mockReq := range c.MockRequests {
		url := req.URL.String()
		methodMatch := mockReq.MatchRule.Method == "" || mockReq.MatchRule.Method == req.Method
		urlMatch, _ := regexp.MatchString(mockReq.MatchRule.UrlRegexPattern, url)
		if methodMatch && urlMatch {
			return mockReq.ResponseInfo.AsResponse()
		}
	}
	return nil
}
