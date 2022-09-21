package easypost_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/EasyPost/easypost-go/v2"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	// TestAPIKey is used for most tests.
	TestAPIKey string
	// ProdAPIKey is used for tests that require a prod key (like the user API).
	ProdAPIKey string
	// PartnerAPIKey is used for tests that require a partner key (like the white label API).
	PartnerAPIKey string
	// ReferralAPIKey is used for tests that require a referral key (like the white label API).
	ReferralAPIKey string
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

func contains(elements []string, target string) bool {
	for _, element := range elements {
		if target == element {
			return true
		}
	}
	return false
}

func (c *ClientTests) applyCensorsToJsonList(list []interface{}, elementsToCensor map[string]interface{}) []interface{} {
	if len(list) == 0 {
		// short circuit and return the list if it is empty
		return list
	}

	for index, value := range list {
		if value == nil {
			// don't need to worry about censoring nil values
			continue
		}

		// nolint:gosimple
		switch value.(type) {
		case map[string]interface{}:
			// value is a dictionary
			list[index] = c.applyCensorsToJsonDictionary(value.(map[string]interface{}), elementsToCensor)
		case []map[string]interface{}:
			// value is a list of dictionaries
			list[index] = c.applyCensorsToJsonList(value.([]interface{}), elementsToCensor)
		default:
			// value is a single value or a normal list, nothing to censor
		}
	}

	return list
}

func (c *ClientTests) applyCensorsToJsonDictionary(dictionary map[string]interface{}, elementsToCensor map[string]interface{}) map[string]interface{} {
	if len(dictionary) == 0 {
		// short circuit and return the dictionary if it is empty
		return dictionary
	}

	for key, value := range dictionary {
		if value == nil {
			// don't need to worry about censoring nil values
			continue
		}

		var censorKeys []string
		for k := range elementsToCensor {
			censorKeys = append(censorKeys, k)
		}

		if contains(censorKeys, key) {
			// element should be censored
			// replace value with corresponding censor value
			dictionary[key] = elementsToCensor[key]
		} else {
			// element doesn't need to be censored
			// nolint:gosimple
			switch value.(type) {
			case map[string]interface{}:
				// value is a dictionary
				dictionary[key] = c.applyCensorsToJsonDictionary(value.(map[string]interface{}), elementsToCensor)
			case []interface{}:
				// value is a list
				dictionary[key] = c.applyCensorsToJsonList(value.([]interface{}), elementsToCensor)
			default:
				// value is a single value, nothing to censor
			}
		}
	}

	return dictionary
}

func (c *ClientTests) censorJsonData(data string, elementsToCensor map[string]interface{}) string {
	var jsonMap map[string]interface{}
	mapErr := json.Unmarshal([]byte(data), &jsonMap)
	if mapErr != nil {
		// data is not a JSON dictionary
		var jsonList []interface{}
		listErr := json.Unmarshal([]byte(data), &jsonList)

		if listErr != nil {
			// data is not a JSON list either
			// short circuit and return the data
			return data
		}

		// data is a JSON list
		censoredList := c.applyCensorsToJsonList(jsonList, elementsToCensor)
		censoredListBytes, _ := json.Marshal(censoredList)
		return string(censoredListBytes)
	}

	// data is a JSON dictionary
	censoredDictionary := c.applyCensorsToJsonDictionary(jsonMap, elementsToCensor)
	censoredDictionaryBytes, _ := json.Marshal(censoredDictionary)
	return string(censoredDictionaryBytes)
}

func (c *ClientTests) SetupTest() {
	pathComponents := append(
		[]string{"cassettes/"}, strings.Split(c.T().Name(), "/")[1],
	)
	r, err := recorder.New(filepath.Join(pathComponents...))
	c.Require().NoError(err)

	// Filter sensitive data from cassettes
	// Replace value has to be type specific to its corresponding struct
	responseBodyElementsToCensor := map[string]interface{}{
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
		bString := b.String()
		if bString == "" && i.Body == "" {
			// short circuit and return true if the body is empty as it should be
			return true
		}
		// run the request body through the same censors before comparing to the recording
		bStringCensored := c.censorJsonData(bString, responseBodyElementsToCensor)
		return cassette.DefaultMatcher(r, i) && (bStringCensored == i.Body)
	})

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

		// Censor request data
		var requestBody = i.Request.Body
		var censoredRequestBody = c.censorJsonData(requestBody, responseBodyElementsToCensor)
		i.Request.Body = censoredRequestBody

		// Censor response data
		var responseBody = i.Response.Body
		var censoredResponseBody = c.censorJsonData(responseBody, responseBodyElementsToCensor)
		i.Response.Body = censoredResponseBody

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

func (c *ClientTests) PartnerClient() *easypost.Client {
	return &easypost.Client{
		APIKey: PartnerAPIKey,
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
	fs.StringVar(
		&PartnerAPIKey,
		"partner-api-key",
		os.Getenv("PARTNER_USER_PROD_API_KEY"),
		"production key for partner account to use against EasyPost API",
	)
	fs.StringVar(
		&ReferralAPIKey,
		"referral-api-key",
		os.Getenv("REFERRAL_USER_PROD_API_KEY"),
		"production key for referral account to use against EasyPost API",
	)

	// Add flags to CommandLine (default) FlagSet:
	fs.VisitAll(func(f *flag.Flag) { flag.Var(f.Value, f.Name, f.Usage) })
	flag.Parse()
	testSuite := m.Run()

	// Fail test suite below desired coverage
	// `testSuite = 0` means it passed, CoverMode will be non empty if run with -cover
	coverageMinPercentage := 0.75 // TODO: This number for whatever reason is about ~5% lower than what the CLI reports which is the real value
	if testSuite == 0 && testing.CoverMode() != "" {
		coverage := testing.Coverage()
		if coverage < coverageMinPercentage {
			fmt.Println("Tests passed but coverage failed with coverage: ", coverage)
			testSuite = -1
		}
	}

	// Exit once tests complete
	os.Exit(testSuite)
}
