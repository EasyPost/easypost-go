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

## install - Install and vendor dependencies
install: | update-examples-submodule
	brew install golangci-lint || exit 0
	go mod vendor
	go build -o $(PROJECT_PATH)

## golangci - Lint the project
golangci:
	golangci-lint run --config examples/style_guides/golang/.golangci.yml

## lint - Lint the project
lint: golangci scan

## release - Cuts a release for the project on GitHub (requires GitHub CLI)
# tag = The associated tag title of the release
# target = Target branch or full commit SHA
release:
	gh release create ${tag} --target ${target}

## scan - Run gosec to scan for security issues
scan:
	gosec -tests --exclude-dir=examples ./...

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
