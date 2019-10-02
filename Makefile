APP:=$(notdir $(CURDIR))
GOPATH:=$(shell go env GOPATH)

GOHOSTOS:=$(shell go env GOHOSTOS)
GOHOSTARCH:=$(shell go env GOHOSTARCH)


.PHONY: all deps test coverage vet build
all: deps build

export GO111MODULE=on

${GOPATH}/bin/goverage:
	GO111MODULE=off go get -u github.com/haya14busa/goverage

GOLANGCI_LINT_VERSION:=1.17.1
${GOPATH}/bin/golangci-lint-${GOLANGCI_LINT_VERSION}:
	wget -qO $@.tar.gz https://github.com/golangci/golangci-lint/releases/download/v${GOLANGCI_LINT_VERSION}/golangci-lint-${GOLANGCI_LINT_VERSION}-${GOHOSTOS}-${GOHOSTARCH}.tar.gz
	@echo "${GOLANGCI_LINT_SHA256}  $@.tar.gz" | ${SHASUM_CHECK} -
	tar -zxf $@.tar.gz --to-stdout golangci-lint-${GOLANGCI_LINT_VERSION}-${GOHOSTOS}-${GOHOSTARCH}/golangci-lint >$@
	@chmod 700 $@
	@$@ --version
	@rm -f $@.tar.gz

deps: vendor
vendor:

build: deps
	CGO_ENABLED=0 go build -o ${APP}-${GOHOSTOS}-${GOHOSTARCH} -ldflags "-s -w" -a -installsuffix cgo .

vet: ${GOPATH}/bin/golangci-lint-${GOLANGCI_LINT_VERSION}
	$< run ./...
