package easypost_test

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/EasyPost/easypost-go"
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
}

func (c *ClientTests) SetupTest() {
	pathComponents := append(
		[]string{"testdata"}, strings.Split(c.T().Name(), "/")...,
	)
	r, err := recorder.New(filepath.Join(pathComponents...))
	c.Require().NoError(err)
	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
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
