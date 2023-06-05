PROJECT_NAME := easypost-go
PROJECT_PATH := $$(go env GOPATH)/bin/$(PROJECT_NAME)
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
	rm $(PROJECT_PATH)

## coverage - Get test coverage and open it in a browser
coverage: 
	go clean -testcache && go test ./tests -v -coverprofile=covprofile -coverpkg=./... && go tool cover -html=covprofile

# TODO: Change branch to master once examples are merged
## install-style - Download style guides
install-style:
	curl -LJs https://raw.githubusercontent.com/EasyPost/examples/style_guides/.golangci_go_cl.yml -o .golangci.yml

## install - Install and vendor dependencies
install: | install-style
	brew install golangci-lint || exit 0
	git submodule init
	git submodule update --remote
	go mod vendor
	go build -o $(PROJECT_PATH)

## lint - Lint the project
lint:
	golangci-lint run -e SA1019

## release - Cuts a release for the project on GitHub (requires GitHub CLI)
# tag = The associated tag title of the release
release:
	gh release create ${tag}

## gosec - Run gosec to scan for security issues
scan:
	gosec -tests --exclude-dir=examples ./...

## test - Test the project
test:
	go clean -testcache && go test ./tests -v

## tidy - Tidies up the vendor directory
tidy:
	go mod tidy

.PHONY: help build clean coverage install install-style lint scan test tidy
