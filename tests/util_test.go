package easypost_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/EasyPost/easypost-go/v2"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/suite"
)

var (
	// TestAPIKey is used for most tests.
	TestAPIKey string
	// ProdAPIKey is used for tests that require a prod key (like the user API).
	ProdAPIKey string
)

type ErrorRoundTripper struct{}

func (ErrorRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("no cassette found for request %s", req.URL)
}

type ClientTests struct {
	suite.Suite
	recorder *recorder.Recorder
	fixture  *Fixture
}

func (c *ClientTests) SetupTest() {
	pathComponents := append(
		[]string{"cassettes/"}, strings.Split(c.T().Name(), "/")[1],
	)
	r, err := recorder.New(filepath.Join(pathComponents...))
	c.Require().NoError(err)

	// Strictly match the URL, method, and body of the requests
	r.SetMatcher(func(r *http.Request, i cassette.Request) bool {
		if r.Body == nil {
			return cassette.DefaultMatcher(r, i)
		}
		var b bytes.Buffer
		if _, err := b.ReadFrom(r.Body); err != nil {
			return false
		}
		r.Body = ioutil.NopCloser(&b)
		return cassette.DefaultMatcher(r, i) && (b.String() == "" || b.String() == i.Body)
	})

	// Filter sensitive data from cassettes
	// Replace value has to be type specific to its corresponding struct
	responseBodyScrubbers := map[string]interface{}{
		"api_keys":         []string{},
		"children":         []string{},
		"client_ip":        "REDACTED",
		"test_credentials": map[string]string{},
		"credentials":      map[string]string{},
		"email":            "REDACTED",
		"keys":             []*easypost.APIKey{},
		"key":              "REDACTED",
		"phone_number":     "REDACTED",
		"phone":            "REDACTED",
	}
	r.AddSaveFilter(func(i *cassette.Interaction) error {
		// Filter headers
		if i.Request.Headers["Authorization"] != nil {
			i.Request.Headers["Authorization"] = []string{"REDACTED"}
		}

		if i.Request.Headers["User-Agent"] != nil {
			i.Request.Headers["User-Agent"] = []string{"REDACTED"}
		}

		if i.Request.Headers["X-Client-User-Agent"] != nil {
			i.Request.Headers["X-Client-User-Agent"] = []string{"REDACTED"}
		}

		// Filter response bodies
		var responseBodyBytes []byte
		var responseBodyString string

		for scrubber, replaceValue := range responseBodyScrubbers {
			if strings.Contains(i.Response.Body, scrubber) {
				var responseBody = i.Response.Body
				var arrayObjects []map[string]interface{}

				if json.Unmarshal([]byte(responseBody), &arrayObjects) == nil {
					err := json.Unmarshal([]byte(responseBody), &arrayObjects)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					for index, object := range arrayObjects {
						for key := range object {
							if key == scrubber {
								arrayObjects[index][key] = replaceValue
							}
						}
					}

					responseBodyBytes, _ = json.Marshal(arrayObjects)
				} else {
					var responseMap map[string]interface{}
					err := json.Unmarshal([]byte(responseBody), &responseMap)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					if _, keyExists := responseMap[scrubber]; keyExists {
						responseMap[scrubber] = replaceValue
					}
					responseBodyBytes, _ = json.Marshal(responseMap)
				}

				responseBodyString = string(responseBodyBytes)
				i.Response.Body = responseBodyString
			}
		}

		return nil
	})

	if TestAPIKey == "" {
		r.SetTransport(ErrorRoundTripper{})
	}

	c.recorder = r
}

func (c *ClientTests) TearDownTest() {
	_ = c.recorder.Stop()
}

func (c *ClientTests) TestClient() *easypost.Client {
	return &easypost.Client{
		APIKey: TestAPIKey,
		Client: &http.Client{Transport: c.recorder},
	}
}

func (c *ClientTests) ProdClient() *easypost.Client {
	return &easypost.Client{
		APIKey: ProdAPIKey,
		Client: &http.Client{Transport: c.recorder},
	}
}

func TestClient(t *testing.T) {
	suite.Run(t, new(ClientTests))
}

func TestMain(m *testing.M) {
	// Create a separate FlagSet just so we can print help specific to our
	// flags--don't need to dump help for all of Go's internal test.* flags.
	fs := flag.NewFlagSet("test", flag.ExitOnError)
	fs.StringVar(
		&TestAPIKey,
		"test-api-key",
		os.Getenv("EASYPOST_TEST_API_KEY"),
		"test key to use against EasyPost API",
	)
	fs.StringVar(
		&ProdAPIKey,
		"prod-api-key",
		os.Getenv("EASYPOST_PROD_API_KEY"),
		"production key to use against EasyPost API",
	)
	// Add flags to CommandLine (default) FlagSet:
	fs.VisitAll(func(f *flag.Flag) { flag.Var(f.Value, f.Name, f.Usage) })
	flag.Parse()
	os.Exit(m.Run())
}
