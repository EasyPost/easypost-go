package easypost_test

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/EasyPost/easypost-go"
)

var (
	// TestClient is used for most tests.
	TestClient = new(easypost.Client)
	// ProdClient is used for tests that require a prod key (like the user API).
	ProdClient = new(easypost.Client)
)

func unique() string {
	t := time.Now().Format("20060102150405")
	r := strconv.FormatInt(int64(rand.Intn(100)), 10)
	return t + r
}

func TestMain(m *testing.M) {
	// Create a separate FlagSet just so we can print help specific to our
	// flags--don't need to dump help for all of Go's internal test.* flags.
	fs := flag.NewFlagSet("test", flag.ExitOnError)
	fs.StringVar(
		&TestClient.APIKey,
		"test-api-key",
		os.Getenv("EASYPOST_TEST_API_KEY"),
		"test key to use against EasyPost API",
	)
	fs.StringVar(
		&ProdClient.APIKey,
		"prod-api-key",
		os.Getenv("EASYPOST_PROD_API_KEY"),
		"production key to use against EasyPost API",
	)
	// Add flags to CommandLine (default) FlagSet:
	fs.VisitAll(func(f *flag.Flag) { flag.Var(f.Value, f.Name, f.Usage) })
	flag.Parse()
	if TestClient.APIKey == "" {
		fmt.Fprintln(os.Stderr, "must specify test API key")
		fs.PrintDefaults()
		os.Exit(1)
	}
	if ProdClient.APIKey == "" {
		fmt.Fprintln(
			os.Stderr,
			"no production API key given--skipping tests requiring a "+
				"production API key",
		)
	}
	os.Exit(m.Run())
}
