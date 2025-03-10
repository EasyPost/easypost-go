package easypost_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/EasyPost/easypost-go/v4"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/suite"
)

var (
	// TestAPIKey is used for most tests.
	TestAPIKey string
	// ProdAPIKey is used for tests that require a prod key (like the user API).
	ProdAPIKey string
	// PartnerAPIKey is used for tests that require a partner key (like the white label API).
	PartnerAPIKey string
	// ReferralAPIKey is used for tests that require a referral key (like the white label API) (internal, access via ReferralAPIKey()).
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

// contains returns true if the string is in the list.
func contains(elements []string, target string) bool {
	for _, element := range elements {
		if target == element {
			return true
		}
	}
	return false
}

// applyCensorsToJsonList applies the censors to a JSON list.
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

// applyCensorsToJsonDictionary applies the censors to a JSON dictionary.
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

// censorJsonData censors data in cassette files before they persist to disk
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

// SetupTest runs all the logic required for a test to run
func (c *ClientTests) SetupTest() {
	pathComponents := append(
		[]string{"cassettes/"}, strings.Split(c.T().Name(), "/")[1],
	)
	r, err := recorder.New(filepath.Join(pathComponents...))
	c.Require().NoError(err)

	c.checkExpiredCassette()

	// Filter sensitive data from cassettes
	// Replace value has to be type specific to its corresponding struct
	redactedString := "REDACTED"
	redactedStringMap := map[string]string{}
	responseBodyElementsToCensor := map[string]interface{}{
		"client_ip":        redactedString,
		"credentials":      redactedStringMap,
		"email":            redactedString,
		"fields":           redactedStringMap,
		"key":              redactedString,
		"phone_number":     redactedString,
		"phone":            redactedString,
		"test_credentials": redactedStringMap,
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
		r.Body = io.NopCloser(&b)
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
		redactedStringList := []string{"REDACTED"}
		if i.Request.Headers["Authorization"] != nil {
			i.Request.Headers["Authorization"] = redactedStringList
		}
		if i.Request.Headers["User-Agent"] != nil {
			i.Request.Headers["User-Agent"] = redactedStringList
		}
		if i.Request.Headers["X-Client-User-Agent"] != nil {
			i.Request.Headers["X-Client-User-Agent"] = redactedStringList
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

// checkExpiredCassette checks for an expired cassette and warns if it is too old and must be re-recorded
func (c *ClientTests) checkExpiredCassette() {
	fullCassettePath := "cassettes/" + strings.Split(c.T().Name(), "/")[1] + ".yaml"
	const expirationDays = 365
	expirationHours := expirationDays * 24 * time.Hour

	if _, err := os.Stat(fullCassettePath); err == nil {
		cassette, _ := os.Stat(fullCassettePath)
		cassetteTimestamp := cassette.ModTime()
		expirationTimestamp := cassetteTimestamp.Add(expirationHours)
		currentTimestamp := time.Now()

		if currentTimestamp.After(expirationTimestamp) {
			c.T().Logf("%s is older than %d days and has expired. Please re-record the cassette", fullCassettePath, expirationDays)
		}
	}
}

// TearDownTest runs all the logic required to teardown a test
func (c *ClientTests) TearDownTest() {
	_ = c.recorder.Stop()
}

// TestClient sets up the test mode client object to be used in the test
func (c *ClientTests) TestClient() *easypost.Client {
	return &easypost.Client{
		APIKey: TestAPIKey,
		Client: &http.Client{Transport: c.recorder},
	}
}

// ProdClient sets up the prod mode client object to be used in the test
func (c *ClientTests) ProdClient() *easypost.Client {
	return &easypost.Client{
		APIKey: ProdAPIKey,
		Client: &http.Client{Transport: c.recorder},
	}
}

// PartnerClient sets up the partner user client object to be used in the test
func (c *ClientTests) PartnerClient() *easypost.Client {
	if len(PartnerAPIKey) == 0 {
		PartnerAPIKey = "123"
	}
	return &easypost.Client{
		APIKey: PartnerAPIKey,
		Client: &http.Client{Transport: c.recorder},
	}
}

// ReferralClient sets up the referral customer client object to be used in the test
func (c *ClientTests) ReferralClient() *easypost.Client {
	if len(ReferralAPIKey) == 0 {
		ReferralAPIKey = "123"
	}
	return &easypost.Client{
		APIKey: ReferralAPIKey,
		Client: &http.Client{Transport: c.recorder},
	}
}

// ReferralAPIKey returns the referral api key or a fallback if the environment variable is not set
func (c *ClientTests) ReferralAPIKey() string {
	if len(ReferralAPIKey) == 0 {
		ReferralAPIKey = "123"
	}
	return ReferralAPIKey
}

// MockClient sets up the mock client object to be used in the test
func (c *ClientTests) MockClient(requests []easypost.MockRequest) *easypost.Client {
	return &easypost.Client{
		APIKey:       "cannot_be_blank",
		Client:       &http.Client{Transport: c.recorder},
		MockRequests: requests,
	}
}

// TestClient runs the entire test suite
func TestClient(t *testing.T) {
	suite.Run(t, new(ClientTests))
}

// TestMain is the entrypoint when running tests via the CLI
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
		os.Getenv("REFERRAL_CUSTOMER_PROD_API_KEY"),
		"production key for referral account to use against EasyPost API",
	)

	// Add flags to CommandLine (default) FlagSet:
	fs.VisitAll(func(f *flag.Flag) { flag.Var(f.Value, f.Name, f.Usage) })
	flag.Parse()
	testSuite := m.Run()

	// Fail test suite below desired coverage
	// `testSuite = 0` means it passed, CoverMode will be non-empty if run with -cover
	coverageMinPercentage := 0.75
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
