export GO15VENDOREXPERIMENT=1

# variable definitions
NAME := inmars
DESC := Server statistics client
VERSION := $(shell git describe --tags --always --dirty)
GOVERSION := $(shell go version)
BUILDTIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILDDATE := $(shell date -u +"%B %d, %Y")
BUILDER := $(shell echo "`git config user.name` <`git config user.email`>")
PKG_RELEASE ?= 1
PROJECT_URL := "https://github.com/usi-lfkeitel/$(NAME)"
BUILDTAGS ?= dball
LDFLAGS := -X 'main.version=$(VERSION)' \
			-X 'main.buildTime=$(BUILDTIME)' \
			-X 'main.builder=$(BUILDER)' \
			-X 'main.goversion=$(GOVERSION)'

.PHONY: all doc fmt alltests test coverage benchmark lint vet inmars dist

all: test inmars

# development tasks
doc:
	@godoc -http=:6060 -index

fmt:
	@go fmt $$(go list ./src/...)

alltests: test lint vet

test:
	@go test -race $$(go list ./src/...)

coverage:
	@go test -cover $$(go list ./src/...)

benchmark:
	@echo "Running tests..."
	@go test -bench=. $$(go list ./src/...)

# https://github.com/golang/lint
# go get github.com/golang/lint/golint
lint:
	@golint ./src/...

vet:
	@go vet $$(go list ./src/...)

inmars:
	GOBIN=$(PWD)/bin go install -v -ldflags "$(LDFLAGS)" -tags '$(BUILDTAGS)' ./cmd/inmars

dist: vet all
	@rm -rf ./dist
	@mkdir -p dist/inmars

	@cp LICENSE dist/inmars/
	@cp README.md dist/inmars/

	@mkdir dist/inmars/bin
	@cp bin/inmars dist/inmars/bin/inmars

	(cd "dist"; tar -cz inmars) > "dist/inmars-dist-$(VERSION).tar.gz"

	@rm -rf dist/inmars