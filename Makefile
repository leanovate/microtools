PACKAGES=$(shell go list ./...)
VERSION = $(shell date -u +.%Y%m%d.%H%M%S)

all: export GOPATH=${PWD}/../../../..
all: format
	@echo "--> Running go build"
	@go build ./...

format: export GOPATH=${PWD}/../../../..
format:
	@echo "--> Running go fmt"
	@go fmt ./...

test: export GOPATH=${PWD}/../../../..
test:
	@echo "--> Running tests"
	@go test -v ./rest/... ./routing/... ./logging/...

coverage:
	@echo "--> Running tests with coverage"
	@echo "" > coverage.txt
	for pkg in $(shell go list ./rest/... ./routing/... ./logging/...); do \
          (go test -coverprofile=.pkg.coverage -covermode=atomic -v $$pkg && \
          cat .pkg.coverage >> coverage.txt) || exit 1; \
  done
	@rm .pkg.coverage

godepssave:
	@echo "--> Godeps save"
	@go get github.com/tools/godep
	@go build -v -o bin/godep github.com/tools/godep
	@bin/godep save ./...

#genmocks:
#	@echo "--> Generate mocks"
#	@go build -v -o bin/mockgen github.com/golang/mock/mockgen
#	bin/mockgen -destination=./routing/logger_mock_test.go -package routing_test github.com/leanovate/microtools/logging Logger 
