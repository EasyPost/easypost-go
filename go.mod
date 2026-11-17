module github.com/EasyPost/easypost-go/v6

go 1.23

require (
	github.com/dnaeon/go-vcr v1.2.0
	github.com/google/go-querystring v1.1.0
	github.com/google/uuid v1.6.0
	// v1.12+ rejects the TestClient() suite helper because it returns a value
	github.com/stretchr/testify v1.11.1
	golang.org/x/text v0.11.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
