# Build the project
build:
	go build -o ../dist/easypost-go

# Clean the project
clean:
	rm -rf dist
	rm $(go env GOPATH)/bin/easypost-go

# Get test coverage and open it in a browser
coverage:
    go clean -testcache
    go test -coverprofile=covprofile ./...
    bash -c 'statement_cov=$(go tool cover -func=covprofile | grep total: | awk "{print substr(\$NF, 1, length(\$NF)-1)}"); if [ $(echo "$statement_cov < 78.0" | bc) -eq 1 ]; then echo "Tests passed but statement coverage failed with coverage: $statement_cov"; exit 1; fi'
    go tool cover -html=covprofile

# Initialize the examples submodule
init-examples-submodule:
	git submodule init
	git submodule update

# Install and vendor dependencies
install: init-examples-submodule
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.64.6
	curl -sSfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
	go mod vendor

# Lint the project
golangci:
	$(go env GOPATH)/bin/golangci-lint run --config examples/style_guides/golang/.golangci.yml

# Lint the project
lint: golangci scan

# Runs formatting tools against the project
lint-fix:
	$(go env GOPATH)/bin/golangci-lint run --config examples/style_guides/golang/.golangci.yml --fix

# Cuts a release for the project on GitHub (requires GitHub CLI)
# tag = The associated tag title of the release
# target = Target branch or full commit SHA
release tag target:
	gh release create {{tag}} --target {{target}}

# Run gosec to scan for security issues
scan:
	$(go env GOPATH)/bin/gosec -tests --exclude-dir=examples ./...

# Test the project
test:
	go clean -testcache
	go test ./... -v

# Tidies up the vendor directory
tidy:
	go mod tidy

# Update the examples submodule
update-examples-submodule:
	git submodule init
	git submodule update --remote
