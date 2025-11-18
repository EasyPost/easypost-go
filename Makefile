PROJECT_NAME := easypost-go
GO_BIN := $(shell go env GOPATH)/bin
DIST_PATH := dist

## help - Display help about make targets for this Makefile
help:
	@cat Makefile | grep '^## ' --color=never | cut -c4- | sed -e "`printf 's/ - /\t- /;'`" | column -s "`printf '\t'`" -t

## build - Build the project
build:
	go build -o ../$(DIST_PATH)/$(PROJECT_NAME)

## clean - Clean the project
clean:
	rm -rf $(DIST_PATH)
	rm $(GO_BIN)/$(PROJECT_NAME)

## coverage - Get test coverage and open it in a browser
coverage:
	go clean -testcache
	go test -coverprofile=covprofile ./...
	@statement_cov=$$(go tool cover -func=covprofile | grep total: | awk '{print substr($$NF, 1, length($$NF)-1)}'); \
	if [ $$(echo "$$statement_cov < 78.0" | bc) -eq 1 ]; then \
		echo "Tests passed but statement coverage failed with coverage: $$statement_cov"; \
		exit 1; \
	fi
	go tool cover -html=covprofile

## init-examples-submodule - Initialize the examples submodule
init-examples-submodule:
	git submodule init
	git submodule update

## install - Install and vendor dependencies
install: | init-examples-submodule
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GO_BIN) v1.64.6
	curl -sSfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $(GO_BIN)
	go mod vendor

## golangci - Lint the project
golangci:
	$(GO_BIN)/golangci-lint run --config examples/style_guides/golang/.golangci.yml

## lint - Lint the project
lint: golangci scan

## lint-fix - Runs formatting tools against the project
lint-fix:
	$(GO_BIN)/golangci-lint run --config examples/style_guides/golang/.golangci.yml --fix

## release - Cuts a release for the project on GitHub (requires GitHub CLI)
# tag = The associated tag title of the release
# target = Target branch or full commit SHA
release:
	gh release create ${tag} --target ${target}

## scan - Run gosec to scan for security issues
scan:
	$(GO_BIN)/gosec -tests --exclude-dir=examples ./...

## test - Test the project
test:
	go clean -testcache && go test ./tests -v

## tidy - Tidies up the vendor directory
tidy:
	go mod tidy

## update-examples-submodule - Update the examples submodule
update-examples-submodule:
	git submodule init
	git submodule update --remote

.PHONY: help build clean coverage install golangci lint scan test tidy update-examples-submodule
